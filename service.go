package main

import (
	"context"
	"fmt"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/straightdave/raft/pb"
)

// ServerServiceImpl ...
type ServerServiceImpl struct {
	server *Server
}

var (
	errContextDone = fmt.Errorf("context done")
	errTimeout     = fmt.Errorf("time out")
)

const (
	respTimeout = 500 * time.Millisecond
)

// NewServerServiceImpl ...
func NewServerServiceImpl(port uint, servers []string) *ServerServiceImpl {
	return &ServerServiceImpl{server: NewServer(port, servers)}
}

// RequestVote ...
func (s *ServerServiceImpl) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	ch := make(chan proto.Message, 1)
	defer close(ch)

	s.server.events <- &protoEvent{
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
func (s *ServerServiceImpl) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	ch := make(chan proto.Message, 1)
	defer close(ch)

	s.server.events <- &protoEvent{
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

// Command ...
func (s *ServerServiceImpl) Command(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	ch := make(chan proto.Message, 1)
	defer close(ch)

	s.server.events <- &protoEvent{
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
