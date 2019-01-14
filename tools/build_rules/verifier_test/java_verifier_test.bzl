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

load(
    ":verifier_test.bzl",
    "KytheVerifierSources",
    "extract",
    "index_compilation",
    "verifier_test",
)

KytheJavaJar = provider(
    doc = "Bundled Java compilation.",
    fields = {
        "jar": "The bundled jar file.",
    },
)

def _invoke(rulefn, name, **kwargs):
    """Invoke rulefn with name and kwargs, returning the label of the rule."""
    rulefn(name = name, **kwargs)
    return "//{}:{}".format(native.package_name(), name)

def _java_extract_kzip_impl(ctx):
    jars = []
    for dep in ctx.attr.deps:
        jars += [dep[KytheJavaJar].jar]

    # Actually compile the sources to be used as a dependency for other tests
    jar = ctx.actions.declare_file(ctx.outputs.kzip.basename + ".jar", sibling = ctx.outputs.kzip)
    info = java_common.compile(
        ctx,
        javac_opts = java_common.default_javac_opts(
            ctx,
            java_toolchain_attr = "_java_toolchain",
        ) + ctx.attr.opts,
        java_toolchain = ctx.attr._java_toolchain,
        host_javabase = ctx.attr._host_javabase,
        source_files = ctx.files.srcs,
        output = jar,
        deps = [JavaInfo(compile_jar = curr_jar, output_jar = curr_jar) for curr_jar in jars],
    )

    args = ctx.attr.opts + [
        "-encoding",
        "utf-8",
        "-cp",
        ":".join([j.path for j in jars]),
    ]
    for src in ctx.files.srcs:
        args += [src.short_path]
    extract(
        srcs = ctx.files.srcs,
        ctx = ctx,
        extractor = ctx.executable.extractor,
        kzip = ctx.outputs.kzip,
        mnemonic = "JavaExtractKZip",
        opts = args,
        vnames_config = ctx.file.vnames_config,
        deps = jars + ctx.files.data,
    )
    return [
        KytheJavaJar(jar = jar),
        KytheVerifierSources(files = depset(ctx.files.srcs)),
    ]

java_extract_kzip = rule(
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_empty = False,
            allow_files = True,
        ),
        "data": attr.label_list(
            allow_files = True,
        ),
        "extractor": attr.label(
            default = Label("@io_kythe//kythe/java/com/google/devtools/kythe/extractors/java/standalone:javac_extractor"),
            executable = True,
            cfg = "host",
        ),
        "opts": attr.string_list(),
        "vnames_config": attr.label(
            default = Label("//external:vnames_config"),
            allow_single_file = True,
        ),
        "deps": attr.label_list(
            providers = [KytheJavaJar],
        ),
        "_host_javabase": attr.label(
            cfg = "host",
            default = Label("@bazel_tools//tools/jdk:current_java_runtime"),
        ),
        "_java_toolchain": attr.label(
            default = Label("@bazel_tools//tools/jdk:toolchain"),
        ),
    },
    fragments = ["java"],
    host_fragments = ["java"],
    outputs = {"kzip": "%{name}.kzip"},
    implementation = _java_extract_kzip_impl,
)

def java_verifier_test(
        name,
        srcs,
        meta = [],
        deps = [],
        size = "small",
        tags = [],
        extractor = None,
        extractor_opts = [
            "-source",
            "9",
            "-target",
            "9",
        ],
        indexer_opts = ["--verbose"],
        verifier_opts = ["--ignore_dups"],
        load_plugin = None,
        extra_goals = [],
        vnames_config = None,
        visibility = None):
    """Extract, analyze, and verify a Java compilation.

    Args:
      srcs: The compilation's source file inputs; each file's verifier goals will be checked
      deps: Optional list of java_verifier_test targets to be used as Java compilation dependencies
      meta: Optional list of Kythe metadata files
      extractor: Executable extractor tool to invoke (defaults to javac_extractor)
      extractor_opts: List of options passed to the extractor tool
      indexer_opts: List of options passed to the indexer tool
      verifier_opts: List of options passed to the verifier tool
      load_plugin: Optional Java analyzer plugin to load
      extra_goals: List of text files containing verifier goals additional to those in srcs
      vnames_config: Optional path to a VName configuration file
    """
    kzip = _invoke(
        java_extract_kzip,
        name = name + "_kzip",
        testonly = True,
        srcs = srcs,
        data = meta,
        extractor = extractor,
        opts = extractor_opts,
        tags = tags,
        visibility = visibility,
        vnames_config = vnames_config,
        # This is a hack to depend on the .jar producer.
        deps = [d + "_kzip" for d in deps],
    )
    indexer = "//kythe/java/com/google/devtools/kythe/analyzers/java:indexer"
    tools = []
    if load_plugin:
        # If loaded plugins have deps, those must be included in the loaded jar
        native.java_binary(
            name = name + "_load_plugin",
            main_class = "not.Used",
            runtime_deps = [load_plugin],
        )
        load_plugin_deploy_jar = ":{}_load_plugin_deploy.jar".format(name)
        indexer_opts = indexer_opts + [
            "--load_plugin",
            "$(location {})".format(load_plugin_deploy_jar),
        ]
        tools += [load_plugin_deploy_jar]

    entries = _invoke(
        index_compilation,
        name = name + "_entries",
        testonly = True,
        indexer = indexer,
        opts = indexer_opts,
        tags = tags,
        tools = tools,
        visibility = visibility,
        deps = [kzip],
    )
    return _invoke(
        verifier_test,
        name = name,
        size = size,
        srcs = [entries] + extra_goals,
        opts = verifier_opts,
        tags = tags,
        visibility = visibility,
        deps = [entries],
    )
