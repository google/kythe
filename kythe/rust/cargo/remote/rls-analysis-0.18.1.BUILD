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
  "notice", # Apache-2.0 from expression "Apache-2.0 OR MIT"
])

load(
    "@io_bazel_rules_rust//rust:rust.bzl",
    "rust_library",
    "rust_binary",
    "rust_test",
)


# Unsupported target "print-crate-id" with type "example" omitted

rust_library(
    name = "rls_analysis",
    crate_type = "lib",
    deps = [
        "@raze__fst__0_3_5//:fst",
        "@raze__itertools__0_8_2//:itertools",
        "@raze__json__0_11_15//:json",
        "@raze__log__0_4_11//:log",
        "@raze__rls_data__0_19_0//:rls_data",
        "@raze__rls_span__0_5_2//:rls_span",
        "@raze__serde__1_0_114//:serde",
        "@raze__serde_json__1_0_57//:serde_json",
    ],
    srcs = glob(["**/*.rs"]),
    crate_root = "src/lib.rs",
    edition = "2018",
    proc_macro_deps = [
        "@raze__derive_new__0_5_8//:derive_new",
    ],
    rustc_flags = [
        "--cap-lints=allow",
    ],
    version = "0.18.1",
    tags = ["cargo-raze"],
    crate_features = [
        "default",
    ],
)

# Unsupported target "std_api_crate" with type "bench" omitted
