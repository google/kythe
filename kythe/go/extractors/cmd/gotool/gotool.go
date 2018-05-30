/*
 * Copyright 2015 Google Inc. All rights reserved.
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

// Binary gotool extracts Kythe compilation information for Go packages named
// by import path on the command line.  The output compilations are written
// into an index pack directory.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	"kythe.io/kythe/go/extractors/golang"
	"kythe.io/kythe/go/platform/indexpack"
	"kythe.io/kythe/go/platform/kindex"
	"kythe.io/kythe/go/platform/kzip"
	"kythe.io/kythe/go/platform/vfs"

	"github.com/pborman/uuid"

	apb "kythe.io/kythe/proto/analysis_go_proto"
)

var (
	bc = build.Default // A shallow copy of the default build settings

	corpus     = flag.String("corpus", "", "Default corpus name to use")
	localPath  = flag.String("local_path", "", "Directory where relative imports are resolved")
	outputPath = flag.String("output", "", "Output path (indexpack directory or .kzip filename)")
	extraFiles = flag.String("extra_files", "", "Additional files to include in each compilation (CSV)")
	indexFiles = flag.Bool("kindex", false, "Write outputs to .kindex files")
	byDir      = flag.Bool("bydir", false, "Import by directory rather than import path")
	keepGoing  = flag.Bool("continue", false, "Continue past errors")
	verbose    = flag.Bool("v", false, "Enable verbose logging")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [options] <import-path>...
Extract Kythe compilation records from Go import paths specified on the command line.
Outputs are written to an index pack unless --kindex is set, in which case they
are written to individual .kindex files in the output directory.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	// Attach flags to the various parts of the go/build context we are using.
	// These will override the system defaults from the environment.
	flag.StringVar(&bc.GOARCH, "goarch", bc.GOARCH, "Go system architecture tag")
	flag.StringVar(&bc.GOOS, "goos", bc.GOOS, "Go operating system tag")
	flag.StringVar(&bc.GOPATH, "gopath", bc.GOPATH, "Go library path")
	flag.StringVar(&bc.GOROOT, "goroot", bc.GOROOT, "Go system root")
	flag.BoolVar(&bc.CgoEnabled, "gocgo", bc.CgoEnabled, "Whether to allow cgo")
	flag.StringVar(&bc.Compiler, "gocompiler", bc.Compiler, "Which Go compiler to use")

	// TODO(fromberger): Attach flags to the build and release tags (maybe).
}

func maybeFatal(msg string, args ...interface{}) {
	log.Printf(msg, args...)
	if !*keepGoing {
		os.Exit(1)
	}
}

func maybeLog(msg string, args ...interface{}) {
	if *verbose {
		log.Printf(msg, args...)
	}
}

func main() {
	flag.Parse()

	if *outputPath == "" {
		log.Fatal("You must provide a non-empty --output path")
	}

	ctx := context.Background()
	ext := &golang.Extractor{
		BuildContext: bc,
		Corpus:       *corpus,
		LocalPath:    *localPath,
	}
	if *extraFiles != "" {
		ext.ExtraFiles = strings.Split(*extraFiles, ",")
	}

	locate := ext.Locate
	if *byDir {
		locate = ext.ImportDir
	}
	for _, path := range flag.Args() {
		pkg, err := locate(path)
		if err == nil {
			maybeLog("Found %q in %s", pkg.Path, pkg.BuildPackage.Dir)
		} else {
			maybeFatal("Error locating %q: %v", path, err)
		}
	}

	if err := ext.Extract(); err != nil {
		maybeFatal("Error in extraction: %v", err)
	}

	var writer packageWriter
	if *indexFiles {
		writer = writeToIndex(ctx, *outputPath)
	} else if filepath.Ext(*outputPath) == ".kzip" {
		writer = writeToKZip(ctx, *outputPath)
	} else {
		writer = writeToPack(ctx, *outputPath)
	}
	maybeLog("Writing %d package(s) to %q", len(ext.Packages), *outputPath)
	for _, pkg := range ext.Packages {
		maybeLog("Package %q:\n\t// %s", pkg.Path, pkg.BuildPackage.Doc)
		if err := writer.writePackage(ctx, pkg); err != nil {
			maybeFatal("Error writing %q: %v", pkg.Path, err)
		}
	}
	if err := writer.close(); err != nil {
		maybeFatal("Error closing output: %v", err)
	}
}

type packageWriter interface {
	writePackage(context.Context, *golang.Package) error
	close() error
}

// writeToKZip returns a packageWriter that stores compilations into a kzip
// file at the specified path.
func writeToKZip(ctx context.Context, path string) packageWriter {
	if err := vfs.MkdirAll(ctx, filepath.Dir(path), 0755); err != nil {
		log.Fatalf("Unable to create output directory: %v", err)
	}
	f, err := vfs.Create(ctx, path)
	if err != nil {
		log.Fatalf("Unable to create output file: %v", err)
	}
	w, err := kzip.NewWriter(f)
	if err != nil {
		log.Fatalf("Unable to create kzip writer: %v", err)
	}
	return kzipWriter{w}
}

type kzipWriter struct{ w *kzip.Writer }

func (k kzipWriter) writePackage(ctx context.Context, pkg *golang.Package) error {
	return pkg.EachUnit(ctx, func(cu *kindex.Compilation) error {
		if _, err := k.w.AddUnit(cu.Proto, nil); err != nil {
			return err
		}
		for _, fd := range cu.Files {
			if _, err := k.w.AddFile(bytes.NewReader(fd.Content)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (k kzipWriter) close() error { return k.w.Close() }

// writeToIndex returns a packageWriter that stores compilations as kindex
// files under the specified directory path.
func writeToIndex(ctx context.Context, path string) packageWriter {
	if err := vfs.MkdirAll(ctx, path, 0755); err != nil {
		log.Fatalf("Unable to create output directory: %v", err)
	}
	return kindexWriter(func(ctx context.Context, pkg *golang.Package) error {
		return pkg.EachUnit(ctx, func(cu *kindex.Compilation) error {
			path := filepath.Join(path, uuid.New()+".kindex")
			f, err := vfs.Create(ctx, path)
			if err != nil {
				return err
			}
			_, err = cu.WriteTo(f)
			cerr := f.Close()
			if err != nil {
				return err
			}
			return cerr
		})
	})
}

type kindexWriter func(context.Context, *golang.Package) error

func (k kindexWriter) writePackage(ctx context.Context, pkg *golang.Package) error {
	return k(ctx, pkg)
}

func (kindexWriter) close() error { return nil }

// writeToPack returns a packageWriter that stores compilations to an index
// pack under the specified directory path.
func writeToPack(ctx context.Context, path string) packageWriter {
	pack, err := indexpack.CreateOrOpen(ctx, path, indexpack.UnitType((*apb.CompilationUnit)(nil)))
	if err != nil {
		log.Fatalf("Unable to open %q: %v", path, err)
	}
	return packWriter{pack}
}

type packWriter struct{ w *indexpack.Archive }

func (p packWriter) writePackage(ctx context.Context, pkg *golang.Package) error {
	_, err := pkg.Store(ctx, p.w)
	return err
}

func (packWriter) close() error { return nil }
