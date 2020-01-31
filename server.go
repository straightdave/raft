package main

import (
	"fmt"
	"sync"

	"github.com/straightdave/raft/pb"
)

// Server ...
type Server struct {
	sessionLock  sync.Mutex
	propertyLock sync.RWMutex

	exe *Executor

	selfID string // ip + port
	peers  []string

	role     Role
	votedFor string
	leader   string
	logs     []LogEntry

	currentTerm uint64
	commitIndex uint64
	lastApplied uint64

	nextIndex  map[string]uint64
	matchIndex map[string]uint64

	// Command requests from clients.
	chCommandReq chan *pb.CommandRequest

	// collects Command responses for clients.
	chCommandRespMap  map[string]chan *pb.CommandResponse
	chCommandRespLock sync.Mutex
}

// NewServer ...
func NewServer(port uint, peers []string) *Server {
	ip, err := getLocalIP()
	if err != nil {
		panic(err)
	}

	s := &Server{
		exe: &Executor{},

		selfID:     fmt.Sprintf("%s:%d", ip, port),
		peers:      peers,
		nextIndex:  make(map[string]uint64),
		matchIndex: make(map[string]uint64),

		chCommandReq:     make(chan *pb.CommandRequest, 50),
		chCommandRespMap: make(map[string]chan *pb.CommandResponse),
	}

	s.logs = append(s.logs, LogEntry{}) // log index starts from 1
	go s.asFollower()
	return s
}
