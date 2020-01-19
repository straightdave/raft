package main

import (
	"sync/atomic"

	"github.com/straightdave/raft/pb"
)

const (
	// CHSIZE default event channel size
	CHSIZE = 10
)

// Role ...
type Role uint

// server roles ...
const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)

// Log ...
type Log struct {
	Term     uint64
	Commands []*pb.CommandEntry
}

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

		appendEntriesCalls: make(chan *pb.AppendEntriesRequest, CHSIZE),
		appendEntriesResps: make(chan *pb.AppendEntriesResponse, CHSIZE),
		requestVoteCalls:   make(chan *pb.RequestVoteRequest, CHSIZE),
		requestVoteResps:   make(chan *pb.RequestVoteResponse, CHSIZE),
		commandCalls:       make(chan *pb.CommandRequest, CHSIZE),
		commandResps:       make(chan *pb.CommandResponse, CHSIZE),
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
