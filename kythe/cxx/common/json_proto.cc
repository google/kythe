/*
 * Copyright 2015 Google Inc. All rights reserved.
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

#include "json_proto.h"

#include "absl/strings/escaping.h"
#include "glog/logging.h"
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/io/zero_copy_stream_impl.h"
#include "google/protobuf/message.h"
#include "google/protobuf/util/json_util.h"
#include "google/protobuf/util/type_resolver.h"
#include "google/protobuf/util/type_resolver_util.h"
#include "rapidjson/document.h"
#include "rapidjson/filewritestream.h"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/writer.h"

namespace kythe {
namespace {
using ::google::protobuf::DescriptorPool;
using ::google::protobuf::util::TypeResolver;

class PermissiveTypeResolver : public TypeResolver {
 public:
  explicit PermissiveTypeResolver(const DescriptorPool *pool)
      : impl_(google::protobuf::util::NewTypeResolverForDescriptorPool("",
                                                                       pool)) {}

  google::protobuf::util::Status ResolveMessageType(
      const std::string &type_url,
      google::protobuf::Type *message_type) override {
    absl::string_view adjusted = type_url;
    adjusted.remove_prefix(type_url.rfind('/') + 1);
    return impl_->ResolveMessageType(absl::StrCat("/", adjusted), message_type);
  }

  google::protobuf::util::Status ResolveEnumType(
      const std::string &type_url, google::protobuf::Enum *enum_type) override {
    absl::string_view adjusted = type_url;
    adjusted.remove_prefix(type_url.rfind('/') + 1);
    return impl_->ResolveEnumType(absl::StrCat("/", adjusted), enum_type);
  }

 private:
  std::unique_ptr<TypeResolver> impl_;
};

TypeResolver *GetGeneratedTypeResolver() {
  static TypeResolver *generated_resolver =
      new PermissiveTypeResolver(DescriptorPool::generated_pool());
  return generated_resolver;
}

struct MaybeDeleteResolver {
  void operator()(TypeResolver *resolver) const {
    if (resolver != GetGeneratedTypeResolver()) {
      delete resolver;
    }
  }
};

std::unique_ptr<TypeResolver, MaybeDeleteResolver> MakeTypeResolverForPool(
    const DescriptorPool *pool) {
  if (pool == DescriptorPool::generated_pool()) {
    return std::unique_ptr<TypeResolver, MaybeDeleteResolver>(
        GetGeneratedTypeResolver());
  }
  return std::unique_ptr<TypeResolver, MaybeDeleteResolver>(
      new PermissiveTypeResolver(pool));
}

}  // namespace

bool WriteMessageAsJsonToString(const google::protobuf::Message &message,
                                std::string *out) {
  auto resolver =
      MakeTypeResolverForPool(message.GetDescriptor()->file()->pool());

  google::protobuf::util::JsonPrintOptions options;
  options.preserve_proto_field_names = true;

  auto status = google::protobuf::util::BinaryToJsonString(
      resolver.get(), message.GetDescriptor()->full_name(),
      message.SerializeAsString(), out, options);

  if (!status.ok()) {
    LOG(ERROR) << status.ToString();
  }
  return status.ok();
}

bool WriteMessageAsJsonToString(const google::protobuf::Message &message,
                                const std::string &format_key,
                                std::string *out) {
  rapidjson::StringBuffer buffer;
  rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
  writer.StartObject();
  writer.Key("format");
  writer.String(format_key.c_str());
  writer.Key("content");
  {
    std::string content;
    if (!WriteMessageAsJsonToString(message, &content)) {
      return false;
    }
    writer.RawValue(content.data(), content.size(), rapidjson::kObjectType);
  }
  writer.EndObject();
  *out = buffer.GetString();
  return true;
}

bool MergeJsonWithMessage(const std::string &in, std::string *format_key,
                          google::protobuf::Message *message) {
  rapidjson::Document document;
  document.Parse(in.c_str());
  if (document.HasParseError()) {
    return false;
  }
  if (!document.IsObject() || !document.HasMember("format") ||
      !document.HasMember("content") || !document["format"].IsString() ||
      !document["content"].IsObject()) {
    return false;
  }
  std::string in_format = document["format"].GetString();
  if (format_key) {
    *format_key = in_format;
  }
  if (in_format == "kythe") {
    std::string content = [&] {
      rapidjson::StringBuffer buffer;
      rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
      CHECK(document["content"].Accept(writer));
      return std::string(buffer.GetString());
    }();
    LOG(ERROR) << content;
    std::string binary;

    auto resolver =
        MakeTypeResolverForPool(message->GetDescriptor()->file()->pool());

    auto status = google::protobuf::util::JsonToBinaryString(
        resolver.get(), message->GetDescriptor()->full_name(), content, &binary,
        {});

    if (!status.ok()) {
      LOG(ERROR) << status.ToString() << ": " << content;
      return false;
    }
    return message->ParseFromString(binary);
  }
  return false;
}

void PackAny(const google::protobuf::Message &message, const char *type_uri,
             google::protobuf::Any *out) {
  out->set_type_url(type_uri);
  google::protobuf::io::StringOutputStream stream(out->mutable_value());
  google::protobuf::io::CodedOutputStream coded_output_stream(&stream);
  message.SerializeToCodedStream(&coded_output_stream);
}

bool UnpackAny(const google::protobuf::Any &any,
               google::protobuf::Message *result) {
  google::protobuf::io::ArrayInputStream stream(any.value().data(),
                                                any.value().size());
  google::protobuf::io::CodedInputStream coded_input_stream(&stream);
  return result->ParseFromCodedStream(&coded_input_stream);
}
}  // namespace kythe
