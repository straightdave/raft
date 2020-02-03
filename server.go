package main

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
)

// Role ...
type Role uint

// server roles ...
const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)

type protoEvent struct {
	req    proto.Message
	respCh chan proto.Message
}

// Server ...
type Server struct {
	sessionLock sync.Mutex

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

	events chan *protoEvent
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

		events: make(chan *protoEvent, 50),
	}

	s.logs = append(s.logs, LogEntry{}) // log index starts from 1
	go s.asFollower()
	return s
}
