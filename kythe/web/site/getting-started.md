---
layout: page
title: Getting Started
permalink: /getting-started/
---

* toc
{:toc}

## Get the Kythe source code

Decide where you want to store kythe code, e.g. `~/my/code/dir` (note that
after we clone from git, it will append 'kythe' as the last directory).

{% highlight bash %}
cd ~/my/code/dir
git clone https://github.com/kythe/kythe.git
{% endhighlight %}

Also set the env var `KYTHE_DIR=~/my/code/dir/kythe` in your `.bashrc`
while you're at it.

If you use ssh to authenticate to github:
{% highlight bash %}
git clone git@github.com:kythe/kythe.git
{% endhighlight %}

### External Dependencies

Kythe relies on the following external dependencies:

* asciidoc
* bison-3.0.4 (2.3 is also acceptable)
* clang >= 3.6
* cmake >= 3.4.3
* [docker](https://www.docker.com/) (for release images `//kythe/release/...` and `//buildtools/docker`)
* flex-2.5
* go >= 1.7
* graphviz
* jdk >= 8
* [leiningen](http://leiningen.org/) (used to build `//kythe/web/ui`)
* libcurl4-openssl-dev
* libncurses-dev
* libssl-dev
* node.js
* parallel
* source-highlight
* uuid-dev
* wget

You will need to ensure these packages are installed on the system where you
intend to build Kythe. There are instructions for using `apt-get` below.
If you are using macOS, see [Instructions for macOS]({{site.baseuri}}/getting-started-macos).


#### Installing Debian Jessie Packages

{% highlight bash %}
echo deb http://http.debian.net/debian jessie-backports main >> /etc/apt/sources.list
apt-get update

apt-get install \
    asciidoc asciidoctor source-highlight graphviz \
    gcc libssl-dev uuid-dev libncurses-dev libcurl4-openssl-dev flex clang-3.5 bison \
    openjdk-8-jdk \
    parallel \
    wget

# https://golang.org/dl/ for Golang installation
# https://docs.docker.com/installation/debian/#debian-jessie-80-64-bit for Docker installation
{% endhighlight %}

### Internal Dependencies

All other Kythe dependencies are hosted within the repository under
`//third_party/...`. Run the `./tools/modules/update.sh` script to update these
dependencies to the exact revision that we test against.

This step may take a little time the first time it is run and should be quick
on subsequent runs.

#### Troubleshooting bazel/clang/llvm errors
You must either have `/usr/bin/clang` aliased properly, or the `CC` env var set
for Bazel:

{% highlight bash %}
echo 'build --client_env=CC=/usr/bin/clang-3.6' >>~/.bazelrc
{% endhighlight %}

OR:

{% highlight bash %}
sudo ln -s /usr/bin/clang-3.6 /usr/bin/clang
sudo ln -s /usr/bin/clang++-3.6 /usr/bin/clang++
{% endhighlight %}

OR:

{% highlight bash %}
echo 'export CC=/usr/bin/clang' >> ~/.bashrc
source ~/.bashrc
{% endhighlight %}

If you ran bazel and get errors like this:

{% highlight bash %}
/home/username/kythe/third_party/zlib/BUILD:10:1: undeclared inclusion(s) in rule '//third_party/zlib:zlib':
this rule is missing dependency declarations for the following files included by 'third_party/zlib/uncompr.c':
  '/usr/lib/llvm-3.6/lib/clang/3.6.0/include/limits.h'
  '/usr/lib/llvm-3.6/lib/clang/3.6.0/include/stddef.h'
  '/usr/lib/llvm-3.6/lib/clang/3.6.0/include/stdarg.h'.
{% endhighlight %}

then you need to clean and rebuild your TOOLCHAIN:

{% highlight bash %}
bazel clean --expunge && bazel build @local_config_cc//:toolchain
{% endhighlight %}

## Building Kythe

### Building using Bazel

Kythe uses [Bazel](http://bazel.io) to build its source code.  After
[installing Bazel](http://bazel.io/docs/install.html) and all external
dependencies, building Kythe should be as simple as:

{% highlight bash %}
./tools/modules/update.sh  # Ensure third_party is updated

bazel build //... # Build all Kythe sources
bazel test  //... # Run all Kythe tests
{% endhighlight %}

Please note that you must use a non-jdk7 version of Bazel. Some package managers
may provide the jdk7 version by default. To determine if you are using an
incompatible version of Bazel, look for `jdk7` in the build label that
is printed by `bazel version`.

Also note that not all targets build with `//...` - some targets are
purposefully omitted.  This includes `//kythe/web/ui`, `//kythe/release`, and
many of the docker images we push.

### Build a release of Kythe using Bazel and unpack it in /opt/kythe

Many examples on the site assume you have installed kythe in /opt/kythe.

{% highlight bash %}
# Build a Kythe release
bazel build //kythe/release
# Set current Kythe version
# check bazel-genfiles/kythe/release/ directory to get current version.
export KYTHE_RELEASE="x.y.z"
# Extract our new Kythe release to /opt/ including its version number
tar -zxf bazel-genfiles/kythe/release/kythe-v${KYTHE_RELEASE}.tar.gz --directory /opt/
# Remove the old pointer to Kythe if we had one
rm -f /opt/kythe
# Point Kythe to our new version
ln -s /opt/kythe-v${KYTHE_RELEASE} /opt/kythe
{% endhighlight %}

### Using the Go tool to build Go sources directly

Kythe's Go sources can be directly built with the `go` tool as well as with
Bazel.

{% highlight bash %}
# Install LevelDB/snappy libraries for https://github.com/jmhodges/levigo
sudo apt-get install libleveldb-dev libsnappy-dev

# With an appropriate GOPATH setup
go get kythe.io/kythe/...

# Using the vendored versions of the needed third_party Go libraries
git clone https://github.com/kythe/kythe.git
GOPATH=$GOPATH:$PWD/kythe/third_party/go go get kythe.io/kythe/...
{% endhighlight %}

The additional benefits of using Bazel are the built-in support for generating
the Go protobuf code in `kythe/proto/` and the automatic usage of the checked-in
`third_party/go` libraries (instead of adding to your `GOPATH`).  However, for
quick access to Kythe's Go sources (which implement most of Kythe's platform and
language-agnostic services), using the Go tool is very convenient.

## Updating and building the website

* Make change in ./kythe/web/site
* Spell check
* Build a local version to verify fixes

Prerequisites:
{% highlight bash %}
apt-get install ruby ruby-dev build-essential
gem install bundler
{% endhighlight %}

Build and serve:
{% highlight bash %}
cd ./kythe/web/site
./build.sh
# Serve website locally on port 4000
bundle exec jekyll serve
{% endhighlight %}
