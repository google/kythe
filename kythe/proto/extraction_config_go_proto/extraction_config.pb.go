// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kythe/proto/extraction_config.proto

package extraction_config_go_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ExtractionConfiguration struct {
	RequiredImage        []*ExtractionConfiguration_Image      `protobuf:"bytes,1,rep,name=required_image,json=requiredImage,proto3" json:"required_image,omitempty"`
	RunCommand           []*ExtractionConfiguration_RunCommand `protobuf:"bytes,2,rep,name=run_command,json=runCommand,proto3" json:"run_command,omitempty"`
	EntryPoint           []string                              `protobuf:"bytes,3,rep,name=entry_point,json=entryPoint,proto3" json:"entry_point,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *ExtractionConfiguration) Reset()         { *m = ExtractionConfiguration{} }
func (m *ExtractionConfiguration) String() string { return proto.CompactTextString(m) }
func (*ExtractionConfiguration) ProtoMessage()    {}
func (*ExtractionConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_extraction_config_2b39cab2ce2c56d9, []int{0}
}
func (m *ExtractionConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtractionConfiguration.Unmarshal(m, b)
}
func (m *ExtractionConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtractionConfiguration.Marshal(b, m, deterministic)
}
func (dst *ExtractionConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtractionConfiguration.Merge(dst, src)
}
func (m *ExtractionConfiguration) XXX_Size() int {
	return xxx_messageInfo_ExtractionConfiguration.Size(m)
}
func (m *ExtractionConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtractionConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_ExtractionConfiguration proto.InternalMessageInfo

func (m *ExtractionConfiguration) GetRequiredImage() []*ExtractionConfiguration_Image {
	if m != nil {
		return m.RequiredImage
	}
	return nil
}

func (m *ExtractionConfiguration) GetRunCommand() []*ExtractionConfiguration_RunCommand {
	if m != nil {
		return m.RunCommand
	}
	return nil
}

func (m *ExtractionConfiguration) GetEntryPoint() []string {
	if m != nil {
		return m.EntryPoint
	}
	return nil
}

type ExtractionConfiguration_Image struct {
	Uri                  string                              `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	Name                 string                              `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	CopySpec             []*ExtractionConfiguration_CopySpec `protobuf:"bytes,3,rep,name=copy_spec,json=copySpec,proto3" json:"copy_spec,omitempty"`
	EnvVar               []*ExtractionConfiguration_EnvVar   `protobuf:"bytes,4,rep,name=env_var,json=envVar,proto3" json:"env_var,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *ExtractionConfiguration_Image) Reset()         { *m = ExtractionConfiguration_Image{} }
func (m *ExtractionConfiguration_Image) String() string { return proto.CompactTextString(m) }
func (*ExtractionConfiguration_Image) ProtoMessage()    {}
func (*ExtractionConfiguration_Image) Descriptor() ([]byte, []int) {
	return fileDescriptor_extraction_config_2b39cab2ce2c56d9, []int{0, 0}
}
func (m *ExtractionConfiguration_Image) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtractionConfiguration_Image.Unmarshal(m, b)
}
func (m *ExtractionConfiguration_Image) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtractionConfiguration_Image.Marshal(b, m, deterministic)
}
func (dst *ExtractionConfiguration_Image) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtractionConfiguration_Image.Merge(dst, src)
}
func (m *ExtractionConfiguration_Image) XXX_Size() int {
	return xxx_messageInfo_ExtractionConfiguration_Image.Size(m)
}
func (m *ExtractionConfiguration_Image) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtractionConfiguration_Image.DiscardUnknown(m)
}

var xxx_messageInfo_ExtractionConfiguration_Image proto.InternalMessageInfo

func (m *ExtractionConfiguration_Image) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

func (m *ExtractionConfiguration_Image) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ExtractionConfiguration_Image) GetCopySpec() []*ExtractionConfiguration_CopySpec {
	if m != nil {
		return m.CopySpec
	}
	return nil
}

func (m *ExtractionConfiguration_Image) GetEnvVar() []*ExtractionConfiguration_EnvVar {
	if m != nil {
		return m.EnvVar
	}
	return nil
}

type ExtractionConfiguration_CopySpec struct {
	Source               string   `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Destination          string   `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExtractionConfiguration_CopySpec) Reset()         { *m = ExtractionConfiguration_CopySpec{} }
func (m *ExtractionConfiguration_CopySpec) String() string { return proto.CompactTextString(m) }
func (*ExtractionConfiguration_CopySpec) ProtoMessage()    {}
func (*ExtractionConfiguration_CopySpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_extraction_config_2b39cab2ce2c56d9, []int{0, 1}
}
func (m *ExtractionConfiguration_CopySpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtractionConfiguration_CopySpec.Unmarshal(m, b)
}
func (m *ExtractionConfiguration_CopySpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtractionConfiguration_CopySpec.Marshal(b, m, deterministic)
}
func (dst *ExtractionConfiguration_CopySpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtractionConfiguration_CopySpec.Merge(dst, src)
}
func (m *ExtractionConfiguration_CopySpec) XXX_Size() int {
	return xxx_messageInfo_ExtractionConfiguration_CopySpec.Size(m)
}
func (m *ExtractionConfiguration_CopySpec) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtractionConfiguration_CopySpec.DiscardUnknown(m)
}

var xxx_messageInfo_ExtractionConfiguration_CopySpec proto.InternalMessageInfo

func (m *ExtractionConfiguration_CopySpec) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *ExtractionConfiguration_CopySpec) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

type ExtractionConfiguration_EnvVar struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExtractionConfiguration_EnvVar) Reset()         { *m = ExtractionConfiguration_EnvVar{} }
func (m *ExtractionConfiguration_EnvVar) String() string { return proto.CompactTextString(m) }
func (*ExtractionConfiguration_EnvVar) ProtoMessage()    {}
func (*ExtractionConfiguration_EnvVar) Descriptor() ([]byte, []int) {
	return fileDescriptor_extraction_config_2b39cab2ce2c56d9, []int{0, 2}
}
func (m *ExtractionConfiguration_EnvVar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtractionConfiguration_EnvVar.Unmarshal(m, b)
}
func (m *ExtractionConfiguration_EnvVar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtractionConfiguration_EnvVar.Marshal(b, m, deterministic)
}
func (dst *ExtractionConfiguration_EnvVar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtractionConfiguration_EnvVar.Merge(dst, src)
}
func (m *ExtractionConfiguration_EnvVar) XXX_Size() int {
	return xxx_messageInfo_ExtractionConfiguration_EnvVar.Size(m)
}
func (m *ExtractionConfiguration_EnvVar) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtractionConfiguration_EnvVar.DiscardUnknown(m)
}

var xxx_messageInfo_ExtractionConfiguration_EnvVar proto.InternalMessageInfo

func (m *ExtractionConfiguration_EnvVar) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ExtractionConfiguration_EnvVar) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type ExtractionConfiguration_RunCommand struct {
	Command              string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Arg                  []string `protobuf:"bytes,2,rep,name=arg,proto3" json:"arg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExtractionConfiguration_RunCommand) Reset()         { *m = ExtractionConfiguration_RunCommand{} }
func (m *ExtractionConfiguration_RunCommand) String() string { return proto.CompactTextString(m) }
func (*ExtractionConfiguration_RunCommand) ProtoMessage()    {}
func (*ExtractionConfiguration_RunCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_extraction_config_2b39cab2ce2c56d9, []int{0, 3}
}
func (m *ExtractionConfiguration_RunCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtractionConfiguration_RunCommand.Unmarshal(m, b)
}
func (m *ExtractionConfiguration_RunCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtractionConfiguration_RunCommand.Marshal(b, m, deterministic)
}
func (dst *ExtractionConfiguration_RunCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtractionConfiguration_RunCommand.Merge(dst, src)
}
func (m *ExtractionConfiguration_RunCommand) XXX_Size() int {
	return xxx_messageInfo_ExtractionConfiguration_RunCommand.Size(m)
}
func (m *ExtractionConfiguration_RunCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtractionConfiguration_RunCommand.DiscardUnknown(m)
}

var xxx_messageInfo_ExtractionConfiguration_RunCommand proto.InternalMessageInfo

func (m *ExtractionConfiguration_RunCommand) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *ExtractionConfiguration_RunCommand) GetArg() []string {
	if m != nil {
		return m.Arg
	}
	return nil
}

func init() {
	proto.RegisterType((*ExtractionConfiguration)(nil), "kythe.proto.ExtractionConfiguration")
	proto.RegisterType((*ExtractionConfiguration_Image)(nil), "kythe.proto.ExtractionConfiguration.Image")
	proto.RegisterType((*ExtractionConfiguration_CopySpec)(nil), "kythe.proto.ExtractionConfiguration.CopySpec")
	proto.RegisterType((*ExtractionConfiguration_EnvVar)(nil), "kythe.proto.ExtractionConfiguration.EnvVar")
	proto.RegisterType((*ExtractionConfiguration_RunCommand)(nil), "kythe.proto.ExtractionConfiguration.RunCommand")
}

func init() {
	proto.RegisterFile("kythe/proto/extraction_config.proto", fileDescriptor_extraction_config_2b39cab2ce2c56d9)
}

var fileDescriptor_extraction_config_2b39cab2ce2c56d9 = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4a, 0xfb, 0x40,
	0x10, 0xc6, 0x49, 0xd3, 0xa6, 0xcd, 0x84, 0xff, 0x1f, 0x59, 0x44, 0x43, 0x10, 0x0c, 0x7a, 0x29,
	0x8a, 0x29, 0xd4, 0x8b, 0x67, 0xdb, 0x1e, 0xf4, 0x54, 0x57, 0xf0, 0xe0, 0x25, 0xac, 0xe9, 0x1a,
	0x83, 0x66, 0x37, 0x6e, 0xb3, 0xc1, 0x3c, 0x9d, 0x4f, 0xe1, 0xfb, 0xc8, 0x4e, 0x92, 0x5a, 0x10,
	0xa1, 0xa7, 0x9d, 0xef, 0x63, 0xe6, 0x97, 0xfd, 0x66, 0x03, 0xa7, 0xaf, 0x75, 0xf9, 0xc2, 0x27,
	0x85, 0x92, 0xa5, 0x9c, 0xf0, 0x8f, 0x52, 0xb1, 0xa4, 0xcc, 0xa4, 0x88, 0x13, 0x29, 0x9e, 0xb3,
	0x34, 0x42, 0x9f, 0x78, 0xd8, 0xd4, 0x88, 0x93, 0xaf, 0x3e, 0x1c, 0x2e, 0x36, 0x8d, 0x33, 0xec,
	0xd3, 0x8a, 0x19, 0x41, 0xee, 0xe0, 0xbf, 0xe2, 0xef, 0x3a, 0x53, 0x7c, 0x15, 0x67, 0x39, 0x4b,
	0xb9, 0x6f, 0x85, 0xf6, 0xd8, 0x9b, 0x9e, 0x45, 0x5b, 0x84, 0xe8, 0x8f, 0xe9, 0xe8, 0xc6, 0x4c,
	0xd0, 0x7f, 0x1d, 0x01, 0x25, 0x59, 0x82, 0xa7, 0xb4, 0xb9, 0x4f, 0x9e, 0x33, 0xb1, 0xf2, 0x7b,
	0xc8, 0x9b, 0xec, 0xc4, 0xa3, 0x5a, 0xcc, 0x9a, 0x31, 0x0a, 0x6a, 0x53, 0x93, 0x63, 0xf0, 0xb8,
	0x28, 0x55, 0x1d, 0x17, 0x32, 0x13, 0xa5, 0x6f, 0x87, 0xf6, 0xd8, 0xa5, 0x80, 0xd6, 0xd2, 0x38,
	0xc1, 0xa7, 0x05, 0x83, 0xe6, 0xe3, 0x7b, 0x60, 0x6b, 0x95, 0xf9, 0x56, 0x68, 0x8d, 0x5d, 0x6a,
	0x4a, 0x42, 0xa0, 0x2f, 0x58, 0xce, 0xfd, 0x1e, 0x5a, 0x58, 0x93, 0x5b, 0x70, 0x13, 0x59, 0xd4,
	0xf1, 0xba, 0xe0, 0x09, 0xe2, 0xbc, 0xe9, 0xc5, 0x4e, 0x17, 0x9c, 0xc9, 0xa2, 0xbe, 0x2f, 0x78,
	0x42, 0x47, 0x49, 0x5b, 0x91, 0x39, 0x0c, 0xb9, 0xa8, 0xe2, 0x8a, 0x29, 0xbf, 0x8f, 0xa4, 0xf3,
	0x9d, 0x48, 0x0b, 0x51, 0x3d, 0x30, 0x45, 0x1d, 0x8e, 0x67, 0x30, 0x87, 0x51, 0xc7, 0x26, 0x07,
	0xe0, 0xac, 0xa5, 0x56, 0x09, 0x6f, 0x63, 0xb4, 0x8a, 0x84, 0xe0, 0xad, 0xf8, 0xba, 0xcc, 0x04,
	0x12, 0xda, 0x40, 0xdb, 0x56, 0x30, 0x05, 0xa7, 0xe1, 0x6e, 0x52, 0x5b, 0x5b, 0xa9, 0xf7, 0x61,
	0x50, 0xb1, 0x37, 0xdd, 0xad, 0xa2, 0x11, 0xc1, 0x15, 0xc0, 0xcf, 0xda, 0x89, 0x0f, 0xc3, 0xee,
	0xe1, 0x9a, 0xd1, 0x4e, 0x9a, 0xcd, 0x32, 0x95, 0xe2, 0x73, 0xba, 0xd4, 0x94, 0xd7, 0x47, 0x8f,
	0xc1, 0xaf, 0xff, 0x2f, 0x4e, 0x65, 0x8c, 0xc1, 0x9f, 0x1c, 0x3c, 0x2e, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x54, 0xeb, 0x38, 0x63, 0xb0, 0x02, 0x00, 0x00,
}
