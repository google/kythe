// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kythe/go/util/riegeli/riegeli_test.proto

package riegeli_test_go_proto

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

type Simple struct {
	Name                 *string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Simple) Reset()         { *m = Simple{} }
func (m *Simple) String() string { return proto.CompactTextString(m) }
func (*Simple) ProtoMessage()    {}
func (*Simple) Descriptor() ([]byte, []int) {
	return fileDescriptor_riegeli_test_b1b33b2af4da371c, []int{0}
}
func (m *Simple) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Simple.Unmarshal(m, b)
}
func (m *Simple) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Simple.Marshal(b, m, deterministic)
}
func (dst *Simple) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Simple.Merge(dst, src)
}
func (m *Simple) XXX_Size() int {
	return xxx_messageInfo_Simple.Size(m)
}
func (m *Simple) XXX_DiscardUnknown() {
	xxx_messageInfo_Simple.DiscardUnknown(m)
}

var xxx_messageInfo_Simple proto.InternalMessageInfo

func (m *Simple) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type Complex struct {
	Str                  *string          `protobuf:"bytes,1,opt,name=str" json:"str,omitempty"`
	I32                  *int32           `protobuf:"varint,2,opt,name=i32" json:"i32,omitempty"`
	I64                  *int64           `protobuf:"varint,3,opt,name=i64" json:"i64,omitempty"`
	Bits                 []byte           `protobuf:"bytes,4,opt,name=bits" json:"bits,omitempty"`
	Rep                  []string         `protobuf:"bytes,5,rep,name=rep" json:"rep,omitempty"`
	SimpleNested         *Simple          `protobuf:"bytes,6,opt,name=simple_nested,json=simpleNested" json:"simple_nested,omitempty"`
	Group                []*Complex_Group `protobuf:"group,7,rep,name=Group,json=group" json:"group,omitempty"`
	ComplexNested        *Complex         `protobuf:"bytes,8,opt,name=complex_nested,json=complexNested" json:"complex_nested,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Complex) Reset()         { *m = Complex{} }
func (m *Complex) String() string { return proto.CompactTextString(m) }
func (*Complex) ProtoMessage()    {}
func (*Complex) Descriptor() ([]byte, []int) {
	return fileDescriptor_riegeli_test_b1b33b2af4da371c, []int{1}
}
func (m *Complex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Complex.Unmarshal(m, b)
}
func (m *Complex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Complex.Marshal(b, m, deterministic)
}
func (dst *Complex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Complex.Merge(dst, src)
}
func (m *Complex) XXX_Size() int {
	return xxx_messageInfo_Complex.Size(m)
}
func (m *Complex) XXX_DiscardUnknown() {
	xxx_messageInfo_Complex.DiscardUnknown(m)
}

var xxx_messageInfo_Complex proto.InternalMessageInfo

func (m *Complex) GetStr() string {
	if m != nil && m.Str != nil {
		return *m.Str
	}
	return ""
}

func (m *Complex) GetI32() int32 {
	if m != nil && m.I32 != nil {
		return *m.I32
	}
	return 0
}

func (m *Complex) GetI64() int64 {
	if m != nil && m.I64 != nil {
		return *m.I64
	}
	return 0
}

func (m *Complex) GetBits() []byte {
	if m != nil {
		return m.Bits
	}
	return nil
}

func (m *Complex) GetRep() []string {
	if m != nil {
		return m.Rep
	}
	return nil
}

func (m *Complex) GetSimpleNested() *Simple {
	if m != nil {
		return m.SimpleNested
	}
	return nil
}

func (m *Complex) GetGroup() []*Complex_Group {
	if m != nil {
		return m.Group
	}
	return nil
}

func (m *Complex) GetComplexNested() *Complex {
	if m != nil {
		return m.ComplexNested
	}
	return nil
}

type Complex_Group struct {
	GrpStr               *string  `protobuf:"bytes,1,opt,name=grp_str,json=grpStr" json:"grp_str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Complex_Group) Reset()         { *m = Complex_Group{} }
func (m *Complex_Group) String() string { return proto.CompactTextString(m) }
func (*Complex_Group) ProtoMessage()    {}
func (*Complex_Group) Descriptor() ([]byte, []int) {
	return fileDescriptor_riegeli_test_b1b33b2af4da371c, []int{1, 0}
}
func (m *Complex_Group) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Complex_Group.Unmarshal(m, b)
}
func (m *Complex_Group) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Complex_Group.Marshal(b, m, deterministic)
}
func (dst *Complex_Group) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Complex_Group.Merge(dst, src)
}
func (m *Complex_Group) XXX_Size() int {
	return xxx_messageInfo_Complex_Group.Size(m)
}
func (m *Complex_Group) XXX_DiscardUnknown() {
	xxx_messageInfo_Complex_Group.DiscardUnknown(m)
}

var xxx_messageInfo_Complex_Group proto.InternalMessageInfo

func (m *Complex_Group) GetGrpStr() string {
	if m != nil && m.GrpStr != nil {
		return *m.GrpStr
	}
	return ""
}

func init() {
	proto.RegisterType((*Simple)(nil), "kythe.proto.riegeli_test.Simple")
	proto.RegisterType((*Complex)(nil), "kythe.proto.riegeli_test.Complex")
	proto.RegisterType((*Complex_Group)(nil), "kythe.proto.riegeli_test.Complex.Group")
}

func init() {
	proto.RegisterFile("kythe/go/util/riegeli/riegeli_test.proto", fileDescriptor_riegeli_test_b1b33b2af4da371c)
}

var fileDescriptor_riegeli_test_b1b33b2af4da371c = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xcd, 0x4e, 0x03, 0x21,
	0x14, 0x85, 0x43, 0xa7, 0xd3, 0xd1, 0x6b, 0x6b, 0x0c, 0x1b, 0x89, 0x71, 0x81, 0xdd, 0xc8, 0x8a,
	0x49, 0x46, 0xe3, 0xce, 0x95, 0x31, 0xba, 0x72, 0x41, 0x1f, 0x60, 0xe2, 0xcf, 0x0d, 0x12, 0xa7,
	0x85, 0x00, 0x4d, 0xf4, 0x49, 0x7d, 0x1d, 0x03, 0x4c, 0xa3, 0x9b, 0xc6, 0x15, 0x87, 0x7b, 0xcf,
	0x39, 0xf9, 0x00, 0xc4, 0xc7, 0x57, 0x7c, 0xc7, 0x56, 0xdb, 0x76, 0x1b, 0xcd, 0xd0, 0x7a, 0x83,
	0x1a, 0x07, 0xb3, 0x3b, 0xfb, 0x88, 0x21, 0x4a, 0xe7, 0x6d, 0xb4, 0x94, 0x65, 0x67, 0xb9, 0xc8,
	0xbf, 0xfb, 0xe5, 0x39, 0xcc, 0x56, 0x66, 0xed, 0x06, 0xa4, 0x14, 0xa6, 0x9b, 0xe7, 0x35, 0x32,
	0xc2, 0x89, 0x38, 0x54, 0x59, 0x2f, 0xbf, 0x27, 0xd0, 0xdc, 0xd9, 0xb4, 0xfe, 0xa4, 0x27, 0x50,
	0x85, 0xe8, 0xc7, 0x75, 0x92, 0x69, 0x62, 0xae, 0x3a, 0x36, 0xe1, 0x44, 0xd4, 0x2a, 0xc9, 0x3c,
	0xb9, 0xb9, 0x66, 0x15, 0x27, 0xa2, 0x52, 0x49, 0xa6, 0xd6, 0x17, 0x13, 0x03, 0x9b, 0x72, 0x22,
	0xe6, 0x2a, 0xeb, 0xe4, 0xf2, 0xe8, 0x58, 0xcd, 0xab, 0xd4, 0xe4, 0xd1, 0xd1, 0x7b, 0x58, 0x84,
	0x4c, 0xd1, 0x6f, 0x30, 0x44, 0x7c, 0x63, 0x33, 0x4e, 0xc4, 0x51, 0xc7, 0xe5, 0x3e, 0x6e, 0x59,
	0xa0, 0xd5, 0xbc, 0xc4, 0x9e, 0x72, 0x8a, 0xde, 0x42, 0xad, 0xbd, 0xdd, 0x3a, 0xd6, 0xf0, 0x4a,
	0x40, 0x77, 0xb9, 0x3f, 0x3e, 0x3e, 0x4a, 0x3e, 0x24, 0xbb, 0x2a, 0x29, 0xfa, 0x08, 0xc7, 0xaf,
	0x65, 0xbe, 0xc3, 0x38, 0xc8, 0x18, 0x17, 0xff, 0xf6, 0xa8, 0xc5, 0x18, 0x2c, 0x20, 0x67, 0x1c,
	0xea, 0xdc, 0x4c, 0x4f, 0xa1, 0xd1, 0xde, 0xf5, 0xbf, 0x1f, 0x37, 0xd3, 0xde, 0xad, 0xa2, 0xff,
	0x09, 0x00, 0x00, 0xff, 0xff, 0x16, 0x1b, 0xcf, 0x85, 0xbc, 0x01, 0x00, 0x00,
}
