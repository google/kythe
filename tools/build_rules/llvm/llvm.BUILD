package(
    default_visibility = ["//visibility:public"],
)

TARGET_DEFAULTS = {
    "LLVMSupport": {
        "linkopts": [
            "-pthread",
            "-lncurses",
            "-ldl",
        ],
        "deps": [
            "@//external:zlib",
        ],
    },
    "LLVMDebugInfoCodeView": {
        "deps": [
            ":LLVMBinaryFormat",
        ],
    },
    "LLVMCore": {
        "additional_header_dirs": [
            # Layering violation.
            "/include/llvm/Analysis",
        ],
        "hdrs": glob([
            "include/llvm/*.h",
        ]),
    },
    "LLVMTransformUtils": {
        "hdrs": glob(["include/llvm-c/Transforms/**/*.h"]),
    },
    "LLVMScalarOpts": {
        "deps": [":LLVMTarget"],
    },
    "LLVMX86CodeGen": {
        "deps": [":LLVMipo"],
    },
    "clangAST": {
        "textual_hdrs": [
            ":tools_clang_include_clang_AST_genhdrs",
        ],
    },
    "clangBasic": {
        "deps": [
            ":LLVMTarget",
        ],
        "textual_hdrs": [
            "tools/clang/include/clang/Basic/Version.inc",
            ":tools_clang_include_clang_Basic_genhdrs",
        ],
    },
    "clangCodeGen": {
        "deps": [":all_targets"],
    },
    "clangDriver": {
        "textual_hdrs": [
            ":tools_clang_include_clang_StaticAnalyzer_Checkers_genhdrs",
        ],
    },
    "clangFrontend": {
        "deps": [
            ":LLVMLinker",
        ],
        "hdrs": [
            "tools/clang/include/clang/StaticAnalyzer/Core/AnalyzerOptions.h",
        ],
        "textual_hdrs": [
            "tools/clang/include/clang/StaticAnalyzer/Core/Analyses.def",
            "tools/clang/include/clang/StaticAnalyzer/Core/AnalyzerOptions.def",
        ],
    },
    "clangIndex": {
        "textual_hdrs": [
            "tools/clang/include/clang/StaticAnalyzer/Core/Analyses.def",
        ],
    },
    "clangParse": {
        "textual_hdrs": [
            ":tools_clang_include_clang_Parse_genhdrs",
        ],
    },
    "clangSema": {
        "textual_hdrs": [
            ":tools_clang_include_clang_Sema_genhdrs",
        ],
        "deps": [":all_targets"],
    },
    "clangSerialization": {
        "textual_hdrs": [
            ":tools_clang_include_clang_Serialization_genhdrs",
        ],
    },
    "ClangDriverOptions": {
        "textual_hdrs": [
            "tools/clang/include/clang/Frontend/LangStandards.def",
        ],
    },
}

cc_library(
    name = "clang-c",
    hdrs = glob(["tools/clang/include/clang-c/*.h"]) + [
        "tools/clang/include/clang/Config/config.h",
    ],
    includes = [
        "tools/clang/include",
    ],
)

cc_library(
    name = "llvm-c",
    hdrs = glob([
        "include/llvm-c/*.h",
    ]) + [
        "include/llvm/Config/abi-breaking.h",
        "include/llvm/Config/config.h",
        "include/llvm/Config/llvm-config.h",
        "include/llvm/Config/AsmParsers.def",
        "include/llvm/Config/AsmPrinters.def",
        "include/llvm/Config/Disassemblers.def",
        "include/llvm/Config/Targets.def",
    ],
    includes = [
        "include",
    ],
)

genrule(
    name = "llvm_vcsrevision_h_gen",
    outs = ["include/llvm/Support/VCSRevision.h"],
    cmd = "touch $@",
)

cc_library(
    name = "llvm_vcsrevision_h",
    hdrs = ["include/llvm/Support/VCSRevision.h"],
)

load("@io_kythe//tools:build_rules/cc_resources.bzl", "cc_resources")

builtin_headers = glob(
    ["tools/clang/lib/Headers/**"],
    exclude = ["tools/clang/lib/Headers/**/CMakeLists.txt"],
) + [
    "tools/clang/lib/Headers/arm_fp16.h",
    "tools/clang/lib/Headers/arm_neon.h",
]

genrule(
    name = "builtin_headers_gen",
    srcs = builtin_headers,
    outs = [hdr.replace("lib/Headers/", "staging/include/") for hdr in builtin_headers],
    cmd = """
      SRCS=($(SRCS))
      OUTS=($(OUTS))
      for i in "$${!SRCS[@]}"; do
        cp $${SRCS[$$i]} $${OUTS[$$i]}
      done""",
    output_to_bindir = True,
)

cc_resources(
    name = "clang_builtin_headers_resources",
    data = [":builtin_headers_gen"],
)

load("@io_kythe//tools/build_rules/llvm:cmake_defines.bzl", "cmake_defines", "LLVM_TARGETS")
load("@io_kythe//tools/build_rules/llvm:generated_llvm_build_deps.bzl", "LLVM_BUILD_DEPS")
load("@io_kythe//tools/build_rules/llvm:llvm.bzl", "make_context")
load("@io_kythe//tools/build_rules/llvm:generated_cmake_targets.bzl", "generated_cmake_targets")

cc_library(
    name = "all_targets",
    deps = [":LLVM%sCodeGen" % t for t in LLVM_TARGETS],
)

generated_cmake_targets(make_context(
    cmake_defines = cmake_defines(),
    llvm_build_deps = LLVM_BUILD_DEPS,
    target_defaults = TARGET_DEFAULTS,
))
