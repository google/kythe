/*
 * Copyright 2018 The Kythe Authors. All rights reserved.
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

// Package validation contains logic for verifing the contents of a kzip.
package validation

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kythe.io/kythe/go/platform/kzip"

	"bitbucket.org/creachadair/stringset"
)

// Harness encapsulates data for validation of a single kzip and repo.
type Harness struct {
	// Path to the kzip file to validate.
	Compilations []string
	// Path to the root of the repo to compare the kzip against.
	Repo string
	// Language extensions to check, e.g. java, cpp, h.
	Langs stringset.Set
	// Optional file to output missing paths to.
	MissingOutput *string
}

// Result explains a validation outcome.
type Result struct {
	// Number of files in the kzip file matching one of Langs extensions.
	FilesInKZIP int
	// Number of files in the repo matching one of Langs extensions.
	FilesInRepo int
	// Number of files in the repo that were not matched in the kzip.
	Missing int
	// A breakdown of subdirectories with the most missing files.
	TopMissing Coverage
}

// Coverage describes the missing subdirectories in a repo.
type Coverage struct {
	// The missing paths.
	Paths []Dir
}

// The number of missing paths to print in output.
const topMissingPathsLimit = 20
const topMissingPathsDepth = 5

// Dir describes a single missing subdirectory and how much is missing.
type Dir struct {
	// Path is a subdirectory from Repo root.
	Path string
	// Missing is the number of files in the Repo from the corresponding
	// subdirectory that were missing from the KZIP.
	Missing int
}

// Len is part of sort.Interface.
func (c *Coverage) Len() int {
	return len(c.Paths)
}

// Swap is part of sort.Interface.
func (c *Coverage) Swap(i, j int) {
	c.Paths[i], c.Paths[j] = c.Paths[j], c.Paths[i]
}

// Less is part of sort.Interface.
func (c *Coverage) Less(i, j int) bool {
	return c.Paths[i].Missing < c.Paths[j].Missing
}

// Validate opens a kzip file and compares all of its contained files with those
// present in the provided repo path.  It returns results explaining what
// percentage of files are covered in the kzip.
func (h Harness) Validate() (*Result, error) {
	log.Printf("Gathering validation data for %s", h.Compilations)
	fromKZIP, err := h.filenamesFromCompilations(h.Compilations)
	if err != nil {
		return nil, fmt.Errorf("reading from kzip %s: %v", h.Compilations, err)
	}
	fromRepo, err := h.filenamesFromPath(h.Repo)
	if err != nil {
		return nil, fmt.Errorf("reading from repo %s: %v", h.Repo, err)
	}
	missingFromKZIP := stringset.FromKeys(fromRepo).Diff(stringset.FromKeys(fromKZIP))

	// Try to calculate the most missing subdirectories.
	topMissing := Coverage{}
	// Only try to go so many paths deep (src/java/com/domain/sub).
	// TODO(danielmoy): be more intelligent about this.
	for i := 1; i <= topMissingPathsDepth; i++ {
		topMissing = getTopMissing(missingFromKZIP, topMissingPathsLimit, i)
		if len(topMissing.Paths) >= topMissingPathsLimit {
			break
		}
	}

	if h.MissingOutput != nil {
		err = outputMissing(missingFromKZIP, *h.MissingOutput)
	}

	return &Result{
		FilesInKZIP: len(fromKZIP),
		FilesInRepo: len(fromRepo),
		Missing:     len(missingFromKZIP),
		TopMissing:  topMissing,
	}, err
}

func (h Harness) filenamesFromCompilations(compilations []string) (stringset.Set, error) {
	ret := stringset.New()
	for _, path := range compilations {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening kzip %s: %v", path, err)
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("getting kzip size %s: %v", path, err)
		}
		r, err := kzip.NewReader(f, fi.Size())
		if err != nil {
			return nil, fmt.Errorf("opening kzip reader %s: %v", path, err)
		}

		if err := r.Scan(func(cu *kzip.Unit) error {
			if cu.Proto != nil {
				for _, f := range cu.Proto.SourceFile {
					if h.Langs.Contains(getExt(f)) {
						ret.Add(f)
					}
				}
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("scanning kzip %s: %v", path, err)
		}
	}

	return ret, nil
}

func (h Harness) filenamesFromPath(path string) (stringset.Set, error) {
	ret := stringset.New()
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if h.Langs.Contains(getExt(p)) {
			// We are only interested in the repo-relative path.
			rp, err := filepath.Rel(path, p)
			if err != nil {
				return err
			}
			ret.Add(rp)
		}
		return nil
	})
	return ret, err
}

// Gets the extension of a file, agnostic of a preceding "."
func getExt(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}

// getTopMissing returns the top num most common subdirectories, down to at most
// level.
func getTopMissing(missing stringset.Set, num, level int) Coverage {
	counts := map[string]int{}
	for f := range missing {
		dirs := strings.Split(f, string(os.PathSeparator))
		if len(dirs) > level-1 {
			counts[strings.Join(dirs[:level-1], string(os.PathSeparator))]++
		} else {
			counts[f]++
		}
	}

	top := Coverage{}
	for f, c := range counts {
		top.Paths = append(top.Paths, Dir{f, c})
	}
	sort.Sort(sort.Reverse(&top))
	if len(top.Paths) >= num {
		top.Paths = top.Paths[:num]
	}
	return top
}

func outputMissing(filenames stringset.Set, output string) (err error) {
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("opening missing file output %s: %v", output, err)
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	for name := range filenames {
		f.WriteString(name + "\n")
	}

	return
}
