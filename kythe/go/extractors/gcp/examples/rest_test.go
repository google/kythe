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

// Binary rest_test checks whether or not you can access a project via cloud
// build rest api:
// https://cloud.google.com/cloud-build/docs/api/reference/rest/
//
// Usage:
// go build rest_test
// ./rest_test -project_id=some-project-id
//
// For authentication steps, see kythe/go/extractors/gcp/README.md
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

var (
	projectID = flag.String("project_id", "", "The GCP Cloud Build project ID to use")
)

func verifyFlags() {
	if *projectID == "" {
		log.Fatalf("Must specify valid -project_id")
	}
}

func main() {
	flag.Parse()
	verifyFlags()

	httpClient, err := google.DefaultClient(context.Background(), cloudbuild.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Failed to create oauth client: %q", err)
	}
	cbs, err := cloudbuild.New(httpClient)
	if err != nil {
		log.Fatalf("Failed to dial cloud build: %q", err)
	}

	pbs := cloudbuild.NewProjectsBuildsService(cbs)

	pbgc := pbs.List(*projectID)
	r, err := pbgc.Do()
	if err != nil {
		log.Fatalf("Failed to list projects for %s: %q", *projectID, err)
	}
	fmt.Printf("Project %s has %d builds\n", *projectID, len(r.Builds))
}
