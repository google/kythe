# Generated file, do not modify by hand
# Generated by 'rbe_default' rbe_autoconfig rule
"""Definitions to be used in rbe_repo attr of an rbe_autoconf rule  """
toolchain_config_spec0 = struct(config_repos = [], create_cc_configs = True, create_java_configs = True, env = {"ABI_LIBC_VERSION": "glibc_2.19", "ABI_VERSION": "clang", "BAZEL_COMPILER": "clang", "BAZEL_HOST_SYSTEM": "i686-unknown-linux-gnu", "BAZEL_TARGET_CPU": "k8", "BAZEL_TARGET_LIBC": "glibc_2.19", "BAZEL_TARGET_SYSTEM": "x86_64-unknown-linux-gnu", "CC": "clang", "CC_TOOLCHAIN_NAME": "linux_gnu_x86"}, java_home = "/usr/lib/jvm/11.29.3-ca-jdk11.0.2/reduced", java_version = "11", name = "default_toolchain_config_spec_name")
_TOOLCHAIN_CONFIG_SPECS = [toolchain_config_spec0]
_BAZEL_TO_CONFIG_SPEC_NAMES = {"4.0.0": ["default_toolchain_config_spec_name"]}
LATEST = "sha256:affd4b4d05144cb524db0d4f4639db74942b53ff499a1dfce23ce0592654e259"
CONTAINER_TO_CONFIG_SPEC_NAMES = {"sha256:affd4b4d05144cb524db0d4f4639db74942b53ff499a1dfce23ce0592654e259": ["default_toolchain_config_spec_name"]}
_DEFAULT_TOOLCHAIN_CONFIG_SPEC = toolchain_config_spec0
TOOLCHAIN_CONFIG_AUTOGEN_SPEC = struct(
        bazel_to_config_spec_names_map = _BAZEL_TO_CONFIG_SPEC_NAMES,
        container_to_config_spec_names_map = CONTAINER_TO_CONFIG_SPEC_NAMES,
        default_toolchain_config_spec = _DEFAULT_TOOLCHAIN_CONFIG_SPEC,
        latest_container = LATEST,
        toolchain_config_specs = _TOOLCHAIN_CONFIG_SPECS,
    )