workspace(
    name = "io_kythe",
    managed_directories = {"@npm": ["node_modules"]},
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("//:version.bzl", "MAX_VERSION", "MIN_VERSION", "check_version")

# Check that the user has a version between our minimum supported version of
# Bazel and our maximum supported version of Bazel.
check_version(MIN_VERSION, MAX_VERSION)

http_archive(
    name = "bazel_toolchains",
    sha256 = "1adf5db506a7e3c465a26988514cfc3971af6d5b3c2218925cd6e71ee443fc3f",
    strip_prefix = "bazel-toolchains-4.0.0",
    urls = [
        "https://github.com/bazelbuild/bazel-toolchains/releases/download/4.0.0/bazel-toolchains-4.0.0.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-toolchains/releases/download/4.0.0/bazel-toolchains-4.0.0.tar.gz",
    ],
)

load("//:setup.bzl", "kythe_rule_repositories", "maybe")

kythe_rule_repositories()

# TODO(schroederc): remove this.  This needs to be loaded before loading the
# go_* rules.  Normally, this is done by go_rules_dependencies in external.bzl,
# but because we want to overload some of those dependencies, we need the go_*
# rules before go_rules_dependencies.  Likewise, we can't precisely control
# when loads occur within a Starlark file so we now need to load this
# manually...
load("@io_bazel_rules_go//go/private:repositories.bzl", "go_name_hack")

maybe(
    go_name_hack,
    name = "io_bazel_rules_go_name_hack",
    is_rules_go = False,
)

# gazelle:repository_macro external.bzl%_go_dependencies
load("//:external.bzl", "kythe_dependencies")

kythe_dependencies()

load("//tools/build_rules/external_tools:external_tools_configure.bzl", "external_tools_configure")

external_tools_configure()

load("@npm//@bazel/labs:package.bzl", "npm_bazel_labs_dependencies")

npm_bazel_labs_dependencies()

load("@maven//:compat.bzl", "compat_repositories")

compat_repositories()

# If the configuration here changes, run tools/platforms/configs/rebuild.sh
load("@bazel_toolchains//rules:environments.bzl", "clang_env")
load("@bazel_toolchains//rules:rbe_repo.bzl", "rbe_autoconfig")
load("//tools/platforms:toolchain_config_suite_spec.bzl", "DEFAULT_TOOLCHAIN_CONFIG_SUITE_SPEC")

rbe_autoconfig(
    name = "rbe_default",
    env = clang_env(),
    export_configs = True,
    toolchain_config_spec_name = "rbe_default_config_spec",
    toolchain_config_suite_spec = DEFAULT_TOOLCHAIN_CONFIG_SUITE_SPEC,
    use_legacy_platform_definition = False,
)

rbe_autoconfig(
    name = "rbe_bazel_minversion",
    bazel_version = MIN_VERSION,
    env = clang_env(),
    export_configs = True,
    toolchain_config_spec_name = "rbe_bazel_minversion_config_spec",
    toolchain_config_suite_spec = DEFAULT_TOOLCHAIN_CONFIG_SUITE_SPEC,
    use_legacy_platform_definition = False,
)

rbe_autoconfig(
    name = "rbe_bazel_maxversion",
    bazel_version = MAX_VERSION,
    env = clang_env(),
    export_configs = True,
    toolchain_config_spec_name = "rbe_bazel_maxversion_config_spec",
    toolchain_config_suite_spec = DEFAULT_TOOLCHAIN_CONFIG_SUITE_SPEC,
    use_legacy_platform_definition = False,
)
