package main

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

type role uint

const (
	follower role = iota
	candidate
	leader
)

type protoEvent struct {
	req    proto.Message
	respCh chan proto.Message
}

// Server ...
type Server struct {
	sessionLock sync.Mutex

	// cache of grpc connections
	connGuard       sync.RWMutex
	grpcConnections map[string]*grpc.ClientConn

	exe Executor

	selfID string // ip + port
	peers  []string

	role     role
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
		exe:             NewExecutor(),
		grpcConnections: make(map[string]*grpc.ClientConn),

		selfID:     fmt.Sprintf("%s:%d", ip, port),
		peers:      peers,
		nextIndex:  make(map[string]uint64),
		matchIndex: make(map[string]uint64),

		events: make(chan *protoEvent, 50),
	}

	s.logs = append(s.logs, emptyLogEntry{}) // log index starts from 1
	go s.asFollower()
	return s
}
