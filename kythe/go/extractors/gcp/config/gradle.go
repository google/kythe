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

package config

import (
	"fmt"
	"path"

	"kythe.io/kythe/go/extractors/constants"

	rpb "kythe.io/kythe/proto/repo_go_proto"

	"google.golang.org/api/cloudbuild/v1"
)

func gradleStep(conf *rpb.ExtractionHint) *cloudbuild.BuildStep {
	return &cloudbuild.BuildStep{
		Name: constants.GCRGradleImage,
		Args: []string{
			"clean",
			// TODO(#3126): If compile-test has to be done as a separate
			// step, then we also need to fix this in the same way as we do
			// for multiple repo support.  Probably this would need to be
			// done with multiple steps (but making sure to not clobber
			// output).
			// The alternative here is to fall back to using clean install,
			// which should also work.
			"build",
			"-s", // Prints stacktraces for user exceptions.
			"-S", // Prints verbose stacktraces.
			"-d", // Logs in debug mode.
			"-b",
			path.Join(codeDirectory, conf.Root, "build.gradle"),
		},
		Volumes: []*cloudbuild.Volume{
			&cloudbuild.Volume{
				Name: javaVolumeName,
				Path: constants.DefaultExtractorsDir,
			},
		},
		Env: []string{
			"KYTHE_CORPUS=" + corpus,
			"KYTHE_OUTPUT_DIRECTORY=" + outputDirectory,
			fmt.Sprintf("KYTHE_OUTPUT_FILE=%s", path.Join(outputDirectory, outputFilePattern)),
			"KYTHE_ROOT_DIRECTORY=" + codeDirectory,
			"JAVAC_EXTRACTOR_JAR=" + constants.DefaultJavaExtractorLocation,
			"REAL_JAVAC=" + constants.DefaultJavacLocation,
			"TMPDIR=" + outputDirectory,
			"KYTHE_JAVAC_RUNTIME_OPTIONS=-Xbootclasspath/p:" + constants.DefaultJava9ToolsLocation,
		},
	}
}
