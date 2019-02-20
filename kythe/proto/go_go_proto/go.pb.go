// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kythe/proto/go.proto

package go_go_proto

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

type GoDetails struct {
	Goos                 string   `protobuf:"bytes,1,opt,name=goos,proto3" json:"goos,omitempty"`
	Goarch               string   `protobuf:"bytes,2,opt,name=goarch,proto3" json:"goarch,omitempty"`
	Goroot               string   `protobuf:"bytes,3,opt,name=goroot,proto3" json:"goroot,omitempty"`
	Gopath               string   `protobuf:"bytes,4,opt,name=gopath,proto3" json:"gopath,omitempty"`
	Compiler             string   `protobuf:"bytes,5,opt,name=compiler,proto3" json:"compiler,omitempty"`
	BuildTags            []string `protobuf:"bytes,6,rep,name=build_tags,json=buildTags,proto3" json:"build_tags,omitempty"`
	CgoEnabled           bool     `protobuf:"varint,7,opt,name=cgo_enabled,json=cgoEnabled,proto3" json:"cgo_enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GoDetails) Reset()         { *m = GoDetails{} }
func (m *GoDetails) String() string { return proto.CompactTextString(m) }
func (*GoDetails) ProtoMessage()    {}
func (*GoDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_go_0251d8e3856873be, []int{0}
}
func (m *GoDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GoDetails.Unmarshal(m, b)
}
func (m *GoDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GoDetails.Marshal(b, m, deterministic)
}
func (dst *GoDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GoDetails.Merge(dst, src)
}
func (m *GoDetails) XXX_Size() int {
	return xxx_messageInfo_GoDetails.Size(m)
}
func (m *GoDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_GoDetails.DiscardUnknown(m)
}

var xxx_messageInfo_GoDetails proto.InternalMessageInfo

func (m *GoDetails) GetGoos() string {
	if m != nil {
		return m.Goos
	}
	return ""
}

func (m *GoDetails) GetGoarch() string {
	if m != nil {
		return m.Goarch
	}
	return ""
}

func (m *GoDetails) GetGoroot() string {
	if m != nil {
		return m.Goroot
	}
	return ""
}

func (m *GoDetails) GetGopath() string {
	if m != nil {
		return m.Gopath
	}
	return ""
}

func (m *GoDetails) GetCompiler() string {
	if m != nil {
		return m.Compiler
	}
	return ""
}

func (m *GoDetails) GetBuildTags() []string {
	if m != nil {
		return m.BuildTags
	}
	return nil
}

func (m *GoDetails) GetCgoEnabled() bool {
	if m != nil {
		return m.CgoEnabled
	}
	return false
}

type GoPackageInfo struct {
	ImportPath           string   `protobuf:"bytes,1,opt,name=import_path,json=importPath,proto3" json:"import_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GoPackageInfo) Reset()         { *m = GoPackageInfo{} }
func (m *GoPackageInfo) String() string { return proto.CompactTextString(m) }
func (*GoPackageInfo) ProtoMessage()    {}
func (*GoPackageInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_go_0251d8e3856873be, []int{1}
}
func (m *GoPackageInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GoPackageInfo.Unmarshal(m, b)
}
func (m *GoPackageInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GoPackageInfo.Marshal(b, m, deterministic)
}
func (dst *GoPackageInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GoPackageInfo.Merge(dst, src)
}
func (m *GoPackageInfo) XXX_Size() int {
	return xxx_messageInfo_GoPackageInfo.Size(m)
}
func (m *GoPackageInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_GoPackageInfo.DiscardUnknown(m)
}

var xxx_messageInfo_GoPackageInfo proto.InternalMessageInfo

func (m *GoPackageInfo) GetImportPath() string {
	if m != nil {
		return m.ImportPath
	}
	return ""
}

func init() {
	proto.RegisterType((*GoDetails)(nil), "kythe.proto.GoDetails")
	proto.RegisterType((*GoPackageInfo)(nil), "kythe.proto.GoPackageInfo")
}

func init() { proto.RegisterFile("kythe/proto/go.proto", fileDescriptor_go_0251d8e3856873be) }

var fileDescriptor_go_0251d8e3856873be = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xd0, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0xc6, 0x71, 0xea, 0x66, 0x5d, 0xdf, 0xe2, 0x25, 0x88, 0x04, 0x41, 0x56, 0x76, 0xea, 0xa9,
	0x13, 0xfc, 0x06, 0xa2, 0x0c, 0x6f, 0xa3, 0x78, 0xf2, 0x52, 0xd2, 0x34, 0xbe, 0x2d, 0x4b, 0xf7,
	0x94, 0x36, 0x0a, 0x7e, 0x3e, 0xbf, 0x98, 0x34, 0x71, 0xb2, 0x53, 0xde, 0xff, 0x2f, 0x04, 0x92,
	0xd0, 0xcd, 0xe1, 0xdb, 0xb5, 0x66, 0x3b, 0x8c, 0x70, 0xd8, 0x32, 0x0a, 0x3f, 0x88, 0xd4, 0x6b,
	0x88, 0xcd, 0x4f, 0x44, 0xc9, 0x0e, 0xcf, 0xc6, 0xa9, 0xce, 0x4e, 0x42, 0xd0, 0x92, 0x81, 0x49,
	0x46, 0x59, 0x94, 0x27, 0xa5, 0x9f, 0xc5, 0x2d, 0xc5, 0x0c, 0x35, 0xea, 0x56, 0x5e, 0x78, 0xfd,
	0xab, 0xe0, 0x23, 0xe0, 0xe4, 0xe2, 0xe4, 0x73, 0x05, 0x1f, 0x94, 0x6b, 0xe5, 0xf2, 0xe4, 0x73,
	0x89, 0x3b, 0x5a, 0x69, 0xf4, 0x43, 0x67, 0xcd, 0x28, 0x2f, 0xfd, 0xce, 0x7f, 0x8b, 0x7b, 0xa2,
	0xfa, 0xb3, 0xb3, 0x4d, 0xe5, 0x14, 0x4f, 0x32, 0xce, 0x16, 0x79, 0x52, 0x26, 0x5e, 0xde, 0x14,
	0x4f, 0x62, 0x4d, 0xa9, 0x66, 0x54, 0xe6, 0xa8, 0x6a, 0x6b, 0x1a, 0x79, 0x95, 0x45, 0xf9, 0xaa,
	0x24, 0xcd, 0x78, 0x09, 0xb2, 0x79, 0xa0, 0xeb, 0x1d, 0xf6, 0x4a, 0x1f, 0x14, 0x9b, 0xd7, 0xe3,
	0x07, 0xe6, 0x13, 0x5d, 0x3f, 0x60, 0x74, 0x95, 0xbf, 0x49, 0x78, 0x0f, 0x05, 0xda, 0x2b, 0xd7,
	0x3e, 0x15, 0xb4, 0xd6, 0xe8, 0x0b, 0x06, 0xd8, 0x9a, 0xa2, 0x31, 0x5f, 0x0e, 0xb0, 0x53, 0x71,
	0xf6, 0x35, 0xef, 0x29, 0xa3, 0x62, 0x54, 0x3e, 0xea, 0xd8, 0x2f, 0x8f, 0xbf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xc6, 0x29, 0xc8, 0xc3, 0x53, 0x01, 0x00, 0x00,
}
