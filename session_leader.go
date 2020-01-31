package main

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

const (
	heartbeatInterval = 100 * time.Millisecond
	retryInterval     = 100 * time.Millisecond
)

func (s *Server) asLeader() {
	// the session lock makes sure that new session must start after
	// the previous session ends.
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	// some acts (as sub-goroutines) of leader are infinite, so here
	// we need to cancel sub goroutines when leader session ends.
	bctx := context.Background()
	ctx, cancelfunc := context.WithCancel(bctx)
	defer cancelfunc()

	// it's safe to change role here;
	// role is only updated when trans-session.
	// role is read-only in other cases which is safe,
	// since during session transforming, no serving, no reading.
	s.role = LEADER
	s.resetIndices()
	log.Printf("%s becomes LEADER {term=%d}", s.ip, s.currentTerm)

	// send the init heartbeat, then wait for the result.
	if !s.heartbeat(ctx) {
		go s.asFollower()
		return
	}

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			s.heartbeat(ctx)

		case resp := <-s.chAppendEntriesResp:
			// leader gets responses from other servers after sending
			// heartbeats or AppendEntries requests.

			if resp.Term > s.currentTerm {
				// reading s.currentTerm is safe. Term is only changed
				// during candidate session transforming. and in that
				// case, no readings of s.currentTerm.

			}

		}
	}
}

func (s *Server) appendEntries(ctx context.Context, req *pb.CommandRequest, others []string) <-chan string {
	validResults := make(chan string, len(others))
	for _, addr := range s.others {
		go s.appendEntry(ctx, validResults, req, addr)
	}
	return validResults
}

func (s *Server) appendEntry(ctx context.Context, validResults chan<- string, req *pb.CommandRequest, addr string) {
	defer rescue()

	prevLogIndex := uint64(len(s.logs) - 1)
	aeReq := &pb.AppendEntriesRequest{
		Term:         s.currentTerm,
		LeaderId:     s.ip,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  s.logs[prevLogIndex].Term,
		Entries:      []*pb.CommandEntry{req.Entry},
		LeaderCommit: s.commitIndex,
	}

	ok := s.call(ctx, addr, eAppendEntries, aeReq)
	if ok {
		return
	}

	retryTicker := time.NewTicker(retryInterval)
	defer retryTicker.Stop()

	for {
		select {
		case <-retryTicker.C:
			ok := s.call(ctx, addr, eAppendEntries, aeReq)
			if ok {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// Send heartbeat to followers then wait for the result.
// It returns `true` if majority is OK;
// It returns `false` if it must fall back to follower.
func (s *Server) heartbeat(ctx context.Context) bool {
	var (
		okCount      = 0
		prevLogIndex = uint64(len(s.logs) - 1)
	)

	for _, addr := range s.others {
		go func(ctx context.Context, addr string) {

			s.call(ctx, addr, eAppendEntries, &pb.AppendEntriesRequest{
				Term:         s.currentTerm,
				LeaderId:     s.ip,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  s.logs[prevLogIndex].Term,
				LeaderCommit: s.commitIndex,
			})
		}(ctx, addr)
	}

	return true
}

type endpoint string

const (
	eAppendEntries endpoint = "AppendEntries"
	eRequestVotes           = "RequestVotes"
)

// call endpoint once and put the possible responses to the channels
func (s *Server) call(ctx context.Context, addr string, end endpoint, req interface{}) bool {
	log.Printf("CALL %s:%s", addr, end)
	cc, err := grpc.Dial(addr)
	if err != nil {
		log.Printf("CALL dial err %v", err)
		return false
	}
	defer cc.Close()

	c := pb.NewRaftClient(cc)
	var resp interface{}
	switch end {
	case eAppendEntries:
		resp, err = c.AppendEntries(ctx, req.(*pb.AppendEntriesRequest))
		if err == nil {
			s.appendEntriesResps <- resp.(*pb.AppendEntriesResponse)
			return true
		}

	case eRequestVotes:
		resp, err = c.RequestVote(ctx, req.(*pb.RequestVoteRequest))
		if err == nil {
			s.requestVoteResps <- resp.(*pb.RequestVoteResponse)
			return true
		}
	}

	log.Printf("CALL %s:%s err %v", addr, end, err)
	return false
}

// Reset leader's nextIndex and matchIndex after election.
//
// NOTE:
//
// 1) not thread-safe.
//
// 2) assumes 'others' don't change.
func (s *Server) resetIndices() {
	lastLogIndex := len(s.logs) - 1
	for _, otherID := range s.others {
		s.nextIndex[otherID] = uint64(lastLogIndex)
		s.matchIndex[otherID] = 0
	}
}
