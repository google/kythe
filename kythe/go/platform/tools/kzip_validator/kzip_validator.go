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

// Binary kzip_validator checks the contents of a kzip against a code repo and
// compares file coverage.
//
// Example:
//  kzip_validator -kzip <kzip-file> -local_repo <repo-root-dir> -lang cc,h
//  kzip_validator -kzip <kzip-file> -repo_url <repo-url> -lang java [-version <hash>]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"kythe.io/kythe/go/extractors/validation"

	"bitbucket.org/creachadair/stringset"
)

var (
	kzip      = flag.String("kzip", "", "The comma separated list of kzip files to check")
	repoURL   = flag.String("repo_url", "", "The repo to compare against")
	version   = flag.String("version", "", "The version of the remote repo to compare")
	localRepo = flag.String("local_repo", "", "The path of an optional local repo to specify instead of -repo_url")
	lang      = flag.String("lang", "", "The comma-separated language file extensions to check, e.g. 'java' or 'cpp,h'")

	missingFile = flag.String("missing_file", "", "An optional file to write all missing filepaths to")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s -kzip <kzip-file> -repo_url <url> [-version <hash>]

Compare a kzip file's contents with a given repo, and print results of file coverage.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func verifyFlags() {
	if len(flag.Args()) > 0 {
		log.Fatalf("Unknown arguments: %v", flag.Args())
	}

	if *kzip == "" {
		log.Fatal("You must provide at least one -kzip file")
	}
	if *repoURL == "" && *localRepo == "" {
		log.Fatal("You must specify either -repo_url or a -local_repo path")
	}
	if *repoURL != "" && *localRepo != "" {
		log.Fatal("You must specify only one of -repo_url or -local_repo")
	}
	if *version != "" && *repoURL == "" {
		log.Fatal("-version is only copmatible with -reop_url, not -local_repo")
	}
	if *lang == "" {
		log.Fatal("You must specify what -lang file extensions to examine")
	}
}

func main() {
	flag.Parse()
	verifyFlags()

	fmt.Printf("Validating kzip %s\n", *kzip)

	repoPath := fetchRepo()

	res, err := validation.Settings{
		Compilations:  strings.Split(*kzip, ","),
		Repo:          repoPath,
		Langs:         stringset.FromKeys(strings.Split(*lang, ",")),
		MissingOutput: missingFile,
	}.Validate()
	if err != nil {
		log.Fatalf("Failure validating: %v", err)
	}

	log.Printf("Result: %v", res)

	fmt.Println("KZIP verification:")
	fmt.Printf("KZIP file count: %d\n", res.NumArchiveFiles)
	fmt.Printf("Repo file count: %d\n", res.NumRepoFiles)
	fmt.Printf("Percent missing: %.1f%%\n", 100*float64(res.NumMissing)/float64(res.NumRepoFiles))
	if len(res.TopMissing.Paths) > 0 {
		fmt.Println("Top missing subdirectories:")
		for _, v := range res.TopMissing.Paths {
			fmt.Printf("%10d %s\n", v.Missing, v.Path)
		}
	}
}

// fetchRepo either just early returns an available -local_repo, or fetches it from
// the specified remote -repo_url
func fetchRepo() string {
	if *localRepo != "" {
		log.Printf("Comparing against local copy of repo %s", *localRepo)
		return *localRepo
	}
	log.Fatal("Unsupported use of -repo_url")
	log.Printf("Comparing against remote repo %s", *repoURL)
	return ""
}
