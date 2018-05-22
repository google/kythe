/*
 * Copyright 2018 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package kzip implements the kzip compilation storage file format.
//
// The package exports two types of interest: A kzip.Reader can be used to read
// the contents of an existing kzip archive, and a kzip.Writer can be used to
// construct a new kzip archive.
//
// Reading an Archive:
//
//   r, err := kzip.NewReader(file, size)
//   ...
//
//   // Look up a compilation record by its digest.
//   unit, err := r.Lookup(unitDigest)
//   ...
//
//   // Scan all the compilation records stored.
//   err := r.Scan(func(unit *kzip.Unit) error {
//      if hasInterestingProperty(unit) {
//         doStuffWith(unit)
//      }
//      return nil
//   })
//
//   // Open a reader for a stored file.
//   rc, err := r.Open(fileDigest)
//   ...
//   defer rc.Close()
//
//   // Read the complete contents of a stored file.
//   bits, err := r.ReadAll(fileDigest)
//   ...
//
// Writing an Archive:
//
//   w, err := kzip.NewWriter(file)
//   ...
//
//   // Add a compilation record and (optional) index data.
//   udigest, err := w.AddUnit(unit, nil)
//   ...
//
//   // Add file contents.
//   fdigest, err := w.AddFile(file)
//   ...
//
package kzip

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"kythe.io/kythe/go/platform/kcd/kythe"

	"github.com/golang/protobuf/jsonpb"

	apb "kythe.io/kythe/proto/analysis_go_proto"
)

// A Reader permits reading and scanning compilation records and file contents
// stored in a .kzip archive. The Lookup and Scan methods are mutually safe for
// concurrent use by multiple goroutines.
type Reader struct {
	zip *zip.Reader

	// The archives written by this library always use "root/" for the root
	// directory, but it's not required by the spec. Use whatever name the
	// archive actually specifies in the leading directory.
	root string
}

// NewReader constructs a new Reader that consumes zip data from r, whose total
// size in bytes is given.
func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	archive, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	if len(archive.File) == 0 {
		return nil, errors.New("archive is empty")
	} else if fi := archive.File[0].FileInfo(); !fi.IsDir() {
		return nil, errors.New("archive root is not a directory")
	}
	return &Reader{
		zip:  archive,
		root: archive.File[0].Name,
	}, nil
}

func (r *Reader) unitPath(digest string) string { return path.Join(r.root, "units", digest) }
func (r *Reader) filePath(digest string) string { return path.Join(r.root, "files", digest) }

// ErrDigestNotFound is returned when a requested compilation unit or file
// digest is not found.
var ErrDigestNotFound = errors.New("digest not found")

// ErrUnitExists is returned by AddUnit when adding the same compilation
// multiple times.
var ErrUnitExists = errors.New("unit already exists")

func readUnit(digest string, f *zip.File) (*Unit, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var msg apb.IndexedCompilation
	if err := jsonpb.Unmarshal(rc, &msg); err != nil {
		return nil, err
	}
	return &Unit{
		Digest: digest,
		Proto:  msg.Unit,
		Index:  msg.Index,
	}, nil
}

// Lookup returns the specified compilation from the archive, if it exists.  If
// the requested digest is not in the archive, ErrDigestNotFound is returned.
func (r *Reader) Lookup(unitDigest string) (*Unit, error) {
	needle := r.unitPath(unitDigest)
	for _, f := range r.zip.File {
		if f.Name == needle {
			return readUnit(unitDigest, f)
		}
	}
	return nil, ErrDigestNotFound
}

// Scan scans all the compilations stored in the archive, and invokes f for
// each compilation record. If f reports an error, the scan is terminated and
// that error is propagated to the caller of Scan.
func (r *Reader) Scan(f func(*Unit) error) error {
	prefix := r.unitPath("") + "/"
	for _, file := range r.zip.File {
		if !strings.HasPrefix(file.Name, prefix) {
			continue // skip root, files
		}
		digest := strings.TrimPrefix(file.Name, prefix)
		unit, err := readUnit(digest, file)
		if err != nil {
			return err
		}
		if err := f(unit); err != nil {
			return err
		}
	}
	return nil
}

// Open opens a reader on the contents of the specified file digest.  If the
// requested digest is not in the archive, ErrDigestNotFound is returned.  The
// caller must close the reader when it is no longer needed.
func (r *Reader) Open(fileDigest string) (io.ReadCloser, error) {
	needle := r.filePath(fileDigest)
	for _, f := range r.zip.File {
		if f.Name == needle {
			return f.Open()
		}
	}
	return nil, ErrDigestNotFound
}

// ReadAll returns the complete contents of the file with the specified digest.
// It is a convenience wrapper for Open followed by ioutil.ReadAll.
func (r *Reader) ReadAll(fileDigest string) ([]byte, error) {
	f, err := r.Open(fileDigest)
	if err == nil {
		defer f.Close()
		return ioutil.ReadAll(f)
	}
	return nil, err
}

// A Unit represents a compilation record read from a kzip archive.
type Unit struct {
	Digest string
	Proto  *apb.CompilationUnit
	Index  *apb.IndexedCompilation_Index
}

// A Writer permits construction of a .kzip archive.
type Writer struct {
	mu  sync.Mutex
	zip *zip.Writer
	fd  map[string]bool // file digests already written
	ud  map[string]bool // unit digests already written
}

// NewWriter constructs a new empty Writer that delivers output to w.  The
// AddUnit and AddFile methods are safe for use by concurrent goroutines.
func NewWriter(w io.Writer) (*Writer, error) {
	archive := zip.NewWriter(w)
	// Create an entry for the root directory, which must be first.
	root := &zip.FileHeader{
		Name:    "root/",
		Comment: "kzip root directory",
	}
	root.SetMode(os.ModeDir | 0755)
	root.SetModTime(time.Now())
	if _, err := archive.CreateHeader(root); err != nil {
		return nil, err
	}
	archive.SetComment("Kythe kzip archive")

	return &Writer{
		zip: archive,
		fd:  make(map[string]bool),
		ud:  make(map[string]bool),
	}, nil
}

// toJSON defines the encoding format for compilation messages.
var toJSON = &jsonpb.Marshaler{OrigName: true}

// AddUnit adds a new compilation record to be added to the archive, returning
// the hex-encoded SHA256 digest of the unit's contents. It is legal for index
// to be nil, in which case no index terms will be added.
//
// If the same compilation is added multiple times, AddUnit returns the digest
// of the duplicated compilation along with ErrUnitExists to all callers after
// the first. The existing unit is not modified.
func (w *Writer) AddUnit(cu *apb.CompilationUnit, index *apb.IndexedCompilation_Index) (string, error) {
	unit := kythe.Unit{Proto: cu}
	unit.Canonicalize()
	hash := sha256.New()
	unit.Digest(hash)
	digest := hex.EncodeToString(hash.Sum(nil))

	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.ud[digest]; ok {
		return digest, ErrUnitExists
	}

	f, err := w.zip.CreateHeader(newFileHeader("root", "units", digest))
	if err != nil {
		return "", err
	}
	if err := toJSON.Marshal(f, &apb.IndexedCompilation{
		Unit:  unit.Proto,
		Index: index,
	}); err != nil {
		return "", err
	}
	w.ud[digest] = true
	return digest, nil
}

// AddFile copies the complete contents of r into the archive as a new file
// entry, returning the hex-encoded SHA256 digest of the file's contents.
func (w *Writer) AddFile(r io.Reader) (string, error) {
	// Buffer the file contents and compute their digest.
	// We have to do this ahead of time, because we have to provide the name of
	// the file before we can start writing its contents.
	var buf bytes.Buffer
	hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(hash, &buf), r); err != nil {
		return "", err
	}
	digest := hex.EncodeToString(hash.Sum(nil))

	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.fd[digest]; ok {
		return digest, nil // already written
	}

	f, err := w.zip.CreateHeader(newFileHeader("root", "files", digest))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, &buf); err != nil {
		return "", err
	}
	w.fd[digest] = true
	return digest, nil
}

// Close closes the writer, flushing any remaining unwritten data out to the
// underlying zip file. It is safe to close w arbitrarily many times; all calls
// after the first will report nil.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.zip != nil {
		err := w.zip.Close()
		w.zip = nil
		return err
	}
	return nil
}

func newFileHeader(parts ...string) *zip.FileHeader {
	fh := &zip.FileHeader{Name: path.Join(parts...), Method: zip.Deflate}
	fh.SetMode(0600)
	fh.SetModTime(time.Now())
	return fh
}
