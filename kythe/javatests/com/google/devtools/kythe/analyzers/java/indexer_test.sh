#!/bin/bash
set -eo pipefail
# Copyright 2015 The Kythe Authors. All rights reserved.
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
# This script tests the java indexer's CLI.

: ${indexer?:missing indexer}
: ${indexpack?:missing indexpack}
: ${entrystream?:missing entrystream}
test_kindex="$PWD/kythe/testdata/test.kindex"

# Test indexing a .kindex file
$indexer $test_kindex | $entrystream >/dev/null

tmp="$(mktemp -d 2>/dev/null || mktemp -d -t 'kythetest')"
trap 'rm -rf "$tmp"' EXIT ERR INT

UNIT="$($indexpack --verbose --to_archive "$tmp/archive" "$test_kindex" | \
  awk '-F\t' '/Wrote compilation:/ { print $2 }')"

# Test indexing compilation in an indexpack
$indexer --index_pack="$tmp/archive" $UNIT | $entrystream >/dev/null
