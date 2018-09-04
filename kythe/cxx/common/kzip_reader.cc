/*
 * Copyright 2018 The Kythe Authors. All rights reserved.
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

#include "kythe/cxx/common/kzip_reader.h"

#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/strip.h"
#include "absl/types/optional.h"
#include "glog/logging.h"
#include "google/protobuf/io/zero_copy_stream.h"
#include "google/protobuf/io/zero_copy_stream_impl.h"
#include "kythe/cxx/common/json_proto.h"
#include "kythe/cxx/common/libzip/error.h"
#include "kythe/proto/analysis.pb.h"

#include <openssl/sha.h>
#include "absl/strings/escaping.h"

namespace kythe {
namespace {

struct ZipFileClose {
  void operator()(zip_file_t* file) {
    if (file != nullptr) {
      CHECK(zip_fclose(file) == 0);
    }
  }
};
using ZipFile = std::unique_ptr<zip_file_t, ZipFileClose>;

class ZipFileInputStream : public google::protobuf::io::ZeroCopyInputStream {
 public:
  explicit ZipFileInputStream(zip_file_t* file) : input_(file) {}

  bool Next(const void** data, int* size) override {
    return impl_.Next(data, size);
  }

  void BackUp(int count) override { impl_.BackUp(count); }
  bool Skip(int count) override { return impl_.Skip(count); }
  google::protobuf::int64 ByteCount() const override {
    return impl_.ByteCount();
  }

 private:
  class CopyingZipInputStream
      : public google::protobuf::io::CopyingInputStream {
   public:
    explicit CopyingZipInputStream(zip_file_t* file) : file_(file) {}

    int Read(void* buffer, int size) override {
      return zip_fread(file_, buffer, size);
    }

    int Skip(int count) override {
      zip_int64_t start = zip_ftell(file_);
      if (start < 0) {
        return 0;
      }
      if (zip_fseek(file_, count, SEEK_CUR) < 0) {
        return 0;
      }
      zip_int64_t end = zip_ftell(file_);
      if (end < 0) {
        return 0;
      }
      return end - start;
    }

   private:
    zip_file_t* file_;
  };

  CopyingZipInputStream input_;
  google::protobuf::io::CopyingInputStreamAdaptor impl_{&input_};
};

absl::optional<absl::string_view> UnitDigest(absl::string_view path) {
  path.remove_prefix(std::min(path.find('/'), path.size()));
  if (!absl::ConsumePrefix(&path, "/units/")) {
    return absl::nullopt;
  }
  return path;
}

StatusOr<absl::string_view> Validate(zip_t* archive) {
  if (!zip_get_num_entries(archive, 0)) {
    return InvalidArgumentError("Empty kzip archive");
  }

  // Pull the root directory from an arbitrary entry.
  absl::string_view root = zip_get_name(archive, 0, 0);
  auto slashpos = root.find('/');
  if (slashpos == 0 || slashpos == absl::string_view::npos) {
    return InvalidArgumentError(
        absl::StrCat("Malformed kzip: invalid root: ", root));
  }
  root.remove_suffix(root.size() - slashpos);
  VLOG(1) << "Using archive root: " << root;
  for (int i = 0; i < zip_get_num_entries(archive, 0); ++i) {
    absl::string_view name = zip_get_name(archive, i, 0);
    if (!absl::StartsWith(name, root)) {
      return InvalidArgumentError(
          absl::StrCat("Malformed kzip: invalid entry: ", name));
    }
  }
  return root;
}

absl::optional<zip_uint64_t> FileSize(zip_t* archive, zip_uint64_t index) {
  zip_stat_t sb;
  zip_stat_init(&sb);

  if (zip_stat_index(archive, index, ZIP_STAT_SIZE, &sb) < 0) {
    return absl::nullopt;
  }
  return sb.size;
}

StatusOr<std::string> ReadTextFile(zip_t* archive, const std::string& path) {
  zip_int64_t index = zip_name_locate(archive, path.c_str(), 0);
  if (index >= 0) {
    if (auto file = ZipFile(zip_fopen_index(archive, index, 0))) {
      if (auto size = FileSize(archive, index)) {
        std::string result(*size, '\0');
        if (zip_fread(file.get(), &result.front(), *size) == *size) {
          return result;
        } else {
          return libzip::ToStatus(zip_file_get_error(file.get()));
        }
      }
    }
  }
  Status status = libzip::ToStatus(zip_get_error(archive));
  if (!status.ok()) {
    return status;
  }
  return UnknownError(absl::StrCat("Unable to read: ", path));
}

}  // namespace

/* static */
StatusOr<IndexReader> KzipReader::Open(absl::string_view path) {
  int error;
  if (auto archive =
          ZipHandle(zip_open(std::string(path).c_str(), ZIP_RDONLY, &error))) {
    if (auto root = Validate(archive.get())) {
      return IndexReader(
          absl::WrapUnique(new KzipReader(std::move(archive), *root)));
    } else {
      return root.status();
    }
  }
  return libzip::Error(error).ToStatus();
}

/* static */
StatusOr<IndexReader> KzipReader::FromSource(zip_source_t* source) {
  libzip::Error error;
  if (auto archive =
          ZipHandle(zip_open_from_source(source, ZIP_RDONLY, error.get()))) {
    if (auto root = Validate(archive.get())) {
      return IndexReader(
          absl::WrapUnique(new KzipReader(std::move(archive), *root)));
    } else {
      return root.status();
    }
  }
  return error.ToStatus();
}

KzipReader::KzipReader(ZipHandle archive, absl::string_view root)
    : archive_(std::move(archive)), root_(root) {}

StatusOr<proto::IndexedCompilation> KzipReader::ReadUnit(
    absl::string_view digest) {
  std::string path = absl::StrCat(root_, "/units/", digest);
  if (auto file = ZipFile(zip_fopen(archive(), path.c_str(), 0))) {
    proto::IndexedCompilation unit;
    ZipFileInputStream input(file.get());
    Status status = ParseFromJsonStream(&input, &unit);
    if (!status.ok()) {
      Status zip_status = libzip::ToStatus(zip_file_get_error(file.get()));
      if (!zip_status.ok()) {
        // Prefer the underlying zip error, if present.
        return zip_status;
      }
      return status;
    }
    return unit;
  }
  Status status = libzip::ToStatus(zip_get_error(archive()));
  if (!status.ok()) {
    return status;
  }
  return UnknownError(absl::StrCat("Unable to open unit ", digest));
}

StatusOr<std::string> KzipReader::ReadFile(absl::string_view digest) {
  return ReadTextFile(archive(), absl::StrCat(root_, "/files/", digest));
}

Status KzipReader::Scan(const ScanCallback& callback) {
  for (int i = 0; i < zip_get_num_entries(archive(), 0); ++i) {
    if (auto digest = UnitDigest(zip_get_name(archive(), i, 0))) {
      if (!callback(*digest)) {
        break;
      }
    }
  }
  return OkStatus();
}

}  // namespace kythe
