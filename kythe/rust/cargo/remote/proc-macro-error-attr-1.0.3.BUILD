"""
cargo-raze crate build file.

DO NOT EDIT! Replaced on runs of cargo-raze
"""
package(default_visibility = [
  # Public for visibility by "@raze__crate__version//" targets.
  #
  # Prefer access through "//kythe/rust/cargo", which limits external
  # visibility to explicit Cargo.toml dependencies.
  "//visibility:public",
])

licenses([
  "notice", # MIT from expression "MIT OR Apache-2.0"
])

load(
    "@io_bazel_rules_rust//rust:rust.bzl",
    "rust_library",
    "rust_binary",
    "rust_test",
)


# Unsupported target "build-script-build" with type "custom-build" omitted

rust_library(
    name = "proc_macro_error_attr",
    crate_type = "proc-macro",
    deps = [
        "@raze__proc_macro2__1_0_19//:proc_macro2",
        "@raze__quote__1_0_7//:quote",
        "@raze__syn__1_0_36//:syn",
        "@raze__syn_mid__0_5_0//:syn_mid",
    ],
    srcs = glob(["**/*.rs"]),
    crate_root = "src/lib.rs",
    edition = "2018",
    rustc_flags = [
        "--cap-lints=allow",
    ],
    version = "1.0.3",
    tags = ["cargo-raze"],
    crate_features = [
    ],
)

