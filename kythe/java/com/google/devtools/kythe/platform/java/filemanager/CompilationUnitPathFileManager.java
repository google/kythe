/*
 * Copyright 2019 The Kythe Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.google.devtools.kythe.platform.java.filemanager;

import com.google.common.collect.ImmutableList;
import com.google.common.collect.ImmutableMap;
import com.google.common.collect.ImmutableSet;
import com.google.common.collect.Iterables;
import com.google.common.collect.Iterators;
import com.google.common.collect.Maps;
import com.google.common.flogger.FluentLogger;
import com.google.devtools.kythe.extractors.java.JavaCompilationUnitExtractor;
import com.google.devtools.kythe.platform.shared.FileDataProvider;
import com.google.devtools.kythe.platform.shared.filesystem.CompilationUnitFileSystem;
import com.google.devtools.kythe.proto.Analysis.CompilationUnit;
import com.google.devtools.kythe.proto.Java.JavaDetails;
import com.google.protobuf.Any;
import com.google.protobuf.InvalidProtocolBufferException;
import com.sun.tools.javac.main.Option;
import com.sun.tools.javac.main.OptionHelper;
import java.io.IOException;
import java.nio.file.FileVisitResult;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.Collection;
import java.util.Iterator;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Stream;
import javax.tools.StandardJavaFileManager;
import javax.tools.StandardLocation;
import org.checkerframework.checker.nullness.qual.Nullable;

/**
 * StandardJavaFileManager which uses a CompilationUnitFileSystem for managing paths based on on the
 * paths provided in the CompilationUnit.
 */
@com.sun.tools.javac.api.ClientCodeWrapper.Trusted
public final class CompilationUnitPathFileManager extends ForwardingStandardJavaFileManager {
  private static final FluentLogger logger = FluentLogger.forEnclosingClass();

  private final CompilationUnitFileSystem fileSystem;
  private final ImmutableSet<String> defaultPlatformClassPath;
  // The path given to us that we are allowed to write in. This will be stored as an absolute
  // path.
  private final @Nullable Path temporaryDirectoryPrefix;
  // A temporary directory inside of temporaryDirectoryPrefix that we will use and delete when the
  // close method is called. This should be will as an absolute path.
  private @Nullable Path temporaryDirectory;

  public CompilationUnitPathFileManager(
      CompilationUnit compilationUnit,
      FileDataProvider fileDataProvider,
      StandardJavaFileManager fileManager,
      @Nullable Path temporaryDirectory) {
    super(fileManager);
    // Store the absolute path so we can safely do startsWith checks later.
    this.temporaryDirectoryPrefix =
        temporaryDirectory == null ? null : temporaryDirectory.toAbsolutePath();
    defaultPlatformClassPath =
        ImmutableSet.copyOf(
            Iterables.transform(
                // TODO(shahms): use getLocationAsPaths on Java 9
                fileManager.getLocation(StandardLocation.PLATFORM_CLASS_PATH),
                f -> f.toPath().normalize().toString()));

    fileSystem = CompilationUnitFileSystem.create(compilationUnit, fileDataProvider);
    // When compiled for Java 9+ this is ambiguous, so disambiguate to the compatibility shim.
    setPathFactory((ForwardingStandardJavaFileManager.PathFactory) this::getPath);
    setLocations(
        findJavaDetails(compilationUnit)
            .map(details -> toLocationMap(details))
            .orElseGet(() -> logEmptyLocationMap()));
  }

  @Override
  public boolean handleOption(String current, Iterator<String> remaining) {
    if (Option.SYSTEM.matches(current)) {
      try {
        Option.SYSTEM.handleOption(
            new OptionHelper.GrumpyHelper(null) {
              @Override
              public void put(String name, String value) {}

              @Override
              public boolean handleFileManagerOption(Option unused, String value) {
                try {
                  setSystemOption(value);
                } catch (IOException exc) {
                  throw new IllegalArgumentException(exc);
                }
                return true;
              }
            },
            current,
            remaining);
      } catch (Option.InvalidValueException exc) {
        throw new IllegalArgumentException(exc);
      }
      return true;
    }
    return super.handleOption(current, remaining);
  }

  @Override
  public void close() throws IOException {
    super.close();
    if (temporaryDirectory != null) {
      Files.walkFileTree(
          temporaryDirectory,
          new SimpleFileVisitor<Path>() {
            @Override
            public FileVisitResult visitFile(Path file, BasicFileAttributes attrs)
                throws IOException {
              Files.delete(file);
              return FileVisitResult.CONTINUE;
            }

            @Override
            public FileVisitResult postVisitDirectory(Path dir, IOException exc)
                throws IOException {
              Files.delete(dir);
              return FileVisitResult.CONTINUE;
            }
          });
    }
  }

  /** Extracts the embedded JavaDetails message, if any, from the CompilationUnit. */
  private static Optional<JavaDetails> findJavaDetails(CompilationUnit unit) {
    for (Any details : unit.getDetailsList()) {
      try {
        if (details.getTypeUrl().equals(JavaCompilationUnitExtractor.JAVA_DETAILS_URL)) {
          return Optional.of(JavaDetails.parseFrom(details.getValue()));
        } else if (details.is(JavaDetails.class)) {
          return Optional.of(details.unpack(JavaDetails.class));
        }
      } catch (InvalidProtocolBufferException ipbe) {
        logger.atWarning().withCause(ipbe).log("error unpacking JavaDetails");
      }
    }
    return Optional.empty();
  }

  /** Logs that path handling will fall back to Javac's option parsing. */
  private static Map<Location, Collection<Path>> logEmptyLocationMap() {
    // It's expected that extractors which use JavaDetails will remove the corresponding
    // arguments from the command line.  Those extractors which don't use JavaDetails
    // (or options not present in the details), will remain on the command line and be
    // parsed as normal, relying on getPath() to map into the compilation unit.
    logger.atInfo().log("Compilation missing JavaDetails; falling back to flag parsing");
    return ImmutableMap.of();
  }

  /** Translates the JavaDetails locations into {@code Map<Location, Collection<Path>>} */
  private Map<Location, Collection<Path>> toLocationMap(JavaDetails details) {
    return Maps.filterValues(
        new ImmutableMap.Builder<Location, Collection<Path>>()
            .put(StandardLocation.CLASS_PATH, toPaths(details.getClasspathList()))
            .put(StandardLocation.locationFor("MODULE_PATH"), toPaths(details.getClasspathList()))
            .put(StandardLocation.SOURCE_PATH, toPaths(details.getSourcepathList()))
            .put(
                StandardLocation.PLATFORM_CLASS_PATH,
                // bootclasspath should fall back to the local filesystem,
                // while the others should only come from the compilation unit.
                toPaths(details.getBootclasspathList(), this::getPath))
            .build(),
        v -> !v.isEmpty());
  }

  private Collection<Path> toPaths(Collection<String> paths) {
    return toPaths(paths, fileSystem::getPath);
  }

  private Collection<Path> toPaths(Collection<String> paths, PathFactory factory) {
    return paths.stream().map(factory::getPath).collect(ImmutableList.toImmutableList());
  }

  private Path getPath(String path, String... rest) {
    // Get the absolute path so we can safely do the startsWith check.
    Path local = Paths.get(path, rest).toAbsolutePath();
    // If this is a path underneath the temporary directory, use it. This is required for --system
    // flags to work correctly.
    if (temporaryDirectory != null) {
      if (local.startsWith(temporaryDirectory)) {
        logger.atInfo().log("Using the filesystem for temporary path %s", local);
        return local;
      }
    }
    // In order to support paths passed via command line options rather than
    // JavaDetails, prevent loading source files from outside the
    // CompilationUnit (#818), and allow JDK classes to be provided by the
    // platform we always form CompilationUnit paths unless that path is not
    // present in the CompilationUnit and is part of the default boot class path.
    Path result = fileSystem.getPath(path, rest);
    if (Files.exists(result) || !defaultPlatformClassPath.contains(result.normalize().toString())) {
      return result;
    }
    logger.atFine().log("Falling back to filesystem for %s", result);
    return Paths.get(path, rest);
  }

  /** For each entry in the provided map, sets the corresponding location in fileManager */
  private void setLocations(Map<Location, Collection<Path>> locations) {
    for (Map.Entry<Location, Collection<Path>> entry : locations.entrySet()) {
      try {
        setLocationFromPaths(entry.getKey(), entry.getValue());
      } catch (IOException ex) {
        logger.atWarning().withCause(ex).log("error setting location %s", entry);
      }
    }
  }

  private void setSystemOption(String value) throws IOException {
    // There are two kinds of --system flags we need to support:
    //   1) Bundled system images, with a lib/jrt-fs.jar and lib/modules image.
    //   2) Exploded system images, where the modules live under a modules subdirectory.
    // The former must reside in the filesystem as there are multifarious asserts and checks that
    // this is so.
    Path sys = fileSystem.getPath(value).normalize();
    if (Files.exists(sys.resolve("lib").resolve("jrt-fs.jar"))) {
      if (temporaryDirectoryPrefix == null) {
        logger.atSevere().log(
            "Can't create temporary directory to store system module because no temporary"
                + " directory was provided");
        throw new IllegalArgumentException("temporary directory needed but not provided");
      }
      if (temporaryDirectory != null) {
        throw new IllegalStateException("Temporary directory set twice");
      }
      temporaryDirectory =
          Files.createTempDirectory(temporaryDirectoryPrefix, "kythe_java_indexer")
              .toAbsolutePath();
      Path systemRoot =
          Files.createDirectory(temporaryDirectory.resolve("system")).toAbsolutePath();
      try (Stream<Path> stream = Files.walk(sys.resolve("lib"))) {
        for (Path path : (Iterable<Path>) stream::iterator) {
          Path p =
              Files.copy(
                  path,
                  systemRoot.resolve(
                      path.subpath(sys.getNameCount(), path.getNameCount()).toString()));
          logger.atInfo().log("Copied file to %s", p.toAbsolutePath());
        }
      }
      logger.atInfo().log("Setting system path to %s", systemRoot);
      super.handleOption("--system", Iterators.singletonIterator(systemRoot.toString()));
    } else if (Files.isDirectory(sys.resolve("modules"))) {
      // TODO(salguarnieri) Due to a bug in the javac argument validation, bypass it and set the
      // location directly.
      setLocationFromPaths(
          StandardLocation.valueOf("SYSTEM_MODULES"), ImmutableList.of(sys.resolve("modules")));
    } else {
      throw new IllegalArgumentException(value);
    }
  }
}
