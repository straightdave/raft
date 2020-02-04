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
			switch reqData := e.req.(type) {
			case *pb.AppendEntriesRequest:

				if reqData.Term < s.currentTerm {
					e.respCh <- &pb.AppendEntriesResponse{
						Term:    s.currentTerm,
						Success: false,
					}
					break
				}

				// apply last logs
				if reqData.LeaderCommit > s.lastApplied {
					s.exe.Apply(s.logs[s.lastApplied])
					s.lastApplied++
				}

				// append new logs
				for _, cmd := range reqData.Entries {
					logs := s.exe.cmd2logs(reqData.Term, cmd)
					s.logs = append(s.logs, logs...)
				}

				e.respCh <- &pb.AppendEntriesResponse{
					Term:    s.currentTerm,
					Success: true,
				}

			case *pb.RequestVoteRequest:
				if reqData.Term < s.currentTerm {
					e.respCh <- &pb.RequestVoteResponse{
						Term:        s.currentTerm,
						VoteGranted: false,
					}
					break
				}

				if s.votedFor == "" || s.votedFor == reqData.CandidateId {
					if reqData.LastLogIndex >= uint64(len(s.logs)-1) {
						s.votedFor = reqData.CandidateId
						s.leader = reqData.CandidateId
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
					Cid:    reqData.Cid,
					Result: fmt.Sprintf("redirect %s", s.leader),
				}
			}
		}
	}
}
