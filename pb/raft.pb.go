// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RequestVoteRequest struct {
	Term                 uint64   `protobuf:"varint,11,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId          string   `protobuf:"bytes,12,opt,name=candidate_id,json=candidateId,proto3" json:"candidate_id,omitempty"`
	LastLogIndex         uint64   `protobuf:"varint,13,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"`
	LastLogTerm          uint64   `protobuf:"varint,14,opt,name=last_log_term,json=lastLogTerm,proto3" json:"last_log_term,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteRequest) Reset()         { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()    {}
func (*RequestVoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{0}
}

func (m *RequestVoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteRequest.Unmarshal(m, b)
}
func (m *RequestVoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteRequest.Marshal(b, m, deterministic)
}
func (m *RequestVoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteRequest.Merge(m, src)
}
func (m *RequestVoteRequest) XXX_Size() int {
	return xxx_messageInfo_RequestVoteRequest.Size(m)
}
func (m *RequestVoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteRequest proto.InternalMessageInfo

func (m *RequestVoteRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteRequest) GetCandidateId() string {
	if m != nil {
		return m.CandidateId
	}
	return ""
}

func (m *RequestVoteRequest) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogTerm() uint64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

type RequestVoteResponse struct {
	Term                 uint64   `protobuf:"varint,11,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted          bool     `protobuf:"varint,12,opt,name=vote_granted,json=voteGranted,proto3" json:"vote_granted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteResponse) Reset()         { *m = RequestVoteResponse{} }
func (m *RequestVoteResponse) String() string { return proto.CompactTextString(m) }
func (*RequestVoteResponse) ProtoMessage()    {}
func (*RequestVoteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{1}
}

func (m *RequestVoteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteResponse.Unmarshal(m, b)
}
func (m *RequestVoteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteResponse.Marshal(b, m, deterministic)
}
func (m *RequestVoteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteResponse.Merge(m, src)
}
func (m *RequestVoteResponse) XXX_Size() int {
	return xxx_messageInfo_RequestVoteResponse.Size(m)
}
func (m *RequestVoteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteResponse proto.InternalMessageInfo

func (m *RequestVoteResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteResponse) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

type AppendEntriesRequest struct {
	Term                 uint64          `protobuf:"varint,11,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId             string          `protobuf:"bytes,12,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	PrevLogIndex         uint64          `protobuf:"varint,13,opt,name=prev_log_index,json=prevLogIndex,proto3" json:"prev_log_index,omitempty"`
	PrevLogTerm          uint64          `protobuf:"varint,14,opt,name=prev_log_term,json=prevLogTerm,proto3" json:"prev_log_term,omitempty"`
	Entries              []*CommandEntry `protobuf:"bytes,15,rep,name=entries,proto3" json:"entries,omitempty"`
	LeaderCommit         uint64          `protobuf:"varint,16,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AppendEntriesRequest) Reset()         { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()    {}
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{2}
}

func (m *AppendEntriesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesRequest.Unmarshal(m, b)
}
func (m *AppendEntriesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesRequest.Marshal(b, m, deterministic)
}
func (m *AppendEntriesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesRequest.Merge(m, src)
}
func (m *AppendEntriesRequest) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesRequest.Size(m)
}
func (m *AppendEntriesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesRequest proto.InternalMessageInfo

func (m *AppendEntriesRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}

func (m *AppendEntriesRequest) GetPrevLogIndex() uint64 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogTerm() uint64 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesRequest) GetEntries() []*CommandEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *AppendEntriesRequest) GetLeaderCommit() uint64 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

type AppendEntriesResponse struct {
	Term                 uint64   `protobuf:"varint,11,opt,name=term,proto3" json:"term,omitempty"`
	Success              bool     `protobuf:"varint,12,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AppendEntriesResponse) Reset()         { *m = AppendEntriesResponse{} }
func (m *AppendEntriesResponse) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesResponse) ProtoMessage()    {}
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{3}
}

func (m *AppendEntriesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesResponse.Unmarshal(m, b)
}
func (m *AppendEntriesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesResponse.Marshal(b, m, deterministic)
}
func (m *AppendEntriesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesResponse.Merge(m, src)
}
func (m *AppendEntriesResponse) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesResponse.Size(m)
}
func (m *AppendEntriesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesResponse proto.InternalMessageInfo

func (m *AppendEntriesResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type CommandEntry struct {
	Command              string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Args                 []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandEntry) Reset()         { *m = CommandEntry{} }
func (m *CommandEntry) String() string { return proto.CompactTextString(m) }
func (*CommandEntry) ProtoMessage()    {}
func (*CommandEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{4}
}

func (m *CommandEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandEntry.Unmarshal(m, b)
}
func (m *CommandEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandEntry.Marshal(b, m, deterministic)
}
func (m *CommandEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandEntry.Merge(m, src)
}
func (m *CommandEntry) XXX_Size() int {
	return xxx_messageInfo_CommandEntry.Size(m)
}
func (m *CommandEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandEntry.DiscardUnknown(m)
}

var xxx_messageInfo_CommandEntry proto.InternalMessageInfo

func (m *CommandEntry) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *CommandEntry) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

type CommandRequest struct {
	Cid                  string        `protobuf:"bytes,1,opt,name=cid,proto3" json:"cid,omitempty"`
	Entry                *CommandEntry `protobuf:"bytes,11,opt,name=entry,proto3" json:"entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CommandRequest) Reset()         { *m = CommandRequest{} }
func (m *CommandRequest) String() string { return proto.CompactTextString(m) }
func (*CommandRequest) ProtoMessage()    {}
func (*CommandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{5}
}

func (m *CommandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandRequest.Unmarshal(m, b)
}
func (m *CommandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandRequest.Marshal(b, m, deterministic)
}
func (m *CommandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandRequest.Merge(m, src)
}
func (m *CommandRequest) XXX_Size() int {
	return xxx_messageInfo_CommandRequest.Size(m)
}
func (m *CommandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommandRequest proto.InternalMessageInfo

func (m *CommandRequest) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *CommandRequest) GetEntry() *CommandEntry {
	if m != nil {
		return m.Entry
	}
	return nil
}

type CommandResponse struct {
	Cid                  string   `protobuf:"bytes,1,opt,name=cid,proto3" json:"cid,omitempty"`
	Result               string   `protobuf:"bytes,11,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandResponse) Reset()         { *m = CommandResponse{} }
func (m *CommandResponse) String() string { return proto.CompactTextString(m) }
func (*CommandResponse) ProtoMessage()    {}
func (*CommandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{6}
}

func (m *CommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponse.Unmarshal(m, b)
}
func (m *CommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponse.Marshal(b, m, deterministic)
}
func (m *CommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponse.Merge(m, src)
}
func (m *CommandResponse) XXX_Size() int {
	return xxx_messageInfo_CommandResponse.Size(m)
}
func (m *CommandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommandResponse proto.InternalMessageInfo

func (m *CommandResponse) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *CommandResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*RequestVoteRequest)(nil), "pb.RequestVoteRequest")
	proto.RegisterType((*RequestVoteResponse)(nil), "pb.RequestVoteResponse")
	proto.RegisterType((*AppendEntriesRequest)(nil), "pb.AppendEntriesRequest")
	proto.RegisterType((*AppendEntriesResponse)(nil), "pb.AppendEntriesResponse")
	proto.RegisterType((*CommandEntry)(nil), "pb.CommandEntry")
	proto.RegisterType((*CommandRequest)(nil), "pb.CommandRequest")
	proto.RegisterType((*CommandResponse)(nil), "pb.CommandResponse")
}

func init() { proto.RegisterFile("raft.proto", fileDescriptor_b042552c306ae59b) }

var fileDescriptor_b042552c306ae59b = []byte{
	// 444 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x95, 0x9b, 0xd0, 0x34, 0xe3, 0x24, 0x8d, 0xb6, 0x50, 0x96, 0x72, 0x09, 0x0b, 0x42, 0x11,
	0x87, 0x1c, 0xc2, 0x11, 0x84, 0x84, 0x42, 0x85, 0x8a, 0x7a, 0x5a, 0x21, 0xae, 0xd1, 0xc6, 0x9e,
	0x46, 0x96, 0x6c, 0xaf, 0xd9, 0xdd, 0x54, 0xf4, 0x87, 0xf0, 0x7f, 0xf8, 0x2d, 0xfc, 0x12, 0xb4,
	0x1f, 0x36, 0x0e, 0x0e, 0xb9, 0xcd, 0x3c, 0xbf, 0x99, 0x7d, 0xef, 0x69, 0x0c, 0xa0, 0xc4, 0x9d,
	0x59, 0x54, 0x4a, 0x1a, 0x49, 0x4e, 0xaa, 0x0d, 0xfb, 0x19, 0x01, 0xe1, 0xf8, 0x7d, 0x87, 0xda,
	0x7c, 0x93, 0x06, 0x43, 0x49, 0x08, 0xf4, 0x0d, 0xaa, 0x82, 0xc6, 0xb3, 0x68, 0xde, 0xe7, 0xae,
	0x26, 0x2f, 0x60, 0x94, 0x88, 0x32, 0xcd, 0x52, 0x61, 0x70, 0x9d, 0xa5, 0x74, 0x34, 0x8b, 0xe6,
	0x43, 0x1e, 0x37, 0xd8, 0x4d, 0x4a, 0x5e, 0xc1, 0x24, 0x17, 0xda, 0xac, 0x73, 0xb9, 0x5d, 0x67,
	0x65, 0x8a, 0x3f, 0xe8, 0xd8, 0x2d, 0x18, 0x59, 0xf4, 0x56, 0x6e, 0x6f, 0x2c, 0x46, 0x18, 0x8c,
	0x1b, 0x96, 0x7b, 0x65, 0xe2, 0x48, 0x71, 0x20, 0x7d, 0x45, 0x55, 0xb0, 0x5b, 0xb8, 0xd8, 0x93,
	0xa5, 0x2b, 0x59, 0x6a, 0xfc, 0x9f, 0xae, 0x7b, 0x69, 0x70, 0xbd, 0x55, 0xa2, 0x34, 0xe8, 0x75,
	0x9d, 0xf1, 0xd8, 0x62, 0x9f, 0x3d, 0xc4, 0x7e, 0x47, 0xf0, 0xf8, 0x63, 0x55, 0x61, 0x99, 0x5e,
	0x97, 0x46, 0x65, 0xa8, 0x8f, 0xf9, 0x7c, 0x0e, 0xc3, 0x1c, 0x45, 0x8a, 0xea, 0xaf, 0xc9, 0x33,
	0x0f, 0x78, 0x87, 0x95, 0xc2, 0xfb, 0xae, 0x43, 0x8b, 0xb6, 0x1d, 0x36, 0xac, 0xb6, 0xc3, 0x40,
	0xb2, 0x0e, 0xc9, 0x1b, 0x18, 0xa0, 0x17, 0x43, 0xcf, 0x67, 0xbd, 0x79, 0xbc, 0x9c, 0x2e, 0xaa,
	0xcd, 0x62, 0x25, 0x8b, 0x42, 0x78, 0x99, 0x0f, 0xbc, 0x26, 0x90, 0x97, 0x30, 0x0e, 0x92, 0x12,
	0x59, 0x14, 0x99, 0xa1, 0xd3, 0x10, 0xab, 0x03, 0x57, 0x0e, 0x63, 0xd7, 0xf0, 0xe4, 0x1f, 0x8f,
	0x47, 0x42, 0xa3, 0x30, 0xd0, 0xbb, 0x24, 0x41, 0xad, 0x43, 0x5e, 0x75, 0xcb, 0xde, 0xc3, 0xa8,
	0x2d, 0xc2, 0x32, 0x13, 0xdf, 0xd3, 0xc8, 0x85, 0x51, 0xb7, 0x76, 0xaf, 0x50, 0x5b, 0x4d, 0x4f,
	0x66, 0xbd, 0xf9, 0x90, 0xbb, 0x9a, 0x7d, 0x81, 0x49, 0x98, 0xae, 0x23, 0x9e, 0x42, 0x2f, 0xc9,
	0xea, 0x59, 0x5b, 0x92, 0xd7, 0xf0, 0xc8, 0x1a, 0x7b, 0x70, 0x82, 0x0e, 0xf9, 0xf6, 0x9f, 0xd9,
	0x3b, 0x38, 0x6f, 0x76, 0x05, 0x2b, 0xdd, 0x65, 0x97, 0x70, 0xaa, 0x50, 0xef, 0x72, 0xe3, 0xb6,
	0x0d, 0x79, 0xe8, 0x96, 0xbf, 0x22, 0xe8, 0x73, 0x71, 0x67, 0xc8, 0x07, 0x88, 0x5b, 0x97, 0x44,
	0x2e, 0xed, 0x6b, 0xdd, 0x8b, 0xbf, 0x7a, 0xda, 0xc1, 0xc3, 0x93, 0x9f, 0x60, 0xbc, 0x17, 0x2b,
	0xa1, 0x96, 0x79, 0xe8, 0x9a, 0xae, 0x9e, 0x1d, 0xf8, 0x12, 0xb6, 0x2c, 0x61, 0xb0, 0xaa, 0x63,
	0x6b, 0xf9, 0xad, 0x27, 0x2f, 0xf6, 0x30, 0x3f, 0xb3, 0x39, 0x75, 0xbf, 0xe9, 0xdb, 0x3f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xbe, 0x38, 0xd5, 0xf5, 0xb4, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RaftClient is the client API for Raft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RaftClient interface {
	RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error)
	AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error)
	// external user command
	Command(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
}

type raftClient struct {
	cc *grpc.ClientConn
}

func NewRaftClient(cc *grpc.ClientConn) RaftClient {
	return &raftClient{cc}
}

func (c *raftClient) RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error) {
	out := new(RequestVoteResponse)
	err := c.cc.Invoke(ctx, "/pb.Raft/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error) {
	out := new(AppendEntriesResponse)
	err := c.cc.Invoke(ctx, "/pb.Raft/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) Command(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, "/pb.Raft/Command", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftServer is the server API for Raft service.
type RaftServer interface {
	RequestVote(context.Context, *RequestVoteRequest) (*RequestVoteResponse, error)
	AppendEntries(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error)
	// external user command
	Command(context.Context, *CommandRequest) (*CommandResponse, error)
}

// UnimplementedRaftServer can be embedded to have forward compatible implementations.
type UnimplementedRaftServer struct {
}

func (*UnimplementedRaftServer) RequestVote(ctx context.Context, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (*UnimplementedRaftServer) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (*UnimplementedRaftServer) Command(ctx context.Context, req *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Command not implemented")
}

func RegisterRaftServer(s *grpc.Server, srv RaftServer) {
	s.RegisterService(&_Raft_serviceDesc, srv)
}

func _Raft_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Raft/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).RequestVote(ctx, req.(*RequestVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Raft_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Raft/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).AppendEntries(ctx, req.(*AppendEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Raft_Command_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).Command(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Raft/Command",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).Command(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Raft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Raft",
	HandlerType: (*RaftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _Raft_RequestVote_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _Raft_AppendEntries_Handler,
		},
		{
			MethodName: "Command",
			Handler:    _Raft_Command_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft.proto",
}
