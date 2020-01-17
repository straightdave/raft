package main

import (
	"fmt"
	"log"

	"github.com/straightdave/raft/pb"
)

func (s *Server) asFollower() {
	log.Printf("> as follower")
	s.role = FOLLOWER

	timeout := randomTimeout150300()
	for {
		select {
		// if timeout, quit current session
		case <-timeout:
			go s.asCandidate()
			return

		// handle calls, not quitting session
		case req := <-s.appendEntriesCalls:
			go s.handleAppendEntriesCallAsFollower(req)
			// only in this case, it would reset the timer
			timeout = randomTimeout150300()
		case req := <-s.requestVoteCalls:
			go s.handleRequestVoteCallAsFollower(req)
		case req := <-s.commandCalls:
			go s.handleCommandsAsFollower(req)
		}
	}
}

func (s *Server) handleAppendEntriesCallAsFollower(req *pb.AppendEntriesRequest) {

}

func (s *Server) handleRequestVoteCallAsFollower(req *pb.RequestVoteRequest) {
	term := s.currentTerm

	// the requesting candidate is out-of-date
	if term > req.Term {
		s.requestVoteResps <- &pb.RequestVoteResponse{
			Term:        term,
			VoteGranted: false,
		}
		return
	}

	// first time voting or last vote is the same candidate
	if s.votedFor == "" || s.votedFor == req.CandidateId {
		if req.LastLogIndex >= s.commitIndex {
			s.requestVoteResps <- &pb.RequestVoteResponse{
				Term:        term,
				VoteGranted: true,
			}
		} else {
			// if req's log is out-of-date, don't grant vote
			s.requestVoteResps <- &pb.RequestVoteResponse{
				Term:        term,
				VoteGranted: false,
			}
		}
		return
	}

	// other cases: vote granted
	s.requestVoteResps <- &pb.RequestVoteResponse{
		Term:        term,
		VoteGranted: true,
	}
}

func (s *Server) handleCommandsAsFollower(req *pb.CommandRequest) {
	// NOTE: here's defined a redirect command
	s.commandResps <- &pb.CommandResponse{
		Result: fmt.Sprintf("redirect %s", s.leader),
	}
}
