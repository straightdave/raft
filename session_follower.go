package main

import (
	"fmt"
	"log"

	"github.com/straightdave/raft/pb"
)

func (s *Server) asFollower() {
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	s.role = follower
	log.Printf("%s becomes FOLLOWER", s.selfID)

	for {
		timeout := randomTimeout150300()

		select {
		case <-timeout:
			go s.asCandidate()
			return

		case e := <-s.events:
			switch req := e.req.(type) {
			case *pb.AppendEntriesRequest:
				if req.Term < s.currentTerm {
					e.respCh <- &pb.AppendEntriesResponse{
						Term:    s.currentTerm,
						Success: false,
					}
					break
				}

				// apply last logs
				if req.LeaderCommit > s.lastApplied {
					_, err := s.exe.Apply(s.logs[s.lastApplied])
					if err == nil {
						s.lastApplied++
					}
				}

				if len(req.Entries) == 0 {
					// heartbeat
					break
				}

				// append new logs
				for _, cmd := range req.Entries {
					logs := s.exe.cmd2logs(req.Term, cmd)
					s.logs = append(s.logs, logs...)
				}

				e.respCh <- &pb.AppendEntriesResponse{
					Term:    s.currentTerm,
					Success: true,
				}

			case *pb.RequestVoteRequest:
				if req.Term < s.currentTerm {
					e.respCh <- &pb.RequestVoteResponse{
						Term:        s.currentTerm,
						VoteGranted: false,
					}
					break
				}

				if s.votedFor == "" || s.votedFor == req.CandidateId {
					if req.LastLogIndex >= uint64(len(s.logs)-1) {
						s.votedFor = req.CandidateId
						s.leader = req.CandidateId
						e.respCh <- &pb.RequestVoteResponse{
							Term:        s.currentTerm,
							VoteGranted: true,
						}
						break
					}
				}

				e.respCh <- &pb.RequestVoteResponse{
					Term:        s.currentTerm,
					VoteGranted: false,
				}

			case *pb.CommandRequest:
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: fmt.Sprintf("redirect %s", s.leader),
				}
			}
		}
	}
}
