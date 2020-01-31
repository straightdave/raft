package main

import (
	"context"
	"fmt"
	"time"

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

	return nil, nil
}

// AppendEntries ...
func (s *ServerServiceImpl) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {

	return nil, nil
}

// Command ...
func (s *ServerServiceImpl) Command(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	ch := make(chan *pb.CommandResponse, 1)
	s.server.chCommandRespLock.Lock()
	s.server.chCommandRespMap[req.Cid] = ch
	s.server.chCommandRespLock.Unlock()

	s.server.chCommandReq <- req

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(respTimeout):
		return nil, errTimeout
	case <-ctx.Done():
		return nil, errContextDone
	}
}
