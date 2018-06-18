workspace(name = "io_kythe")

load("//:version.bzl", "check_version")

# Check that the user has a version between our minimum supported version of
# Bazel and our maximum supported version of Bazel.
check_version("0.13", "0.14")

load("//tools/cpp:clang_configure.bzl", "clang_configure")

clang_configure()

bind(
    name = "libuuid",
    actual = "//third_party:libuuid",
)

new_http_archive(
    name = "org_libmemcached_libmemcached",
    build_file = "third_party/libmemcached.BUILD",
    sha256 = "e22c0bb032fde08f53de9ffbc5a128233041d9f33b5de022c0978a2149885f82",
    strip_prefix = "libmemcached-1.0.18",
    url = "https://launchpad.net/libmemcached/1.0/1.0.18/+download/libmemcached-1.0.18.tar.gz",
)

bind(
    name = "libmemcached",
    actual = "@org_libmemcached_libmemcached//:libmemcached",
)

bind(
    name = "guava",  # required by @com_google_protobuf
    actual = "//third_party/guava",
)

bind(
    name = "gson",  # required by @com_google_protobuf
    actual = "@com_google_code_gson_gson//jar",
)

new_http_archive(
    name = "net_zlib",
    build_file = "third_party/zlib.BUILD",
    sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
    strip_prefix = "zlib-1.2.11",
    urls = [
        "https://zlib.net/zlib-1.2.11.tar.gz",
    ],
)

bind(
    name = "zlib",  # required by @com_google_protobuf
    actual = "@net_zlib//:zlib",
)

http_archive(
    name = "boringssl",  # Must match upstream workspace name.
    # Gitiles creates gzip files with an embedded timestamp, so we cannot use
    # sha256 to validate the archives.  We must rely on the commit hash and https.
    # Commits must come from the master-with-bazel branch.
    url = "https://boringssl.googlesource.com/boringssl/+archive/4be3aa87917b20fedc45fa1fc5b6a2f3738612ad.tar.gz",
)

# Make sure to update regularly in accordance with Abseil's principle of live at HEAD
http_archive(
    name = "com_google_absl",
    strip_prefix = "abseil-cpp-da336a84e9c1f86409b21996164ae9602b37f9ca",
    url = "https://github.com/abseil/abseil-cpp/archive/da336a84e9c1f86409b21996164ae9602b37f9ca.zip",
)

http_archive(
    name = "com_google_googletest",
    sha256 = "89cebb92b9a7eb32c53e180ccc0db8f677c3e838883c5fbd07e6412d7e1f12c7",
    strip_prefix = "googletest-d175c8bf823e709d570772b038757fadf63bc632",
    url = "https://github.com/google/googletest/archive/d175c8bf823e709d570772b038757fadf63bc632.zip",
)

http_archive(
    name = "com_github_gflags_gflags",
    sha256 = "94ad0467a0de3331de86216cbc05636051be274bf2160f6e86f07345213ba45b",
    strip_prefix = "gflags-77592648e3f3be87d6c7123eb81cbad75f9aef5a",
    url = "https://github.com/gflags/gflags/archive/77592648e3f3be87d6c7123eb81cbad75f9aef5a.zip",
)

http_archive(
    name = "com_googlesource_code_re2",
    # Gitiles creates gzip files with an embedded timestamp, so we cannot use
    # sha256 to validate the archives.  We must rely on the commit hash and https.
    url = "https://code.googlesource.com/re2/+archive/2c220e7df3c10d42d74cb66290ec89116bb5e6be.tar.gz",
)

new_http_archive(
    name = "com_github_google_glog",
    build_file = "third_party/googlelog.BUILD",
    sha256 = "ce61883437240d650be724043e8b3c67e257690f876ca9fd53ace2a791cfea6c",
    strip_prefix = "glog-bac8811710c77ac3718be1c4801f43d37c1aea46",
    url = "https://github.com/google/glog/archive/bac8811710c77ac3718be1c4801f43d37c1aea46.zip",
)

new_http_archive(
    name = "com_github_tencent_rapidjson",
    build_file = "third_party/rapidjson.BUILD",
    sha256 = "8e00c38829d6785a2dfb951bb87c6974fa07dfe488aa5b25deec4b8bc0f6a3ab",
    strip_prefix = "rapidjson-1.1.0",
    url = "https://github.com/Tencent/rapidjson/archive/v1.1.0.zip",
)

new_http_archive(
    name = "com_github_stedolan_jq",
    build_file = "third_party/jq.BUILD",
    sha256 = "998c41babeb57b4304e65b4eb73094279b3ab1e63801b6b4bddd487ce009b39d",
    strip_prefix = "jq-1.4",
    url = "https://github.com/stedolan/jq/releases/download/jq-1.4/jq-1.4.tar.gz",
)

new_http_archive(
    name = "com_github_google_snappy",
    build_file = "third_party/snappy.BUILD",
    sha256 = "61e05a0295fd849072668b1f3494801237d809427cfe8fd014cda455036c3ef7",
    strip_prefix = "snappy-1.1.7",
    url = "https://github.com/google/snappy/archive/1.1.7.zip",
)

new_http_archive(
    name = "com_github_google_leveldb",
    build_file = "third_party/leveldb.BUILD",
    sha256 = "5b2bd7a91489095ad54bb81ca6544561025b48ec6d19cc955325f96755d88414",
    strip_prefix = "leveldb-1.20",
    url = "https://github.com/google/leveldb/archive/v1.20.zip",
)

maven_jar(
    name = "com_google_code_gson_gson",
    artifact = "com.google.code.gson:gson:2.8.5",
    sha1 = "f645ed69d595b24d4cf8b3fbb64cc505bede8829",
)

maven_jar(
    name = "com_google_guava_guava",
    artifact = "com.google.guava:guava:25.1-jre",
    sha1 = "6c57e4b22b44e89e548b5c9f70f0c45fe10fb0b4",
)

maven_jar(
    name = "junit_junit",
    artifact = "junit:junit:4.12",
    sha1 = "2973d150c0dc1fefe998f834810d68f278ea58ec",
)

maven_jar(
    name = "com_google_re2j_re2j",
    artifact = "com.google.re2j:re2j:1.2",
    sha1 = "4361eed4abe6f84d982cbb26749825f285996dd2",
)

maven_jar(
    name = "com_beust_jcommander",
    artifact = "com.beust:jcommander:1.48",
    sha1 = "bfcb96281ea3b59d626704f74bc6d625ff51cbce",
)

maven_jar(
    name = "com_google_truth_truth",
    artifact = "com.google.truth:truth:0.41",
    sha1 = "846cd094934911f635ba2dadc016d538b8c30927",
)

maven_jar(
    name = "com_googlecode_java_diff_utils",
    artifact = "com.googlecode.java-diff-utils:diffutils:1.3.0",
    sha1 = "7e060dd5b19431e6d198e91ff670644372f60fbd",
)

maven_jar(
    name = "com_google_code_findbugs_jsr305",
    artifact = "com.google.code.findbugs:jsr305:3.0.1",
    sha1 = "f7be08ec23c21485b9b5a1cf1654c2ec8c58168d",
)

maven_jar(
    name = "com_google_auto_value_auto_value",
    artifact = "com.google.auto.value:auto-value:1.5.4",
    sha1 = "65183ddd1e9542d69d8f613fdae91540d04e1476",
)

maven_jar(
    name = "com_google_auto_service_auto_service",
    artifact = "com.google.auto.service:auto-service:1.0-rc4",
    sha1 = "44954d465f3b9065388bbd2fc08a3eb8fd07917c",
)

maven_jar(
    name = "com_google_auto_auto_common",
    artifact = "com.google.auto:auto-common:0.10",
    sha1 = "c8f153ebe04a17183480ab4016098055fb474364",
)

maven_jar(
    name = "javax_annotation_jsr250_api",
    artifact = "javax.annotation:jsr250-api:1.0",
    sha1 = "5025422767732a1ab45d93abfea846513d742dcf",
)

maven_jar(
    name = "com_google_common_html_types",
    artifact = "com.google.common.html.types:types:1.0.8",
    sha1 = "9e9cf7bc4b2a60efeb5f5581fe46d17c068e0777",
)

maven_jar(
    name = "org_ow2_asm_asm",
    artifact = "org.ow2.asm:asm:6.0",
    sha1 = "bc6fa6b19424bb9592fe43bbc20178f92d403105",
)

maven_jar(
    name = "com_google_errorprone_error_prone_annotations",
    artifact = "com.google.errorprone:error_prone_annotations:2.3.1",
    sha1 = "a6a2b2df72fd13ec466216049b303f206bd66c5d",
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "8b68d0630d63d95dacc0016c3bb4b76154fe34fca93efd65d1c366de3fcb4294",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.12.1/rules_go-0.12.1.tar.gz",
)

# proto_library, cc_proto_library, and java_proto_library rules implicitly
# depend on @com_google_protobuf for protoc and proto runtimes.
#
# N.B. We have a near-clone of the protobuf BUILD file overriding upstream so
# that we can set the unexported config variable to enable zlib. Without this,
# protobuf silently yields link errors.
new_http_archive(
    name = "com_google_protobuf",
    build_file = "third_party/protobuf.BUILD",
    sha256 = "091d4263d9a55eccb6d3c8abde55c26eaaa933dea9ecabb185cdf3795f9b5ca2",
    strip_prefix = "protobuf-3.5.1.1",
    urls = ["https://github.com/google/protobuf/archive/v3.5.1.1.zip"],
)

# This is required by the proto_library implementation for its
# :cc_toolchain rule.
http_archive(
    name = "com_google_protobuf_cc",
    strip_prefix = "protobuf-106ffc04be1abf3ff3399f54ccf149815b287dd9",
    urls = ["https://github.com/google/protobuf/archive/106ffc04be1abf3ff3399f54ccf149815b287dd9.zip"],
)

# This is required by the proto_library implementation for its
# :java_toolchain rule.
http_archive(
    name = "com_google_protobuf_java",
    strip_prefix = "protobuf-106ffc04be1abf3ff3399f54ccf149815b287dd9",
    urls = ["https://github.com/google/protobuf/archive/106ffc04be1abf3ff3399f54ccf149815b287dd9.zip"],
)

http_archive(
    name = "google_bazel_common",
    strip_prefix = "bazel-common-370b397507d9bab9d9cdad8dfe7e6ccc8c2d0c67",
    urls = ["https://github.com/google/bazel-common/archive/370b397507d9bab9d9cdad8dfe7e6ccc8c2d0c67.zip"],
)

git_repository(
    name = "com_google_common_flogger",
    commit = "b08ed99eb6dcd62afe81fd0fafd97299b1870fbf",
    remote = "https://github.com/google/flogger",
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "ddedc7aaeb61f2654d7d7d4fd7940052ea992ccdb031b8f9797ed143ac7e8d43",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.12.0/bazel-gazelle-0.12.0.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

# Go imports are managed by Gazelle:
# https://github.com/bazelbuild/bazel-gazelle.  Use the
# `gazelle update-repos <go_import_path>` command to add or update any
# third_party Go library.
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("//tools:build_rules/shims.bzl", "go_repository")

go_repository(
    name = "com_github_golang_protobuf",
    commit = "b4deda0973fb4c70b50d226b1af49f3da59f5265",
    custom = "protobuf",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_jmhodges_levigo",
    commit = "c42d9e0ca023e2198120196f842701bb4c55d7b9",
    custom = "levigo",
    importpath = "github.com/jmhodges/levigo",
)

go_repository(
    name = "com_github_google_go_cmp",
    commit = "5411ab924f9ffa6566244a9e504bc347edacffd3",
    custom = "cmp",
    importpath = "github.com/google/go-cmp",
)

go_repository(
    name = "org_golang_x_sync",
    commit = "1d60e4601c6fd243af51cc01ddf169918a5407ca",
    custom = "sync",
    custom_git = "https://github.com/golang/sync.git",
    importpath = "golang.org/x/sync",
)

go_repository(
    name = "com_github_sourcegraph_jsonrpc2",
    commit = "a3d86c792f0f5a0c0c2c4ed9157125e914cb5534",
    custom = "jsonrpc2",
    importpath = "github.com/sourcegraph/jsonrpc2",
)

go_repository(
    name = "com_github_golang_snappy",
    commit = "553a641470496b2327abcac10b36396bd98e45c9",
    custom = "snappy",
    importpath = "github.com/golang/snappy",
)

go_repository(
    name = "com_github_sourcegraph_go_langserver",
    commit = "e526744fd766a8f42e55bd92a3843c2afcdbf08c",
    custom = "langserver",
    importpath = "github.com/sourcegraph/go-langserver",
)

go_repository(
    name = "com_github_pborman_uuid",
    commit = "c65b2f87fee37d1c7854c9164a450713c28d50cd",
    custom = "uuid",
    importpath = "github.com/pborman/uuid",
)

go_repository(
    name = "com_github_sergi_go_diff",
    commit = "da645544ed44df016359bd4c0e3dc60ee3a0da43",
    custom = "diff",
    importpath = "github.com/sergi/go-diff",
)

go_repository(
    name = "com_github_google_subcommands",
    commit = "a3682377147edf596d303faabd89f81977b3f678",
    custom = "subcommands",
    importpath = "github.com/google/subcommands",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "48418e5732e1b1e2a10207c8007a5f959e422f20",
    custom = "x_tools",
    custom_git = "https://github.com/golang/tools.git",
    importpath = "golang.org/x/tools",
)

go_repository(
    name = "org_golang_x_text",
    commit = "7922cc490dd5a7dbaa7fd5d6196b49db59ac042f",
    custom = "x_text",
    custom_git = "https://github.com/golang/text.git",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "org_golang_x_net",
    commit = "f73e4c9ed3b7ebdd5f699a16a880c2b1994e50dd",
    custom = "x_net",
    custom_git = "https://github.com/golang/net.git",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "816c9085562cd7ee03e7f8188a1cfd942858cded",
    custom = "errors",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "org_bitbucket_creachadair_stringset",
    commit = "e974a3c1694da0d5a14216ce46dbceef6a680978",
    custom = "stringset",
    custom_git = "https://bitbucket.org/creachadair/stringset.git",
    importpath = "bitbucket.org/creachadair/stringset",
)

go_repository(
    name = "org_bitbucket_creachadair_shell",
    commit = "3dcd505a7ca5845388111724cc2e094581e92cc6",
    custom = "shell",
    custom_git = "https://bitbucket.org/creachadair/shell.git",
    importpath = "bitbucket.org/creachadair/shell",
)

go_repository(
    name = "com_github_google_go_github",
    commit = "8ea2e2657df890db8fb434a9274799d641bd698c",
    custom = "github",
    importpath = "github.com/google/go-github",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "d07538b1475ec5b0ac85319e4a6706b2d2d8cab7",
    custom = "grpc",
    custom_git = "https://github.com/grpc/grpc-go.git",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "cdc340f7c179dbbfa4afd43b7614e8fcadde4269",
    custom = "x_oauth2",
    custom_git = "https://github.com/golang/oauth2.git",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "com_github_google_go_querystring",
    commit = "53e6ce116135b80d037921a7fdd5138cf32d7a8a",
    custom = "querystring",
    importpath = "github.com/google/go-querystring",
)

go_repository(
    name = "com_github_apache_beam",
    commit = "0ea97a562f82e98ff5cbe5a0825d298663112cdb",
    custom = "beam",
    importpath = "github.com/apache/beam",
)

go_repository(
    name = "org_golang_google_api",
    commit = "3097bf831ede4a24e08a3316254e29ca726383e3",
    custom = "google_api",
    custom_git = "https://github.com/google/google-api-go-client.git",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "01301d1df8060594708d76bda9062b0205b77e8b",
    custom = "google_cloud",
    custom_git = "https://github.com/GoogleCloudPlatform/google-cloud-go.git",
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "io_opencensus_go",
    commit = "c40611a83b49d279ee5203c85e4fe169dcb158b6",
    custom = "opencensus",
    custom_git = "https://github.com/census-instrumentation/opencensus-go.git",
    importpath = "go.opencensus.io",
)

go_repository(
    name = "com_github_syndtr_goleveldb",
    commit = "5d6fca44a948d2be89a9702de7717f0168403d3d",
    importpath = "github.com/syndtr/goleveldb",
)
