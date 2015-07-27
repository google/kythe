// Code generated by protoc-gen-go.
// source: kythe/proto/storage.proto
// DO NOT EDIT!

/*
Package storage_proto is a generated protocol buffer package.

It is generated from these files:
	kythe/proto/storage.proto

It has these top-level messages:
	VName
	VNameMask
	Entry
	Entries
	ReadRequest
	WriteRequest
	WriteReply
	ScanRequest
	CountRequest
	CountReply
	ShardRequest
	SearchRequest
	SearchReply
*/
package storage_proto

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

// VName is a proto representation of a vector name.
//
// Rules:
//  - All fields must be optional, and must have default values.
//  - No field may ever be removed.  If a field is deprecated, it may be
//    renamed or marked with a comment, but must not be deleted.
//  - New fields are always added to the end of the message.
//  - All fields must be strings, not messages.
//
// One of the key principles is that we want as few fields as possible in a
// vname.  We're not trying to exhaust the possible dimensions along which a
// name could vary, but to find a minimal basis. Be conservative.
type VName struct {
	// A language-specific signature assigned by the analyzer.
	// e.g., "com.google.common.collect.Lists.newLinkedList<#1>()"
	Signature string `protobuf:"bytes,1,opt,name=signature" json:"signature,omitempty"`
	// The corpus this name belongs to.
	// e.g., "kythe", "chromium", "github.com/creachadair/imath", "aosp"
	// The corpus label "kythe" is reserved for internal use.
	Corpus string `protobuf:"bytes,2,opt,name=corpus" json:"corpus,omitempty"`
	// A corpus-specific root label, designating a subordinate collection within
	// the corpus.  If a corpus stores files in unrelated directory structures,
	// for example, the root can be used to distinguish them.  Or, of a corpus
	// incorporates subprojects, the root can be a project ID that it governs.
	// This may also be used to distinguish virtual subgroups of a corpus such as
	// generated files.
	Root string `protobuf:"bytes,3,opt,name=root" json:"root,omitempty"`
	// A path-structured label describing the location of this object relative to
	// the corpus and the root.  For code, this will generally be the relative
	// path to the file containing the code, e.g., "storage/service.go" in kythe.
	//
	// However, this need not be a true file path; virtual objects like figments
	// can assign an ad-hoc abstract ID, or omit it entirely.
	//
	// Examples:
	//   "devools/kythe/platform/go/datastore.go" (a file)
	//   "type/cpp/void.cc" (a type figment)
	Path string `protobuf:"bytes,4,opt,name=path" json:"path,omitempty"`
	// The language this name belongs to.
	// e.g., "c++", "python", "elisp", "haskell", "java"
	//
	// The schema will define specific labels for each supported language, so we
	// don't wind up with a confusion of names like "cxx", "cpp", "C++", etc.
	// Prototype: Official language name converted to lowercase.  If a version
	// number is necessary, include it, e.g., "python3".
	Language string `protobuf:"bytes,5,opt,name=language" json:"language,omitempty"`
}

func (m *VName) Reset()         { *m = VName{} }
func (m *VName) String() string { return proto.CompactTextString(m) }
func (*VName) ProtoMessage()    {}

type VNameMask struct {
	Signature bool `protobuf:"varint,1,opt,name=signature" json:"signature,omitempty"`
	Corpus    bool `protobuf:"varint,2,opt,name=corpus" json:"corpus,omitempty"`
	Root      bool `protobuf:"varint,3,opt,name=root" json:"root,omitempty"`
	Path      bool `protobuf:"varint,4,opt,name=path" json:"path,omitempty"`
	Language  bool `protobuf:"varint,5,opt,name=language" json:"language,omitempty"`
}

func (m *VNameMask) Reset()         { *m = VNameMask{} }
func (m *VNameMask) String() string { return proto.CompactTextString(m) }
func (*VNameMask) ProtoMessage()    {}

// An Entry associates a fact with a graph object (node or edge).  This is the
// the primary unit of storage.
type Entry struct {
	Source *VName `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	// The following two fields must either be both empty, or both nonempty.
	EdgeKind string `protobuf:"bytes,2,opt,name=edge_kind" json:"edge_kind,omitempty"`
	Target   *VName `protobuf:"bytes,3,opt,name=target" json:"target,omitempty"`
	// The grammar for fact_name:
	//  name   = "/" | 1*path
	//  path   = "/" word
	//  word   = 1*{LETTER|DIGIT|PUNCT}
	//  LETTER = [A-Za-z]
	//  DIGIT  = [0-9]
	//  PUNCT  = [-.@#$%&_+:()]
	FactName  string `protobuf:"bytes,4,opt,name=fact_name" json:"fact_name,omitempty"`
	FactValue []byte `protobuf:"bytes,5,opt,name=fact_value,proto3" json:"fact_value,omitempty"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}

func (m *Entry) GetSource() *VName {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Entry) GetTarget() *VName {
	if m != nil {
		return m.Target
	}
	return nil
}

// A collection of Entry instances.
type Entries struct {
	Entries []*Entry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *Entries) Reset()         { *m = Entries{} }
func (m *Entries) String() string { return proto.CompactTextString(m) }
func (*Entries) ProtoMessage()    {}

func (m *Entries) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

// Request for a stream of Entry objects from a GraphStore.  Read operations
// should be implemented with time complexity proportional to the size of the
// return set.
type ReadRequest struct {
	// Return entries having this source VName, which may not be empty.
	Source *VName `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	// Return entries having this edge kind; if empty, only entries with an empty
	// edge kind are returned; if "*", entries of any edge kind are returned.
	EdgeKind string `protobuf:"bytes,2,opt,name=edge_kind" json:"edge_kind,omitempty"`
}

func (m *ReadRequest) Reset()         { *m = ReadRequest{} }
func (m *ReadRequest) String() string { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()    {}

func (m *ReadRequest) GetSource() *VName {
	if m != nil {
		return m.Source
	}
	return nil
}

// Request to write Entry objects to a GraphStore
type WriteRequest struct {
	Source *VName                 `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	Update []*WriteRequest_Update `protobuf:"bytes,2,rep,name=update" json:"update,omitempty"`
}

func (m *WriteRequest) Reset()         { *m = WriteRequest{} }
func (m *WriteRequest) String() string { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()    {}

func (m *WriteRequest) GetSource() *VName {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *WriteRequest) GetUpdate() []*WriteRequest_Update {
	if m != nil {
		return m.Update
	}
	return nil
}

type WriteRequest_Update struct {
	EdgeKind  string `protobuf:"bytes,1,opt,name=edge_kind" json:"edge_kind,omitempty"`
	Target    *VName `protobuf:"bytes,2,opt,name=target" json:"target,omitempty"`
	FactName  string `protobuf:"bytes,3,opt,name=fact_name" json:"fact_name,omitempty"`
	FactValue []byte `protobuf:"bytes,4,opt,name=fact_value,proto3" json:"fact_value,omitempty"`
}

func (m *WriteRequest_Update) Reset()         { *m = WriteRequest_Update{} }
func (m *WriteRequest_Update) String() string { return proto.CompactTextString(m) }
func (*WriteRequest_Update) ProtoMessage()    {}

func (m *WriteRequest_Update) GetTarget() *VName {
	if m != nil {
		return m.Target
	}
	return nil
}

// Response to a WriteRequest
type WriteReply struct {
}

func (m *WriteReply) Reset()         { *m = WriteReply{} }
func (m *WriteReply) String() string { return proto.CompactTextString(m) }
func (*WriteReply) ProtoMessage()    {}

// Request for a stream of Entry objects resulting from a full scan of a
// GraphStore.
type ScanRequest struct {
	// Return entries having this target VName; if empty, any target field is
	// matched, including empty.
	Target *VName `protobuf:"bytes,1,opt,name=target" json:"target,omitempty"`
	// Return entries having this kind; if empty, any kind is matched, including
	// empty.
	EdgeKind string `protobuf:"bytes,2,opt,name=edge_kind" json:"edge_kind,omitempty"`
	// Return entries having fact labels with this prefix; if empty, any fact
	// label is matched,
	FactPrefix string `protobuf:"bytes,3,opt,name=fact_prefix" json:"fact_prefix,omitempty"`
}

func (m *ScanRequest) Reset()         { *m = ScanRequest{} }
func (m *ScanRequest) String() string { return proto.CompactTextString(m) }
func (*ScanRequest) ProtoMessage()    {}

func (m *ScanRequest) GetTarget() *VName {
	if m != nil {
		return m.Target
	}
	return nil
}

// Request for the size of the shard at the given index.
type CountRequest struct {
	Index  int64 `protobuf:"varint,1,opt,name=index" json:"index,omitempty"`
	Shards int64 `protobuf:"varint,2,opt,name=shards" json:"shards,omitempty"`
}

func (m *CountRequest) Reset()         { *m = CountRequest{} }
func (m *CountRequest) String() string { return proto.CompactTextString(m) }
func (*CountRequest) ProtoMessage()    {}

// Response for a CountRequest
type CountReply struct {
	// Total number of entries in the specified shard.
	Entries int64 `protobuf:"varint,1,opt,name=entries" json:"entries,omitempty"`
}

func (m *CountReply) Reset()         { *m = CountReply{} }
func (m *CountReply) String() string { return proto.CompactTextString(m) }
func (*CountReply) ProtoMessage()    {}

// Request for a stream of Entry objects in the given shard.
type ShardRequest struct {
	Index  int64 `protobuf:"varint,1,opt,name=index" json:"index,omitempty"`
	Shards int64 `protobuf:"varint,2,opt,name=shards" json:"shards,omitempty"`
}

func (m *ShardRequest) Reset()         { *m = ShardRequest{} }
func (m *ShardRequest) String() string { return proto.CompactTextString(m) }
func (*ShardRequest) ProtoMessage()    {}

// Request for the set of node tickets matching a partial VName and collection
// of known facts.
type SearchRequest struct {
	// Partial VName to match against nodes.  Each non-empty field becomes a
	// constraint (i.e. the signature/corpus/etc. must be exactly case-sensitively
	// equal to the given string) on the set of returned nodes.  Exact matching
	// turns into prefix matching if the corresponding field in partial_prefix is
	// set to true.
	Partial *VName `protobuf:"bytes,1,opt,name=partial" json:"partial,omitempty"`
	// Facts that a node must have to be matched.  Exact matching turns into
	// prefix matching if a Fact has its prefix field set to true.
	Fact []*SearchRequest_Fact `protobuf:"bytes,2,rep,name=fact" json:"fact,omitempty"`
	// Setting any field in this mask to true converts exact value matching to
	// prefix value matching for the corresponding VName component in partial.
	PartialPrefix *VNameMask `protobuf:"bytes,3,opt,name=partial_prefix" json:"partial_prefix,omitempty"`
}

func (m *SearchRequest) Reset()         { *m = SearchRequest{} }
func (m *SearchRequest) String() string { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()    {}

func (m *SearchRequest) GetPartial() *VName {
	if m != nil {
		return m.Partial
	}
	return nil
}

func (m *SearchRequest) GetFact() []*SearchRequest_Fact {
	if m != nil {
		return m.Fact
	}
	return nil
}

func (m *SearchRequest) GetPartialPrefix() *VNameMask {
	if m != nil {
		return m.PartialPrefix
	}
	return nil
}

type SearchRequest_Fact struct {
	Name   string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value  []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Prefix bool   `protobuf:"varint,3,opt,name=prefix" json:"prefix,omitempty"`
}

func (m *SearchRequest_Fact) Reset()         { *m = SearchRequest_Fact{} }
func (m *SearchRequest_Fact) String() string { return proto.CompactTextString(m) }
func (*SearchRequest_Fact) ProtoMessage()    {}

// Response for a SearchRequest.
type SearchReply struct {
	// Set of node tickets matching a given SearchRequest.
	Ticket []string `protobuf:"bytes,1,rep,name=ticket" json:"ticket,omitempty"`
}

func (m *SearchReply) Reset()         { *m = SearchReply{} }
func (m *SearchReply) String() string { return proto.CompactTextString(m) }
func (*SearchReply) ProtoMessage()    {}

// Client API for GraphStore service

type GraphStoreClient interface {
	// Read responds with all Entry messages that match the given ReadRequest.
	// The Read operation should be implemented with time complexity proportional
	// to the size of the return set.
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (GraphStore_ReadClient, error)
	// Scan responds with all Entry messages matching the given ScanRequest.  If a
	// ScanRequest field is empty, any entry value for that field matches and will
	// be returned.  Scan is similar to Read, but with no time complexity
	// restrictions.
	Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (GraphStore_ScanClient, error)
	// Write atomically inserts or updates a collection of entries into the store.
	// Each update is a tuple of the form (kind, target, fact, value).  For each
	// such update, an entry (source, kind, target, fact, value) is written into
	// the store, replacing any existing entry (source, kind, target, fact,
	// value') that may exist.  Note that this operation cannot delete any data
	// from the store; entries are only ever inserted or updated.  Apart from
	// acting atomically, no other constraints are placed on the implementation.
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteReply, error)
}

type graphStoreClient struct {
	cc *grpc.ClientConn
}

func NewGraphStoreClient(cc *grpc.ClientConn) GraphStoreClient {
	return &graphStoreClient{cc}
}

func (c *graphStoreClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (GraphStore_ReadClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GraphStore_serviceDesc.Streams[0], c.cc, "/kythe.proto.GraphStore/Read", opts...)
	if err != nil {
		return nil, err
	}
	x := &graphStoreReadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GraphStore_ReadClient interface {
	Recv() (*Entry, error)
	grpc.ClientStream
}

type graphStoreReadClient struct {
	grpc.ClientStream
}

func (x *graphStoreReadClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *graphStoreClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (GraphStore_ScanClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GraphStore_serviceDesc.Streams[1], c.cc, "/kythe.proto.GraphStore/Scan", opts...)
	if err != nil {
		return nil, err
	}
	x := &graphStoreScanClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GraphStore_ScanClient interface {
	Recv() (*Entry, error)
	grpc.ClientStream
}

type graphStoreScanClient struct {
	grpc.ClientStream
}

func (x *graphStoreScanClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *graphStoreClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteReply, error) {
	out := new(WriteReply)
	err := grpc.Invoke(ctx, "/kythe.proto.GraphStore/Write", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GraphStore service

type GraphStoreServer interface {
	// Read responds with all Entry messages that match the given ReadRequest.
	// The Read operation should be implemented with time complexity proportional
	// to the size of the return set.
	Read(*ReadRequest, GraphStore_ReadServer) error
	// Scan responds with all Entry messages matching the given ScanRequest.  If a
	// ScanRequest field is empty, any entry value for that field matches and will
	// be returned.  Scan is similar to Read, but with no time complexity
	// restrictions.
	Scan(*ScanRequest, GraphStore_ScanServer) error
	// Write atomically inserts or updates a collection of entries into the store.
	// Each update is a tuple of the form (kind, target, fact, value).  For each
	// such update, an entry (source, kind, target, fact, value) is written into
	// the store, replacing any existing entry (source, kind, target, fact,
	// value') that may exist.  Note that this operation cannot delete any data
	// from the store; entries are only ever inserted or updated.  Apart from
	// acting atomically, no other constraints are placed on the implementation.
	Write(context.Context, *WriteRequest) (*WriteReply, error)
}

func RegisterGraphStoreServer(s *grpc.Server, srv GraphStoreServer) {
	s.RegisterService(&_GraphStore_serviceDesc, srv)
}

func _GraphStore_Read_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GraphStoreServer).Read(m, &graphStoreReadServer{stream})
}

type GraphStore_ReadServer interface {
	Send(*Entry) error
	grpc.ServerStream
}

type graphStoreReadServer struct {
	grpc.ServerStream
}

func (x *graphStoreReadServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func _GraphStore_Scan_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ScanRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GraphStoreServer).Scan(m, &graphStoreScanServer{stream})
}

type GraphStore_ScanServer interface {
	Send(*Entry) error
	grpc.ServerStream
}

type graphStoreScanServer struct {
	grpc.ServerStream
}

func (x *graphStoreScanServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

func _GraphStore_Write_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(WriteRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(GraphStoreServer).Write(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _GraphStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kythe.proto.GraphStore",
	HandlerType: (*GraphStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _GraphStore_Write_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Read",
			Handler:       _GraphStore_Read_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Scan",
			Handler:       _GraphStore_Scan_Handler,
			ServerStreams: true,
		},
	},
}

// Client API for ShardedGraphStore service

type ShardedGraphStoreClient interface {
	// Count returns the number of entries in the given shard.
	Count(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*CountReply, error)
	// Shard responds with each Entry in the given shard.
	Shard(ctx context.Context, in *ShardRequest, opts ...grpc.CallOption) (ShardedGraphStore_ShardClient, error)
}

type shardedGraphStoreClient struct {
	cc *grpc.ClientConn
}

func NewShardedGraphStoreClient(cc *grpc.ClientConn) ShardedGraphStoreClient {
	return &shardedGraphStoreClient{cc}
}

func (c *shardedGraphStoreClient) Count(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*CountReply, error) {
	out := new(CountReply)
	err := grpc.Invoke(ctx, "/kythe.proto.ShardedGraphStore/Count", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shardedGraphStoreClient) Shard(ctx context.Context, in *ShardRequest, opts ...grpc.CallOption) (ShardedGraphStore_ShardClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ShardedGraphStore_serviceDesc.Streams[0], c.cc, "/kythe.proto.ShardedGraphStore/Shard", opts...)
	if err != nil {
		return nil, err
	}
	x := &shardedGraphStoreShardClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ShardedGraphStore_ShardClient interface {
	Recv() (*Entry, error)
	grpc.ClientStream
}

type shardedGraphStoreShardClient struct {
	grpc.ClientStream
}

func (x *shardedGraphStoreShardClient) Recv() (*Entry, error) {
	m := new(Entry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for ShardedGraphStore service

type ShardedGraphStoreServer interface {
	// Count returns the number of entries in the given shard.
	Count(context.Context, *CountRequest) (*CountReply, error)
	// Shard responds with each Entry in the given shard.
	Shard(*ShardRequest, ShardedGraphStore_ShardServer) error
}

func RegisterShardedGraphStoreServer(s *grpc.Server, srv ShardedGraphStoreServer) {
	s.RegisterService(&_ShardedGraphStore_serviceDesc, srv)
}

func _ShardedGraphStore_Count_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(CountRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(ShardedGraphStoreServer).Count(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _ShardedGraphStore_Shard_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ShardRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ShardedGraphStoreServer).Shard(m, &shardedGraphStoreShardServer{stream})
}

type ShardedGraphStore_ShardServer interface {
	Send(*Entry) error
	grpc.ServerStream
}

type shardedGraphStoreShardServer struct {
	grpc.ServerStream
}

func (x *shardedGraphStoreShardServer) Send(m *Entry) error {
	return x.ServerStream.SendMsg(m)
}

var _ShardedGraphStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kythe.proto.ShardedGraphStore",
	HandlerType: (*ShardedGraphStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Count",
			Handler:    _ShardedGraphStore_Count_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Shard",
			Handler:       _ShardedGraphStore_Shard_Handler,
			ServerStreams: true,
		},
	},
}

// Client API for SearchService service

type SearchServiceClient interface {
	// Search responds with the set of node tickets that match the given
	// SearchRequest.
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchReply, error)
}

type searchServiceClient struct {
	cc *grpc.ClientConn
}

func NewSearchServiceClient(cc *grpc.ClientConn) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchReply, error) {
	out := new(SearchReply)
	err := grpc.Invoke(ctx, "/kythe.proto.SearchService/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SearchService service

type SearchServiceServer interface {
	// Search responds with the set of node tickets that match the given
	// SearchRequest.
	Search(context.Context, *SearchRequest) (*SearchReply, error)
}

func RegisterSearchServiceServer(s *grpc.Server, srv SearchServiceServer) {
	s.RegisterService(&_SearchService_serviceDesc, srv)
}

func _SearchService_Search_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(SearchRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(SearchServiceServer).Search(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _SearchService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kythe.proto.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _SearchService_Search_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
