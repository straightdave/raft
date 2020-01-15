package main

import (
	"context"
	"sync"

	pb "github.com/straightdave/raft/pb"
)

// Node ...
type Node struct {
	ip     string // local outbound IP (used as ID)
	others *nodeList

	exe *Executor

	role      Role
	roleGuard sync.RWMutex

	currentTerm      uint64
	currentTermGuard sync.RWMutex

	// votedFor: who I've voted in the current term
	votedFor string

	commitIndex uint64

	// raft calls (reqs & resps)
	// since we are using async pattern here
	appendEntriesCalls chan *pb.AppendEntriesRequest
	appendEntriesResps chan *pb.AppendEntriesResponse
	requestVoteCalls   chan *pb.RequestVoteRequest
	requestVoteResps   chan *pb.RequestVoteResponse
}

// NewNode ...
func NewNode(others []string) *Node {
	ip, err := getLocalIP()
	if err != nil {
		panic(err)
	}

	n := &Node{
		ip:     ip,
		others: newNodeList(others),
		exe:    &Executor{},

		appendEntriesCalls: make(chan *pb.AppendEntriesRequest, 10),
		appendEntriesResps: make(chan *pb.AppendEntriesResponse, 10),
		requestVoteCalls:   make(chan *pb.RequestVoteRequest, 10),
		requestVoteResps:   make(chan *pb.RequestVoteResponse, 10),
	}

	go n.asFollower()
	return n
}

func (n *Node) setRole(r Role) {
	n.roleGuard.Lock()
	defer n.roleGuard.Unlock()
	n.role = r
}

// ====== public ======

// Command ...
func (n *Node) Command(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	res, err := n.exe.Apply(ctx, req.Entry)
	if err != nil {
		return &pb.CommandResponse{
			Result: err.Error(),
		}, nil
	}
	return &pb.CommandResponse{
		Result: res,
	}, nil
}

// ====== Support gRPC communication ======

func (n *Node) onAppendEntriesRequestReceived(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	n.currentTermGuard.RLock()
	defer n.currentTermGuard.RUnlock()

	if req.Term < n.currentTerm {
		return &pb.AppendEntriesResponse{
			Term:    n.currentTerm,
			Success: false,
		}, nil
	}

	return nil, nil
}

func (n *Node) onRequestVoteRequestReceived(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return nil, nil
}

// ====== private works ======

func (n *Node) incrTerm() {
	n.currentTermGuard.Lock()
	defer n.currentTermGuard.Unlock()
	n.currentTerm++
}

func (n *Node) getTerm() uint64 {
	n.currentTermGuard.RLock()
	defer n.currentTermGuard.RUnlock()
	return n.currentTerm
}

// Role gets the role
func (n *Node) Role() Role {
	n.roleGuard.RLock()
	defer n.roleGuard.RUnlock()
	return n.role
}
