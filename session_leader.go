package main

import (
	"context"
	"log"
	"time"

	// "google.golang.org/grpc"

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
	s.role = leader
	s.resetIndices()
	log.Printf("%s becomes LEADER {term=%d}", s.selfID, s.currentTerm)

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:

		case e := <-s.events:
			switch reqData := e.req.(type) {
			case *pb.CommandRequest:
				// TODO: implement
				// important!

			case *pb.RequestVoteRequest:
				if reqData.Term > s.currentTerm {
					s.currentTerm = reqData.Term
					go s.asFollower()
					return
				}

			case *pb.AppendEntriesRequest:
				if reqData.Term > s.currentTerm {
					s.currentTerm = reqData.Term
					go s.asFollower()
					return
				}
			}
		}
	}
}

func (s *Server) resetIndices() {
	lastLogIndex := uint64(len(s.logs) - 1)
	for _, peer := range s.peers {
		s.nextIndex[peer] = lastLogIndex
		s.matchIndex[peer] = 0
	}
}
