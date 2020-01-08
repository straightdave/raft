package main

import (
	"context"

	pb "github.com/straightdave/raft/server/pb"
)

// ServerServiceImpl is the Raft server service implementing RaftServer interface (in pb/raft.pb.go)
type ServerServiceImpl struct {
	node *Node
}

// NewServerServiceImpl ...
func NewServerServiceImpl(ctx context.Context, nodes []string) *ServerServiceImpl {
	return &ServerServiceImpl{
		node: NewNode(ctx, nodes),
	}
}

// RequestVote ...
func (s *ServerServiceImpl) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return nil, nil
}

// AppendEntries ...
func (s *ServerServiceImpl) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	return nil, nil
}
