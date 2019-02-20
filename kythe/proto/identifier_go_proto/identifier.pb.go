// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kythe/proto/identifier.proto

package identifier_go_proto

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

type FindRequest struct {
	Identifier           string   `protobuf:"bytes,1,opt,name=identifier,proto3" json:"identifier,omitempty"`
	Corpus               []string `protobuf:"bytes,2,rep,name=corpus,proto3" json:"corpus,omitempty"`
	Languages            []string `protobuf:"bytes,3,rep,name=languages,proto3" json:"languages,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindRequest) Reset()         { *m = FindRequest{} }
func (m *FindRequest) String() string { return proto.CompactTextString(m) }
func (*FindRequest) ProtoMessage()    {}
func (*FindRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_identifier_35694b3e713c8801, []int{0}
}
func (m *FindRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindRequest.Unmarshal(m, b)
}
func (m *FindRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindRequest.Marshal(b, m, deterministic)
}
func (dst *FindRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindRequest.Merge(dst, src)
}
func (m *FindRequest) XXX_Size() int {
	return xxx_messageInfo_FindRequest.Size(m)
}
func (m *FindRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindRequest proto.InternalMessageInfo

func (m *FindRequest) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

func (m *FindRequest) GetCorpus() []string {
	if m != nil {
		return m.Corpus
	}
	return nil
}

func (m *FindRequest) GetLanguages() []string {
	if m != nil {
		return m.Languages
	}
	return nil
}

type FindReply struct {
	Matches              []*FindReply_Match `protobuf:"bytes,1,rep,name=matches,proto3" json:"matches,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *FindReply) Reset()         { *m = FindReply{} }
func (m *FindReply) String() string { return proto.CompactTextString(m) }
func (*FindReply) ProtoMessage()    {}
func (*FindReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_identifier_35694b3e713c8801, []int{1}
}
func (m *FindReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindReply.Unmarshal(m, b)
}
func (m *FindReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindReply.Marshal(b, m, deterministic)
}
func (dst *FindReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindReply.Merge(dst, src)
}
func (m *FindReply) XXX_Size() int {
	return xxx_messageInfo_FindReply.Size(m)
}
func (m *FindReply) XXX_DiscardUnknown() {
	xxx_messageInfo_FindReply.DiscardUnknown(m)
}

var xxx_messageInfo_FindReply proto.InternalMessageInfo

func (m *FindReply) GetMatches() []*FindReply_Match {
	if m != nil {
		return m.Matches
	}
	return nil
}

type FindReply_Match struct {
	Ticket               string   `protobuf:"bytes,1,opt,name=ticket,proto3" json:"ticket,omitempty"`
	NodeKind             string   `protobuf:"bytes,2,opt,name=node_kind,json=nodeKind,proto3" json:"node_kind,omitempty"`
	NodeSubkind          string   `protobuf:"bytes,3,opt,name=node_subkind,json=nodeSubkind,proto3" json:"node_subkind,omitempty"`
	BaseName             string   `protobuf:"bytes,4,opt,name=base_name,json=baseName,proto3" json:"base_name,omitempty"`
	QualifiedName        string   `protobuf:"bytes,5,opt,name=qualified_name,json=qualifiedName,proto3" json:"qualified_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindReply_Match) Reset()         { *m = FindReply_Match{} }
func (m *FindReply_Match) String() string { return proto.CompactTextString(m) }
func (*FindReply_Match) ProtoMessage()    {}
func (*FindReply_Match) Descriptor() ([]byte, []int) {
	return fileDescriptor_identifier_35694b3e713c8801, []int{1, 0}
}
func (m *FindReply_Match) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindReply_Match.Unmarshal(m, b)
}
func (m *FindReply_Match) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindReply_Match.Marshal(b, m, deterministic)
}
func (dst *FindReply_Match) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindReply_Match.Merge(dst, src)
}
func (m *FindReply_Match) XXX_Size() int {
	return xxx_messageInfo_FindReply_Match.Size(m)
}
func (m *FindReply_Match) XXX_DiscardUnknown() {
	xxx_messageInfo_FindReply_Match.DiscardUnknown(m)
}

var xxx_messageInfo_FindReply_Match proto.InternalMessageInfo

func (m *FindReply_Match) GetTicket() string {
	if m != nil {
		return m.Ticket
	}
	return ""
}

func (m *FindReply_Match) GetNodeKind() string {
	if m != nil {
		return m.NodeKind
	}
	return ""
}

func (m *FindReply_Match) GetNodeSubkind() string {
	if m != nil {
		return m.NodeSubkind
	}
	return ""
}

func (m *FindReply_Match) GetBaseName() string {
	if m != nil {
		return m.BaseName
	}
	return ""
}

func (m *FindReply_Match) GetQualifiedName() string {
	if m != nil {
		return m.QualifiedName
	}
	return ""
}

func init() {
	proto.RegisterType((*FindRequest)(nil), "kythe.proto.FindRequest")
	proto.RegisterType((*FindReply)(nil), "kythe.proto.FindReply")
	proto.RegisterType((*FindReply_Match)(nil), "kythe.proto.FindReply.Match")
}

func init() {
	proto.RegisterFile("kythe/proto/identifier.proto", fileDescriptor_identifier_35694b3e713c8801)
}

var fileDescriptor_identifier_35694b3e713c8801 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x41, 0x6b, 0xf2, 0x30,
	0x18, 0xc7, 0xa9, 0x55, 0xdf, 0xb7, 0x4f, 0xb7, 0xc1, 0x32, 0x90, 0xe0, 0x64, 0x73, 0xc2, 0xc0,
	0x53, 0x05, 0x07, 0xb2, 0xf3, 0x0e, 0x83, 0x31, 0xdc, 0xa1, 0xde, 0x76, 0x29, 0xb1, 0x79, 0x56,
	0x83, 0x6d, 0xa2, 0x4d, 0x2a, 0xf8, 0x79, 0xf6, 0xe9, 0xf6, 0x2d, 0x46, 0x52, 0xa7, 0x3d, 0x78,
	0x2a, 0xcf, 0xef, 0xff, 0x0b, 0x7d, 0xf2, 0x0f, 0x0c, 0xd6, 0x7b, 0xb3, 0xc2, 0xc9, 0xa6, 0x54,
	0x46, 0x4d, 0x04, 0x47, 0x69, 0xc4, 0x97, 0xc0, 0x32, 0x72, 0x80, 0x84, 0x2e, 0xad, 0x87, 0x51,
	0x0a, 0xe1, 0xab, 0x90, 0x3c, 0xc6, 0x6d, 0x85, 0xda, 0x90, 0x3b, 0x80, 0x93, 0x4f, 0xbd, 0xa1,
	0x37, 0x0e, 0xe2, 0x06, 0x21, 0x3d, 0xe8, 0xa6, 0xaa, 0xdc, 0x54, 0x9a, 0xb6, 0x86, 0xfe, 0x38,
	0x88, 0x0f, 0x13, 0x19, 0x40, 0x90, 0x33, 0x99, 0x55, 0x2c, 0x43, 0x4d, 0x7d, 0x17, 0x9d, 0xc0,
	0xe8, 0xc7, 0x83, 0xa0, 0xfe, 0xcb, 0x26, 0xdf, 0x93, 0x19, 0xfc, 0x2b, 0x98, 0x49, 0x57, 0xa8,
	0xa9, 0x37, 0xf4, 0xc7, 0xe1, 0x74, 0x10, 0x35, 0x36, 0x8a, 0x8e, 0x62, 0x34, 0xb7, 0x56, 0xfc,
	0x27, 0xf7, 0xbf, 0x3d, 0xe8, 0x38, 0x64, 0xb7, 0x30, 0x22, 0x5d, 0xa3, 0x39, 0x6c, 0x78, 0x98,
	0xc8, 0x2d, 0x04, 0x52, 0x71, 0x4c, 0xd6, 0x42, 0x72, 0xda, 0x72, 0xd1, 0x7f, 0x0b, 0xde, 0x85,
	0xe4, 0xe4, 0x01, 0x2e, 0x5c, 0xa8, 0xab, 0xa5, 0xcb, 0x7d, 0x97, 0x87, 0x96, 0x2d, 0x6a, 0x64,
	0xcf, 0x2f, 0x99, 0xc6, 0x44, 0xb2, 0x02, 0x69, 0xbb, 0x3e, 0x6f, 0xc1, 0x07, 0x2b, 0x90, 0x3c,
	0xc2, 0xd5, 0xb6, 0x62, 0xb9, 0xed, 0x81, 0xd7, 0x46, 0xc7, 0x19, 0x97, 0x47, 0x6a, 0xb5, 0xe9,
	0x1c, 0xae, 0xdf, 0x8e, 0x7d, 0x2d, 0xb0, 0xdc, 0x89, 0x14, 0xc9, 0x33, 0xb4, 0xed, 0xb5, 0x08,
	0x3d, 0x73, 0x53, 0x57, 0x7c, 0xbf, 0x77, 0xbe, 0x83, 0x97, 0x19, 0xdc, 0xa7, 0xaa, 0x88, 0x32,
	0xa5, 0xb2, 0x1c, 0x23, 0x8e, 0x3b, 0xa3, 0x54, 0xae, 0x9b, 0xf2, 0xe7, 0xcd, 0xe9, 0x7d, 0x92,
	0x4c, 0x25, 0x0e, 0x2e, 0xbb, 0xee, 0xf3, 0xf4, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x7c, 0x37,
	0x1b, 0x0b, 0x02, 0x00, 0x00,
}
