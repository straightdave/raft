package main

import (
	"log"
	"sync"
	"time"

	"github.com/straightdave/raft/server/event"
	pb "github.com/straightdave/raft/server/pb"
)

// NodeRole is role for the raft nodes.
type NodeRole uint

// Roles ...
const (
	FOLLOWER NodeRole = iota
	CANDIDATE
	LEADER
)

type heartbeat struct{}

// Node is the raft server node.
type Node struct {
	mgr *event.Manager

	// role init as FOLLOWER
	role      NodeRole
	roleGuard sync.RWMutex

	currentTerm      uint64
	currentTermGuard sync.RWMutex

	others *nodeList
	logs   []pb.Entry

	// channel for perceiving AppendEntries RPC calls from the leader;
	// the AppendEntries call also acts like the leader's heartbeat;
	appendingEntriesCalls chan *pb.AppendEntriesRequest

	// channel for perceiving RequestVote RPC calls from other candidates;
	requestVoteCalls chan *pb.RequestVoteRequest
}

// NewNode ...
func NewNode(nodes []string) *Node {
	n := &Node{
		others: newNodeList(nodes),

		appendingEntriesCalls: make(chan *pb.AppendEntriesRequest, 10),
		requestVoteCalls:      make(chan *pb.RequestVoteRequest, 10),
	}

	// log index starts from 1
	n.logs = append(n.logs, pb.Entry{Key: "__0", Value: "__0"})

	// init as FOLLOWER, so wait for leader's heartbeat
	n.startHeartbeatMonitor()

	return n
}

// Role ...
func (n *Node) Role() NodeRole {
	n.roleGuard.RLock()
	defer n.roleGuard.RUnlock()
	return n.role
}

// SetRole ...
func (n *Node) SetRole(newRole NodeRole) {
	if n.Role() == newRole {
		return
	}

	n.roleGuard.Lock()
	defer n.roleGuard.Unlock()

	n.role = newRole
	go n.onRoleChanged(newRole)
}

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

func (n *Node) onRoleChanged(newRole NodeRole) {
	switch newRole {
	case FOLLOWER: // things to do as soon as becoming a FOLLOWER
		n.startHeartbeatMonitor()
	case CANDIDATE: // things to do as soon as becoming a CANDIDATE
		n.startElection()
	case LEADER: // things to do as soon as becoming a LEADER
		n.startHeartbeating()
	}
	log.Printf("role changed to %d", newRole)
}

func (n *Node) onAppendEntriesRequestReceived(req *pb.AppendEntriesRequest) {

}

func (n *Node) onRequestVoteRequestReceived(req *pb.RequestVoteRequest) {

}

func (n *Node) startHeartbeatMonitor() {
	go func(ctx context.Context) {
		for {
			timeout := time.After(time.Duration(getRandom150300()) * time.Millisecond)
			select {
			case req := <-n.appendingEntriesCalls:
				n.onAppendEntriesRequestReceived(req)
			case <-timeout:
				n.SetRole(CANDIDATE)
				return // end monitoring goroutine
			case <-ctx.Done():
				return // end monitoring goroutine
			}
		}
	}(n.ctx)
}

func (n *Node) startElection() {
	n.incrTerm()
	go func(ctx context.Context) {
		for {
			timeout := time.After(time.Duration(getRandom150300()) * time.Millisecond)
			select {
			case req := <-n.appendingEntriesCalls:
				n.onAppendEntriesRequestReceived(req)
			case <-timeout:
				n.SetRole(CANDIDATE)
				return // end monitoring goroutine
			case <-ctx.Done():
				return // end monitoring goroutine
			}
		}
	}(n.ctx)
}

func (n *Node) startHeartbeating() {
	go func() {

	}()
}
