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

package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	ecpb "kythe.io/kythe/proto/extraction_config_go_proto"
)

const testDataDir = "testdata"

var multipleNewLines = regexp.MustCompile("\n{2,}")

func imagesEqual(got, want []byte) (bool, string) {
	// remove superfluous whitespace
	gotStr := strings.Trim(string(got[:]), " \n")
	wantStr := strings.Trim(string(want[:]), " \n")
	gotStr = multipleNewLines.ReplaceAllString(gotStr, "\n")
	wantStr = multipleNewLines.ReplaceAllString(wantStr, "\n")

	// diff want vs got
	diff := cmp.Diff(gotStr, wantStr)
	if diff != "" {
		return false, diff
	}

	return true, ""
}

func mustLoadDockerFile(t *testing.T, testConfigFile string) []byte {
	t.Helper()
	fileName := fmt.Sprintf("expected_%s.Dockerfile", strings.Replace(filepath.Base(testConfigFile), ".json", "", 1))
	content, err := ioutil.ReadFile(os.ExpandEnv(filepath.Join(testDataDir, fileName)))
	if err != nil {
		t.Fatalf("Failed to open test docker file: %v\n", err)
	}

	return content
}

type testFiles []*os.File

func (t testFiles) Close() error {
	var err error
	for _, f := range t {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}
	return err
}

func mustOpenTestData(t *testing.T) testFiles {
	t.Helper()
	fileNames, err := filepath.Glob(os.ExpandEnv(fmt.Sprintf("%s/%s", testDataDir, "*.json")))
	if err != nil {
		t.Fatalf("Failed to glob for test data files: %v\n", err)
	}

	if len(fileNames) == 0 {
		t.Fatal("No test config data found!\n")
	}

	var files []*os.File
	for _, fileName := range fileNames {
		file, err := os.Open(os.ExpandEnv(fileName))
		if err != nil {
			t.Fatalf("Failed to load test data: %v", err)
		}

		files = append(files, file)
	}

	return files
}

func TestNewImageGeneratesExpectedDockerFiles(t *testing.T) {
	testData := mustOpenTestData(t)
	defer testData.Close()
	for _, file := range testData {
		config, err := Load(file)
		if err != nil {
			t.Fatalf("Failed to load extraction config: %v", err)
		}

		got, err := newImage(config, imageSettings{})
		if err != nil {
			t.Fatalf("Failed to parse test config: %v\n", err)
		}

		want := mustLoadDockerFile(t, file.Name())
		if eq, diff := imagesEqual(want, got); !eq {

			t.Fatalf("Images were not equal, diff:\n%s", diff)
		}
	}
}

func TestLoadReturnsProperData(t *testing.T) {
	testData := mustOpenTestData(t)
	defer testData.Close()
	for _, file := range testData {
		got, err := Load(file)
		if err != nil {
			t.Fatalf("Failed to load extraction config: %v", err)
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatalf("Failed to seek test data: %v", err)
		}
		want := &ecpb.ExtractionConfiguration{}
		if err := jsonpb.Unmarshal(file, want); err != nil {
			t.Fatalf("Failed to unmarshal test data: %v", err)
		}

		if !proto.Equal(got, want) {
			t.Fatalf("Expected: %v\nGot: %v", want, got)
		}
	}
}

func TestCreateImageWritesProperData(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "tempOutputDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testData := mustOpenTestData(t)
	defer testData.Close()
	for _, file := range testData {
		config, err := Load(file)
		if err != nil {
			t.Fatalf("Failed to load extraction config: %v", err)
		}

		tmpImageFile, err := ioutil.TempFile(tempDir, "tmpExtractionImage.Docker")
		if err != nil {
			t.Fatalf("Failed to create temp image file: %v", err)
		}

		err = CreateImage(config, imageSettings{}, tmpImageFile.Name())
		if err != nil {
			t.Fatalf("Failed to create image: %v", err)
		}

		got, err := ioutil.ReadAll(tmpImageFile)
		if err != nil {
			t.Fatalf("Failed to read created image: %v", err)
		}

		want := mustLoadDockerFile(t, file.Name())
		if eq, diff := imagesEqual(got, want); !eq {
			t.Fatalf("Images were not equal, diff:\n%s\n", diff)
		}
	}
}
