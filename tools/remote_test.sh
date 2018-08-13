#!/bin/bash -e
#
# Copyright 2018 The Kythe Authors. All rights reserved.
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
#
# Tests the Kythe repository using experimental Remote Build Execution.

exclude_tags=(
  manual
)
exclude_kinds=(
  asciidoc
  go_entries
  go_extract
)
exclude_targets=(
  //kythe/cxx/verifier:all
  //third_party/llvm
  //third_party:libuuid
  @com_google_common_flogger//api:gen_platform_provider
)

query='//...'
for tag in ${exclude_tags[@]}; do
  query="$query - attr(tags, $tag, //...)"
done
for kind in ${exclude_kinds[@]}; do
  query="$query - kind($kind, //...)"
done
for target in ${exclude_targets[@]}; do
  query="$query - rdeps(//..., $target)"
done
targets=($(bazel query "$query"))

bazel test --config=remote "$@" "${targets[@]}"
