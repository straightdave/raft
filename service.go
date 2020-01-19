package main

import (
	"context"

	"github.com/straightdave/raft/pb"
)

// ServerServiceImpl handles communication, but the main
// logic is handled inside the Server instance.
type ServerServiceImpl struct {
	server *Server
}

// NewServerServiceImpl ... 'servers' is the initial list of other nodes.
func NewServerServiceImpl(servers []string) *ServerServiceImpl {
	return &ServerServiceImpl{server: NewServer(servers)}
}

// RequestVote ...
func (s *ServerServiceImpl) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	s.server.requestVoteCalls <- req
	select {
	case resp := <-s.server.requestVoteResps:
		return resp, nil
	}
}

// AppendEntries ...
func (s *ServerServiceImpl) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	s.server.appendEntriesCalls <- req
	select {
	case resp := <-s.server.appendEntriesResps:
		return resp, nil
	}
}

// Command ...
func (s *ServerServiceImpl) Command(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	s.server.commandCalls <- req
	select {
	case resp := <-s.server.commandResps:
		return resp, nil
	}
}
