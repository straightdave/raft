package main

import (
	"context"

	pb "github.com/straightdave/raft/pb"
)

// ServerServiceImpl handles communication, but the main
// logic is handled inside the Node instance.
type ServerServiceImpl struct {
	node *Node
}

// NewServerServiceImpl ... 'nodes' is the initial list of other nodes.
func NewServerServiceImpl(nodes []string) *ServerServiceImpl {
	return &ServerServiceImpl{node: NewNode(nodes)}
}

// RequestVote ...
func (s *ServerServiceImpl) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return s.node.onRequestVoteRequestReceived(ctx, req)
}

// AppendEntries ...
func (s *ServerServiceImpl) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	return s.node.onAppendEntriesRequestReceived(ctx, req)
}

// Command ...
func (s *ServerServiceImpl) Command(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	return s.node.Command(ctx, req)
}
