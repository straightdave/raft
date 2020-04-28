package raft

import (
	"log"

	"github.com/straightdave/raft/pb"
)

func (r *Raft) asFollower() {
	r.sessionLock.Lock()
	defer r.sessionLock.Unlock()

	r.role = Follower
	log.Printf("%s becomes FOLLOWER", r.selfID)

	for {
		timeout := randomTimeout150300()

		select {
		case <-timeout:
			go r.asCandidate()
			return

		case e := <-r.events:
			switch req := e.req.(type) {
			case *pb.AppendEntriesRequest:
				if req.Term < r.currentTerm {
					e.respCh <- &pb.AppendEntriesResponse{
						Term:    r.currentTerm,
						Success: false,
					}
					break
				}

				// apply last logs
				if req.LeaderCommit > r.lastApplied {
					_, err := r.exe.Apply(r.logs[r.lastApplied])
					if err == nil {
						r.lastApplied++
					}
				}

				if len(req.Entries) == 0 {
					// heartbeat
					break
				}

				// append new logs
				for _, cmd := range req.Entries {
					log, _ := r.exe.ToRaftLog(req.Term, cmd)
					r.logs = append(r.logs, log)
				}

				e.respCh <- &pb.AppendEntriesResponse{
					Term:    r.currentTerm,
					Success: true,
				}

			case *pb.RequestVoteRequest:
				if req.Term < r.currentTerm {
					e.respCh <- &pb.RequestVoteResponse{
						Term:        r.currentTerm,
						VoteGranted: false,
					}
					break
				}

				if r.votedFor == "" || r.votedFor == req.CandidateId {
					if req.LastLogIndex >= uint64(len(r.logs)-1) {
						r.votedFor = req.CandidateId
						r.leader = req.CandidateId
						e.respCh <- &pb.RequestVoteResponse{
							Term:        r.currentTerm,
							VoteGranted: true,
						}
						break
					}
				}

				e.respCh <- &pb.RequestVoteResponse{
					Term:        r.currentTerm,
					VoteGranted: false,
				}

			case *pb.CommandRequest:
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: r.exe.MakeRedirectionResponse(r.selfID),
				}
			}
		}
	}
}
