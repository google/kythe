load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def maybe(repo_rule, name, **kwargs):
    """Defines a repository if it does not already exist.
    """
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def kythe_rule_repositories():
    """Defines external repositories for Kythe Bazel rules.

    These repositories must be loaded before calling external.bzl%kythe_dependencies.
    """
    maybe(
        http_archive,
        name = "io_bazel_rules_go",
        url = "https://github.com/bazelbuild/rules_go/releases/download/0.18.5/rules_go-0.18.5.tar.gz",
        sha256 = "a82a352bffae6bee4e95f68a8d80a70e87f42c4741e6a448bec11998fcc82329",
    )

    maybe(
        http_archive,
        name = "rules_java",
        url = "https://github.com/bazelbuild/rules_java/archive/973a93dd2d698929264d1028836f6b9cc60ff817.zip",
        sha256 = "a6cb0dbe343b67c7d4f3f11a68e327acdfc71fee1e17affa4e605129fc56bb15",
        strip_prefix = "rules_java-973a93dd2d698929264d1028836f6b9cc60ff817",
    )

    maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "e4fe70af52135d2ee592a07f916e6e1fc7c94cf8786c15e8c0d0f08b1fe5ea16",
        strip_prefix = "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
        url = "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.zip",
    )

    maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
        urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
    )

    maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "e04a82a72146bfbca2d0575947daa60fda1878c8d3a3afe868a8ec39a6b968bb",
        urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/0.31.1/rules_nodejs-0.31.1.tar.gz"],
    )
