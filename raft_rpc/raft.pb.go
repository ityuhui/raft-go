// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: raft_rpc/raft.proto

package raft_rpc

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type LogEntries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *LogEntries) Reset() {
	*x = LogEntries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_rpc_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntries) ProtoMessage() {}

func (x *LogEntries) ProtoReflect() protoreflect.Message {
	mi := &file_raft_rpc_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntries.ProtoReflect.Descriptor instead.
func (*LogEntries) Descriptor() ([]byte, []int) {
	return file_raft_rpc_raft_proto_rawDescGZIP(), []int{0}
}

func (x *LogEntries) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type AppendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeaderId     string        `protobuf:"bytes,1,opt,name=leaderId,proto3" json:"leaderId,omitempty"`
	Term         int64         `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	PrevLogIndex int64         `protobuf:"varint,3,opt,name=prevLogIndex,proto3" json:"prevLogIndex,omitempty"`
	PrevLogTerm  int64         `protobuf:"varint,4,opt,name=prevLogTerm,proto3" json:"prevLogTerm,omitempty"`
	LogEntries   []*LogEntries `protobuf:"bytes,5,rep,name=logEntries,proto3" json:"logEntries,omitempty"`
	LeaderCommit int64         `protobuf:"varint,6,opt,name=leaderCommit,proto3" json:"leaderCommit,omitempty"`
}

func (x *AppendRequest) Reset() {
	*x = AppendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_rpc_raft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendRequest) ProtoMessage() {}

func (x *AppendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_rpc_raft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendRequest.ProtoReflect.Descriptor instead.
func (*AppendRequest) Descriptor() ([]byte, []int) {
	return file_raft_rpc_raft_proto_rawDescGZIP(), []int{1}
}

func (x *AppendRequest) GetLeaderId() string {
	if x != nil {
		return x.LeaderId
	}
	return ""
}

func (x *AppendRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendRequest) GetPrevLogIndex() int64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendRequest) GetPrevLogTerm() int64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendRequest) GetLogEntries() []*LogEntries {
	if x != nil {
		return x.LogEntries
	}
	return nil
}

func (x *AppendRequest) GetLeaderCommit() int64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

type AppendReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term    int64  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AppendReply) Reset() {
	*x = AppendReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_rpc_raft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendReply) ProtoMessage() {}

func (x *AppendReply) ProtoReflect() protoreflect.Message {
	mi := &file_raft_rpc_raft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendReply.ProtoReflect.Descriptor instead.
func (*AppendReply) Descriptor() ([]byte, []int) {
	return file_raft_rpc_raft_proto_rawDescGZIP(), []int{2}
}

func (x *AppendReply) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AppendReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId  string `protobuf:"bytes,2,opt,name=candidateId,proto3" json:"candidateId,omitempty"`
	LastLogIndex int64  `protobuf:"varint,3,opt,name=lastLogIndex,proto3" json:"lastLogIndex,omitempty"`
	LastLogTerm  int64  `protobuf:"varint,4,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_rpc_raft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_rpc_raft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_raft_rpc_raft_proto_rawDescGZIP(), []int{3}
}

func (x *VoteRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteRequest) GetCandidateId() string {
	if x != nil {
		return x.CandidateId
	}
	return ""
}

func (x *VoteRequest) GetLastLogIndex() int64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *VoteRequest) GetLastLogTerm() int64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type VoteReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted bool  `protobuf:"varint,2,opt,name=voteGranted,proto3" json:"voteGranted,omitempty"`
}

func (x *VoteReply) Reset() {
	*x = VoteReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_rpc_raft_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteReply) ProtoMessage() {}

func (x *VoteReply) ProtoReflect() protoreflect.Message {
	mi := &file_raft_rpc_raft_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteReply.ProtoReflect.Descriptor instead.
func (*VoteReply) Descriptor() ([]byte, []int) {
	return file_raft_rpc_raft_proto_rawDescGZIP(), []int{4}
}

func (x *VoteReply) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteReply) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

var File_raft_rpc_raft_proto protoreflect.FileDescriptor

var file_raft_rpc_raft_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x72, 0x70, 0x63, 0x22,
	0x20, 0x0a, 0x0a, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72,
	0x6d, 0x22, 0xdf, 0x01, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c,
	0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c,
	0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x70, 0x72,
	0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x34, 0x0a, 0x0a, 0x6c, 0x6f, 0x67,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x72, 0x61, 0x66, 0x74, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x52, 0x0a, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12,
	0x22, 0x0a, 0x0c, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x22, 0x55, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x89, 0x01, 0x0a, 0x0b, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20,
	0x0a, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64,
	0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54,
	0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c,
	0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x41, 0x0a, 0x09, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x76, 0x6f,
	0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x32, 0x8d, 0x01, 0x0a, 0x0b, 0x52, 0x61,
	0x66, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x0d, 0x41, 0x70, 0x70,
	0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x17, 0x2e, 0x72, 0x61, 0x66,
	0x74, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x41,
	0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x15, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x13, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x56, 0x6f,
	0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_raft_rpc_raft_proto_rawDescOnce sync.Once
	file_raft_rpc_raft_proto_rawDescData = file_raft_rpc_raft_proto_rawDesc
)

func file_raft_rpc_raft_proto_rawDescGZIP() []byte {
	file_raft_rpc_raft_proto_rawDescOnce.Do(func() {
		file_raft_rpc_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_rpc_raft_proto_rawDescData)
	})
	return file_raft_rpc_raft_proto_rawDescData
}

var file_raft_rpc_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_raft_rpc_raft_proto_goTypes = []interface{}{
	(*LogEntries)(nil),    // 0: raft_rpc.LogEntries
	(*AppendRequest)(nil), // 1: raft_rpc.AppendRequest
	(*AppendReply)(nil),   // 2: raft_rpc.AppendReply
	(*VoteRequest)(nil),   // 3: raft_rpc.VoteRequest
	(*VoteReply)(nil),     // 4: raft_rpc.VoteReply
}
var file_raft_rpc_raft_proto_depIdxs = []int32{
	0, // 0: raft_rpc.AppendRequest.logEntries:type_name -> raft_rpc.LogEntries
	1, // 1: raft_rpc.RaftService.AppendEntries:input_type -> raft_rpc.AppendRequest
	3, // 2: raft_rpc.RaftService.RequestVote:input_type -> raft_rpc.VoteRequest
	2, // 3: raft_rpc.RaftService.AppendEntries:output_type -> raft_rpc.AppendReply
	4, // 4: raft_rpc.RaftService.RequestVote:output_type -> raft_rpc.VoteReply
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_raft_rpc_raft_proto_init() }
func file_raft_rpc_raft_proto_init() {
	if File_raft_rpc_raft_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_raft_rpc_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntries); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_rpc_raft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_rpc_raft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_rpc_raft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_rpc_raft_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_raft_rpc_raft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_raft_rpc_raft_proto_goTypes,
		DependencyIndexes: file_raft_rpc_raft_proto_depIdxs,
		MessageInfos:      file_raft_rpc_raft_proto_msgTypes,
	}.Build()
	File_raft_rpc_raft_proto = out.File
	file_raft_rpc_raft_proto_rawDesc = nil
	file_raft_rpc_raft_proto_goTypes = nil
	file_raft_rpc_raft_proto_depIdxs = nil
}
