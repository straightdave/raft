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
	log.Printf("Becomes FOLLOWER")

	for {
		timeout := randomTimeout150300()

		select {
		case <-timeout:
			go s.asCandidate()
			return

		case e := <-s.events:
			switch reqData := e.req.(type) {
			case *pb.AppendEntriesRequest:
				// TODO: implement
				// the most important logic

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
