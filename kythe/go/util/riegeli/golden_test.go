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

package riegeli_test

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"kythe.io/kythe/go/storage/stream"
	"kythe.io/kythe/go/util/riegeli"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	spb "kythe.io/kythe/proto/storage_go_proto"
	rmpb "kythe.io/third_party/riegeli/records_metadata_go_proto"
)

var (
	goldenJSONFile            = "testdata/golden.entries.json"
	goldenMetadataFile        = "testdata/golden.records_metadata.textproto"
	goldenRiegeliFilePrefix   = "testdata/golden.entries"
	goldenRiegeliFileVariants = []string{
		"uncompressed",
		"uncompressed_transpose",
		"brotli",
		"brotli_transpose",
	}
)

func BenchmarkGoldenTestData(b *testing.B) {
	for _, variant := range goldenRiegeliFileVariants {
		file := strings.Join([]string{goldenRiegeliFilePrefix, variant, "riegeli"}, ".")
		b.Run(variant, func(b *testing.B) { benchGoldenData(b, file) })
	}
}

// benchGoldenData benchmarks the sequential reading of a Riegeli file.  MB/s is
// measured by the size of each record read.
func benchGoldenData(b *testing.B, goldenRiegeliFile string) {
	f, err := os.Open(goldenRiegeliFile)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	b.ReportAllocs()
	rd := riegeli.NewReader(f)
	for {
		rec, err := rd.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			b.Fatal(err)
		}
		b.SetBytes(int64(len(rec)))
	}
}

type jsonReader struct{ ch <-chan *spb.Entry }

func (j *jsonReader) Next() (*spb.Entry, error) {
	e, ok := <-j.ch
	if !ok {
		return nil, io.EOF
	}
	return e, nil
}

func TestGoldenTestData(t *testing.T) {
	for _, variant := range goldenRiegeliFileVariants {
		file := strings.Join([]string{goldenRiegeliFilePrefix, variant, "riegeli"}, ".")
		t.Run(variant, func(t *testing.T) { checkGoldenData(t, file) })
	}
}

func checkGoldenData(t *testing.T, goldenRiegeliFile string) {
	jsonFile, err := os.Open(goldenJSONFile)
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()
	riegeliFile, err := os.Open(goldenRiegeliFile)
	if err != nil {
		t.Fatal(err)
	}
	defer riegeliFile.Close()

	jsonReader := &jsonReader{stream.ReadJSONEntries(jsonFile)}
	riegeliReader := riegeli.NewReader(riegeliFile)

	mdTextProto, err := ioutil.ReadFile(goldenMetadataFile)
	if err != nil {
		t.Fatalf("Error reading %s: %v", goldenMetadataFile, err)
	}
	var expectedMetadata rmpb.RecordsMetadata
	if err := proto.UnmarshalText(string(mdTextProto), &expectedMetadata); err != nil {
		t.Fatalf("Error unmarshaling %s: %v", goldenMetadataFile, err)
	}

	md, err := riegeliReader.RecordsMetadata()
	if err != nil {
		t.Fatalf("Error reading RecordsMetadata: %v", err)
	} else if diff := cmp.Diff(md, &expectedMetadata, ignoreProtoXXXFields, ignoreOptionsField); diff != "" {
		t.Errorf("Bad RecordsMetadata: (-: found; +: expected)\n%s", diff)
	}

	for {
		expected, err := jsonReader.Next()
		if err == io.EOF {
			if rec, err := riegeliReader.Next(); err != io.EOF {
				t.Fatalf("Unexpected error/record at end of Riegeli file: %q %v", hex.EncodeToString(rec), err)
			}
			return
		} else if err != nil {
			t.Fatalf("Error reading JSON golden data: %v", err)
		}

		found := &spb.Entry{}
		if err := riegeliReader.NextProto(found); err != nil {
			t.Fatalf("Error reading Riegeli golden data: %v", err)
		}

		if !proto.Equal(expected, found) {
			t.Errorf("Unexpected record: found: {%+v}; expected: {%+v}", found, expected)
		}
	}
}

var ignoreProtoXXXFields = cmp.FilterPath(func(p cmp.Path) bool {
	for _, s := range p {
		if strings.HasPrefix(s.String(), ".XXX_") {
			return true
		}
	}
	return false
}, cmp.Ignore())

var ignoreOptionsField = cmp.FilterPath(func(p cmp.Path) bool {
	for _, s := range p {
		if s.String() == ".RecordWriterOptions" {
			return true
		}
	}
	return false
}, cmp.Ignore())
