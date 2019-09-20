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

package info

import (
	"testing"

	"github.com/golang/protobuf/proto"

	apb "kythe.io/kythe/proto/analysis_go_proto"
)

func TestMergeKzipInfo(t *testing.T) {
	infos := []*apb.KzipInfo{
		{
			TotalUnits: 1,
			TotalFiles: 2,
			Corpora: map[string]*apb.KzipInfo_CorpusInfo{
				"corpus1": {
					Files: map[string]int32{"python": 2},
					Units: map[string]int32{"python": 1},
				},
			},
		},
		{
			TotalUnits: 1,
			TotalFiles: 2,
			Corpora: map[string]*apb.KzipInfo_CorpusInfo{
				"corpus1": {
					Files: map[string]int32{"go": 2},
					Units: map[string]int32{"python": 1},
				},
			},
		},
	}

	want := &apb.KzipInfo{
		TotalUnits: 2,
		TotalFiles: 4,
		Corpora: map[string]*apb.KzipInfo_CorpusInfo{
			"corpus1": {
				Files: map[string]int32{
					"go":     2,
					"python": 2,
				},
				Units: map[string]int32{"python": 2},
			},
		},
	}

	got := MergeKzipInfo(infos)
	if !proto.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
