package raft

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/straightdave/raft/pb"
)

const (
	respTimeout = 500 * time.Millisecond
)

var (
	errContextDone = fmt.Errorf("context done")
	errTimeout     = fmt.Errorf("time out")
)

// implementation of pb.RaftServer
type impl struct {
	raft *Raft
}

// NewRaftServer ...
func NewRaftServer(port uint, opts ...Option) pb.RaftServer {
	return &impl{raft: NewRaft(port, opts...)}
}

// Command ...
func (s *impl) Command(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	ch := make(chan proto.Message, 1)
	defer close(ch)

	s.raft.events <- &protoEvent{
		req:    req,
		respCh: ch,
	}

	select {
	case resp := <-ch:
		return resp.(*pb.CommandResponse), nil
	case <-time.After(respTimeout):
		return nil, errTimeout
	case <-ctx.Done():
		return nil, errContextDone
	}
}

// RequestVote ...
func (s *impl) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	ch := make(chan proto.Message, 1)
	defer close(ch)

	s.raft.events <- &protoEvent{
		req:    req,
		respCh: ch,
	}

	select {
	case resp := <-ch:
		return resp.(*pb.RequestVoteResponse), nil
	case <-time.After(respTimeout):
		return nil, errTimeout
	case <-ctx.Done():
		return nil, errContextDone
	}
}

// AppendEntries ...
func (s *impl) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	ch := make(chan proto.Message, 1)
	defer close(ch)

	s.raft.events <- &protoEvent{
		req:    req,
		respCh: ch,
	}

	select {
	case resp := <-ch:
		return resp.(*pb.AppendEntriesResponse), nil
	case <-time.After(respTimeout):
		return nil, errTimeout
	case <-ctx.Done():
		return nil, errContextDone
	}
}
