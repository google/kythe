"""Rules for verifying proto indexer output"""

# Copyright 2019 The Kythe Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@bazel_skylib//lib:collections.bzl", "collections")
load(
    "@io_kythe//tools/build_rules/verifier_test:verifier_test.bzl",
    "KytheVerifierSources",
    "extract",
    "index_compilation",
    "verifier_test",
)
load("@rules_proto//proto:defs.bzl", "ProtoInfo")

def _invoke(rulefn, name, **kwargs):
    """Invoke rulefn with name and kwargs, returning the label of the rule."""
    rulefn(name = name, **kwargs)
    return "//{}:{}".format(native.package_name(), name)

def get_proto_files_and_proto_path_opts(protolibs):
    """Given a list of proto_library targets, returns a list of all transitively
    -included .proto files and a list of options to pass to protoc (containing
    --proto_path locations).
    """
    srcs = []
    proto_paths = []
    for lib in protolibs:
        if ProtoInfo not in lib:
            fail("Expected a proto_library, got: ", lib)
        for src in lib[ProtoInfo].transitive_sources.to_list():
            if src.path.endswith(".proto"):
                srcs.append(src)
            proto_paths += lib[ProtoInfo].transitive_proto_path.to_list()
    srcs = collections.uniq(srcs)
    proto_path_opts = collections.before_each("--proto_path", collections.uniq(proto_paths))
    return srcs, proto_path_opts

def _proto_extract_kzip_impl(ctx):
    srcs, pathopt = get_proto_files_and_proto_path_opts(ctx.attr.srcs)

    extract(
        srcs = srcs,
        ctx = ctx,
        extractor = ctx.executable.extractor,
        kzip = ctx.outputs.kzip,
        mnemonic = "ProtoExtractKZip",
        opts = ["--"] + ctx.attr.opts + pathopt,
        vnames_config = ctx.file.vnames_config,
        deps = ctx.files.deps + srcs,
    )
    return [KytheVerifierSources(files = depset(srcs))]

proto_extract_kzip = rule(
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_empty = False,
            allow_files = True,
        ),
        "deps": attr.label_list(allow_files = True),
        "extractor": attr.label(
            default = Label("//kythe/cxx/extractor/proto:proto_extractor"),
            executable = True,
            cfg = "host",
        ),
        "opts": attr.string_list(),
        "vnames_config": attr.label(
            default = Label("//external:vnames_config"),
            allow_single_file = True,
        ),
    },
    outputs = {"kzip": "%{name}.kzip"},
    implementation = _proto_extract_kzip_impl,
)

def proto_verifier_test(
        name,
        srcs,
        deps = [],
        size = "small",
        tags = [],
        extractor = None,
        extractor_opts = [],
        indexer_opts = [],
        verifier_opts = [],
        convert_marked_source = False,
        vnames_config = None,
        visibility = None):
    """Extract, analyze, and verify a proto compilation.

    Args:
      name: Name of the test
      srcs: The compilation's source file inputs; each file's verifier goals will be checked
      deps: Optional list of proto_verifier_test targets to be used as proto compilation dependencies
      size: Test size
      tags: Test tags
      extractor: Executable extractor tool to invoke (defaults to protoc_extractor)
      extractor_opts: List of options passed to the extractor tool
      indexer_opts: List of options passed to the indexer tool
      verifier_opts: List of options passed to the verifier tool
      convert_marked_source: Whether the verifier should convert marked source.
      vnames_config: Optional path to a VName configuration file
      visibility: Visibility of underlying build targets
    Returns:
      Name of the test rule
    """
    kzip = _invoke(
        proto_extract_kzip,
        name = name + "_kzip",
        testonly = True,
        srcs = srcs,
        extractor = extractor,
        opts = extractor_opts,
        tags = tags,
        visibility = visibility,
        vnames_config = vnames_config,
        deps = deps,
    )
    entries = _invoke(
        index_compilation,
        name = name + "_entries",
        testonly = True,
        indexer = "//kythe/cxx/indexer/proto:indexer",
        opts = indexer_opts + ["--index_file"],
        tags = tags,
        visibility = visibility,
        deps = [kzip],
    )
    vopts = ["--ignore_dups"] + verifier_opts
    if convert_marked_source:
        vopts += ["--convert_marked_source"]
    return _invoke(
        verifier_test,
        name = name,
        size = size,
        srcs = [entries],
        opts = vopts,
        tags = tags,
        visibility = visibility,
        deps = [entries],
    )
