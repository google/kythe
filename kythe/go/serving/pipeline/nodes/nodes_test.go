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

package nodes

import (
	"testing"

	"kythe.io/kythe/go/util/schema/edges"
	"kythe.io/kythe/go/util/schema/facts"
	"kythe.io/kythe/go/util/schema/nodes"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/passert"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/ptest"
	"github.com/apache/beam/sdks/go/pkg/beam/x/debug"

	ppb "kythe.io/kythe/proto/pipeline_go_proto"
	scpb "kythe.io/kythe/proto/schema_go_proto"
	spb "kythe.io/kythe/proto/storage_go_proto"
)

func TestFromEntries(t *testing.T) {
	entries := []*spb.Entry{{
		Source:    &spb.VName{Signature: "node1"},
		FactName:  facts.NodeKind,
		FactValue: []byte(nodes.Record),
	}, {
		Source:    &spb.VName{Signature: "node1"},
		FactName:  facts.Subkind,
		FactValue: []byte(nodes.Class),
	}, {
		Source:    &spb.VName{Signature: "node2"},
		FactName:  facts.NodeKind,
		FactValue: []byte("unknown_nodekind"),
	}, {
		Source:    &spb.VName{Signature: "node2"},
		FactName:  facts.Subkind,
		FactValue: []byte("unknown_subkind"),
	}, {
		Source:    &spb.VName{Signature: "node2"},
		FactName:  facts.Text, // schema-known fact name
		FactValue: []byte("text"),
	}, {
		// Duplicate fact
		Source:    &spb.VName{Signature: "node2"},
		FactName:  facts.Text,
		FactValue: []byte("text"),
	}, {
		Source:    &spb.VName{Signature: "node2"},
		FactName:  "/unknown/fact/name",
		FactValue: []byte("blah"),
	}, {
		Source:   &spb.VName{Signature: "node2"},
		EdgeKind: edges.Typed, // schema-known edge kind
		Target:   &spb.VName{Signature: "node1"},
	}, {
		// Duplicate edge
		Source:   &spb.VName{Signature: "node2"},
		EdgeKind: edges.Typed,
		Target:   &spb.VName{Signature: "node1"},
	}, {
		Source:   &spb.VName{Signature: "node2"},
		EdgeKind: "/unknown/edge/kind",
		Target:   &spb.VName{Signature: "node2"},
	}, {
		Source:    &spb.VName{Signature: "anchor"},
		FactName:  facts.NodeKind,
		FactValue: []byte(nodes.Anchor),
	}, {
		Source:    &spb.VName{Signature: "anchor"},
		FactName:  facts.AnchorStart,
		FactValue: []byte("1"),
	}, {
		Source:    &spb.VName{Signature: "anchor"},
		FactName:  facts.SnippetEnd,
		FactValue: []byte("5"),
	}, {
		Source:    &spb.VName{Signature: "anchor"},
		FactName:  facts.SnippetStart,
		FactValue: []byte("0"),
	}, {
		Source:    &spb.VName{Signature: "anchor"},
		FactName:  facts.AnchorEnd,
		FactValue: []byte("2"),
	}}
	expected := []*ppb.Node{{
		Source:  &spb.VName{Signature: "node1"},
		Kind:    &ppb.Node_KytheKind{scpb.NodeKind_RECORD},
		Subkind: &ppb.Node_KytheSubkind{scpb.Subkind_CLASS},
	}, {
		Source:  &spb.VName{Signature: "node2"},
		Kind:    &ppb.Node_GenericKind{"unknown_nodekind"},
		Subkind: &ppb.Node_GenericSubkind{"unknown_subkind"},
		Fact: []*ppb.Fact{{
			Name:  &ppb.Fact_GenericName{"/unknown/fact/name"},
			Value: []byte("blah"),
		}, {
			Name:  &ppb.Fact_KytheName{scpb.FactName_TEXT},
			Value: []byte("text"),
		}},
		Edge: []*ppb.Edge{{
			Kind:   &ppb.Edge_GenericKind{"/unknown/edge/kind"},
			Target: &spb.VName{Signature: "node2"},
		}, {
			Kind:   &ppb.Edge_KytheKind{scpb.EdgeKind_TYPED},
			Target: &spb.VName{Signature: "node1"},
		}},
	}, {
		Source: &spb.VName{Signature: "anchor"},
		Kind:   &ppb.Node_KytheKind{scpb.NodeKind_ANCHOR},
		Fact: []*ppb.Fact{{
			Name:  &ppb.Fact_KytheName{scpb.FactName_LOC_END},
			Value: []byte("2"),
		}, {
			Name:  &ppb.Fact_KytheName{scpb.FactName_LOC_START},
			Value: []byte("1"),
		}, {
			Name:  &ppb.Fact_KytheName{scpb.FactName_SNIPPET_END},
			Value: []byte("5"),
		}, {
			Name:  &ppb.Fact_KytheName{scpb.FactName_SNIPPET_START},
			Value: []byte("0"),
		}},
	}}

	p, s := beam.NewPipelineWithRoot()
	nodes := FromEntries(s, beam.CreateList(s, entries))
	debug.Print(s, nodes)
	passert.Equals(s, nodes, beam.CreateList(s, expected))

	if err := ptest.Run(p); err != nil {
		t.Fatalf("Pipeline error: %+v", err)
	}
}

func TestFilter(t *testing.T) {
	nodes := []*ppb.Node{{
		Source:  &spb.VName{Signature: "node1"},
		Kind:    &ppb.Node_KytheKind{scpb.NodeKind_RECORD},
		Subkind: &ppb.Node_KytheSubkind{scpb.Subkind_CLASS},
	}, {
		Source:  &spb.VName{Signature: "node2"},
		Kind:    &ppb.Node_GenericKind{"unknown_nodekind"},
		Subkind: &ppb.Node_GenericSubkind{"unknown_subkind"},
		Fact: []*ppb.Fact{{
			Name:  &ppb.Fact_KytheName{scpb.FactName_TEXT},
			Value: []byte("text"),
		}, {
			Name:  &ppb.Fact_GenericName{"/unknown/fact/name"},
			Value: []byte("blah"),
		}},
		Edge: []*ppb.Edge{{
			Kind:   &ppb.Edge_KytheKind{scpb.EdgeKind_TYPED},
			Target: &spb.VName{Signature: "node1"},
		}, {
			Kind:   &ppb.Edge_GenericKind{"/unknown/edge/kind"},
			Target: &spb.VName{Signature: "node2"},
		}},
	}}

	tests := []struct {
		filter   Filter
		expected []*ppb.Node
	}{
		{Filter{}, nodes},
		{Filter{FilterByKind: []string{"record", "unknown_nodekind"}}, nodes},
		{Filter{FilterBySubkind: []string{"class", "unknown_subkind"}}, nodes},

		{Filter{FilterByKind: []string{"record"}}, []*ppb.Node{nodes[0]}},
		{Filter{FilterByKind: []string{"unknown_nodekind"}}, []*ppb.Node{nodes[1]}},
		{Filter{FilterBySubkind: []string{"class"}}, []*ppb.Node{nodes[0]}},
		{Filter{FilterBySubkind: []string{"unknown_subkind"}}, []*ppb.Node{nodes[1]}},

		{Filter{
			FilterByKind: []string{"unknown_nodekind"},
			IncludeFacts: []string{}, // exclude all facts
		}, []*ppb.Node{{
			Source:  nodes[1].Source,
			Kind:    nodes[1].Kind,
			Subkind: nodes[1].Subkind,
			Edge:    nodes[1].Edge,
		}}},
		{Filter{
			FilterByKind: []string{"unknown_nodekind"},
			IncludeEdges: []string{}, // exclude all edges
		}, []*ppb.Node{{
			Source:  nodes[1].Source,
			Kind:    nodes[1].Kind,
			Subkind: nodes[1].Subkind,
			Fact:    nodes[1].Fact,
		}}},

		{Filter{
			FilterByKind: []string{"unknown_nodekind"},
			IncludeFacts: []string{"/kythe/text"},
			IncludeEdges: []string{},
		}, []*ppb.Node{{
			Source:  nodes[1].Source,
			Kind:    nodes[1].Kind,
			Subkind: nodes[1].Subkind,
			Fact: []*ppb.Fact{{
				Name:  &ppb.Fact_KytheName{scpb.FactName_TEXT},
				Value: []byte("text"),
			}},
		}}},
		{Filter{
			FilterByKind: []string{"unknown_nodekind"},
			IncludeFacts: []string{"/unknown/fact/name"},
			IncludeEdges: []string{},
		}, []*ppb.Node{{
			Source:  nodes[1].Source,
			Kind:    nodes[1].Kind,
			Subkind: nodes[1].Subkind,
			Fact: []*ppb.Fact{{
				Name:  &ppb.Fact_GenericName{"/unknown/fact/name"},
				Value: []byte("blah"),
			}},
		}}},

		{Filter{
			FilterByKind: []string{"unknown_nodekind"},
			IncludeFacts: []string{},
			IncludeEdges: []string{"/kythe/edge/typed"},
		}, []*ppb.Node{{
			Source:  nodes[1].Source,
			Kind:    nodes[1].Kind,
			Subkind: nodes[1].Subkind,
			Edge: []*ppb.Edge{{
				Kind:   &ppb.Edge_KytheKind{scpb.EdgeKind_TYPED},
				Target: &spb.VName{Signature: "node1"},
			}},
		}}},
		{Filter{
			FilterByKind: []string{"unknown_nodekind"},
			IncludeFacts: []string{},
			IncludeEdges: []string{"/unknown/edge/kind"},
		}, []*ppb.Node{{
			Source:  nodes[1].Source,
			Kind:    nodes[1].Kind,
			Subkind: nodes[1].Subkind,
			Edge: []*ppb.Edge{{
				Kind:   &ppb.Edge_GenericKind{"/unknown/edge/kind"},
				Target: &spb.VName{Signature: "node2"},
			}},
		}}},
	}

	for _, test := range tests {
		p, s, coll := ptest.CreateList(nodes)
		filtered := beam.ParDo(s, &test.filter, coll)
		debug.Print(s, filtered)
		passert.Equals(s, filtered, beam.CreateList(s, test.expected))

		if err := ptest.Run(p); err != nil {
			t.Fatalf("Pipeline error: %+v", err)
		}
	}
}
