# Copyright 2020 The Kythe Authors. All rights reserved.
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

def cc_release_binary(name, **kwargs):
    """Define two cc_binary rules, one with the _static suffix which uses full_static_link."""
    native.cc_binary(name = name, **kwargs)
    native.cc_binary(
        name = name + "_static",
        features = kwargs.pop("features", []) + ["fully_static_link"],
        linkopts = kwargs.pop("linkopts", []) + ["-l:libstdc++.a"],
        **kwargs
    )
