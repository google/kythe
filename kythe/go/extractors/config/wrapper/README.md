# Kythe Extractor Wrapping

This package contains resources for preprocessing a repository's build
configuration before running Kythe Extraction via `extractrepo.go`.

## Gradle

`gradle_build_modifier.go` takes an input `gradle.build` file and appends the
bits necessary for replacing javac calls with Kythe's `javac-wrapper.sh`.

```
allprojects {
  gradle.projectsEvaluated {
    tasks.withType(JavaCompile) {
      options.fork = true
      options.forkOptions.executable = '/opt/kythe/extractors/javac-wrapper.sh'
    }
  }
}
```

If the input file already contains reference to
`options.forkOptions.executable`, then `gradle_build_modifier.go` does nothing.

### Future work

The current implementation uses simple string-based matching, without actually
understanding the structure.  If that becomes necessary in the future, it might
be better to use the existing Java libraries for `org.codehaus.groovy.ast` to
properly parse the build.gradle file and have more precise picture.  In
particular `org.codehaus.groovy.ast.CodeVisitorSupport` might be sufficient.

## Maven

Maven is handled in Java by
`com.google.devtools.kythe.platform.tools.MvnPomPreprocessor`.  It utilizes an
xml library to parse and modify the mvn `pom.xml` config file in a similar way
as described above for gradle.

### Future work

In theory if we can find a nice xml library for golang that supports reflecting
into specific elements and modifying without knowing the whole file structure,
then we could do away with the Java binary for maven preprocessing.

## CMake

Coming soon™.  https://github.com/google/kythe/issues/2861

## Bazel

Actually we have no custom work here.  We extract compilation records from Bazel
using the extra action mechanism.  The extractrepo tool therefore doesn't handle
Bazel directly, but repositories using Bazel for languages we already support
should work without extra effort.

