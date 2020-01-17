package main

import (
	"sync/atomic"

	pb "github.com/straightdave/raft/pb"
)

// Server ...
type Server struct {
	ip     string
	others *ServerList
	exe    *Executor

	role     Role
	votedFor string
	leader   string
	log      []Log

	currentTerm uint64
	commitIndex uint64
	lastApplied uint64

	appendEntriesCalls chan *pb.AppendEntriesRequest
	appendEntriesResps chan *pb.AppendEntriesResponse
	requestVoteCalls   chan *pb.RequestVoteRequest
	requestVoteResps   chan *pb.RequestVoteResponse
	commandCalls       chan *pb.CommandRequest
	commandResps       chan *pb.CommandResponse
}

// NewServer ...
func NewServer(others []string) *Server {
	ip, err := getLocalIP()
	if err != nil {
		panic(err)
	}

	s := &Server{
		ip:     ip,
		others: NewServerList(others),
		exe:    &Executor{},

		appendEntriesCalls: make(chan *pb.AppendEntriesRequest, 10),
		appendEntriesResps: make(chan *pb.AppendEntriesResponse, 10),
		requestVoteCalls:   make(chan *pb.RequestVoteRequest, 10),
		requestVoteResps:   make(chan *pb.RequestVoteResponse, 10),
		commandCalls:       make(chan *pb.CommandRequest, 10),
		commandResps:       make(chan *pb.CommandResponse, 10),
	}

	s.log = append(s.log, Log{}) // log index starts from 1
	go s.asFollower()
	return s
}

func (s *Server) incrTerm() {
	atomic.AddUint64(&s.currentTerm, 1)
}

func (s *Server) incrCommitIndex() {
	atomic.AddUint64(&s.commitIndex, 1)
}

func (s *Server) incrLastApplied() {
	atomic.AddUint64(&s.lastApplied, 1)
}
