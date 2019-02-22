/*
 * Copyright 2019 The Kythe Authors. All rights reserved.
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

// Package createcmd provides the kzip command for creating simple kzip archives.
package createcmd

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"strings"

	"kythe.io/kythe/go/platform/kzip"
	"kythe.io/kythe/go/util/cmdutil"
	"kythe.io/kythe/go/util/kytheuri"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/subcommands"

	anypb "github.com/golang/protobuf/ptypes/any"
	apb "kythe.io/kythe/proto/analysis_go_proto"
	spb "kythe.io/kythe/proto/storage_go_proto"
)

type createCommand struct {
	cmdutil.Info

	uri     string
	output  string
	source  string
	details string
}

// New creates a new subcommand for merging kzip files.
func New() subcommands.Command {
	return &createCommand{
		Info: cmdutil.NewInfo("create", "create simple kzip archives", `-uri <uri> -output <path> -source <path> required_input*

Construct a kzip file consisting of a single input file name by -source.
The resulting file is written to -output, and the vname of the compilation
will be attributed to the values specified within -uri.

Additional required inputs, if any, can be provided as positional parameters.
`),
	}
}

// SetFlags implements the subcommands interface and provides command-specific flags
// for creating kzip files.
func (c *createCommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.uri, "uri", "", "A Kythe URI naming the compilation unit")
	fs.StringVar(&c.output, "output", "", "Path for output kzip file")
	fs.StringVar(&c.source, "source", "", "Path for input source file")
	fs.StringVar(&c.details, "details", "", "JSON-encoded Any messages to embed as compilation details.")
}

// Execute implements the subcommands interface and creates the requested file.
func (c *createCommand) Execute(ctx context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	switch {
	case c.uri == "":
		return c.Fail("missing required -uri")
	case c.output == "":
		return c.Fail("missing required -output")
	case c.source == "":
		return c.Fail("missing required -source")
	}

	// Create a new compilation populating its VName with the values specified
	// within the specified Kythe URI.
	uri, err := kytheuri.Parse(c.uri)
	if err != nil {
		return c.Fail("error parsing -uri: %v", err)
	}
	cu := &apb.CompilationUnit{
		VName: &spb.VName{
			Corpus:    uri.Corpus,
			Language:  uri.Language,
			Signature: uri.Signature,
			Root:      uri.Root,
			Path:      uri.Path,
		},
		SourceFile: []string{c.source},
	}

	details, err := unmarshalDetails(c.details)
	if err != nil {
		return c.Fail("error parsing -details: %v", err)
	}
	cu.Details = details

	out, err := openWriter(c.output)
	if err != nil {
		return c.Fail("error opening -output: %v", err)
	}

	for _, path := range append([]string{c.source}, fs.Args()...) {
		err := addFile(out, cu, path)
		if err != nil {
			return c.Fail("error adding file to -output: %v", err)
		}
	}

	_, err = out.AddUnit(cu, nil)
	if err != nil {
		return c.Fail("error writing compilation to -output: %v", err)
	}

	if err := out.Close(); err != nil {
		return c.Fail("error closing output file: %v", err)
	}
	return subcommands.ExitSuccess
}

func openWriter(path string) (*kzip.Writer, error) {
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return kzip.NewWriteCloser(out)
}

func addFile(out *kzip.Writer, cu *apb.CompilationUnit, path string) error {
	input, err := os.Open(path)
	if err != nil {
		return err
	}
	defer input.Close()

	digest, err := out.AddFile(input)
	if err != nil {
		return err
	}
	cu.RequiredInput = append(cu.RequiredInput, &apb.CompilationUnit_FileInput{
		VName: &spb.VName{
			Corpus: cu.VName.Corpus,
			Root:   cu.VName.Root,
			Path:   path,
		},
		Info: &apb.FileInfo{
			Path:   path,
			Digest: digest,
		},
	})
	return nil
}

func unmarshalDetails(details string) ([]*anypb.Any, error) {
	var msg []*anypb.Any
	dec := json.NewDecoder(strings.NewReader(details))
	for dec.More() {
		var detail anypb.Any
		if err := jsonpb.UnmarshalNext(dec, &detail); err != nil {
			return msg, err
		}
		msg = append(msg, &detail)
	}
	return msg, nil
}
