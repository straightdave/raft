package raft

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
)

type protoEvent struct {
	req    proto.Message
	respCh chan proto.Message
}

// Raft ...
type Raft struct {
	sessionLock sync.Mutex
	role        Role
	exe         Executor

	selfID string // ip + port
	peers  []string

	votedFor string
	leader   string
	logs     []LogEntry

	currentTerm uint64
	commitIndex uint64
	lastApplied uint64
	nextIndex   map[string]uint64
	matchIndex  map[string]uint64

	events chan *protoEvent
}

// NewRaft ...
func NewRaft(port uint, opts ...Option) *Raft {
	ip, err := getLocalIP()
	if err != nil {
		panic(err)
	}

	r := &Raft{
		exe:        DummyExecutor{},
		selfID:     fmt.Sprintf("%s:%d", ip, port),
		nextIndex:  make(map[string]uint64),
		matchIndex: make(map[string]uint64),
		events:     make(chan *protoEvent, 50),
	}

	for _, o := range opts {
		o(r)
	}

	r.logs = append(r.logs, &EmptyLogEntry{}) // log index starts from 1
	go r.asFollower()
	return r
}
