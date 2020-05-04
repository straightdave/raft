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
		// Timeout for a new election.
		// Renewed at the beginning of each loop.
		electionTimeout := randomTimeoutInMSRange(150, 300)

		// Only two kinds of go-events are handled as a follower:
		// 1. Timeout for a new election;
		// 2. Incoming requests (both internal and external).
		select {
		case <-electionTimeout:
			// I'm tired of waiting. I am becoming a new leader.
			r.currentTerm++ // WARN: concurrent access
			go r.asCandidate()
			return

		case e := <-r.events:
			switch req := e.req.(type) {
			case *pb.AppendEntriesRequest:

				// Sorry your term seems out-of-date.
				// I don't listen to an old king. Thank you. Next.
				if req.Term < r.currentTerm {
					e.respCh <- &pb.AppendEntriesResponse{
						Term:    r.currentTerm,
						Success: false,
					}
					break
				}

				// Oh I see you, my boss, have already committed many.
				// I'll commit those just like you.
				if req.LeaderCommit > r.lastApplied {
					for i := r.lastApplied + 1; i <= req.LeaderCommit; i++ {
						_, err := r.exe.Apply(r.logs[i])
						if err != nil {
							break
						}
						r.lastApplied = i
					}
				}

				if len(req.Entries) == 0 {
					// This is a heartbeat.
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
				// You are old. You cannot be a leader. Even our terms are same now.
				// ('cause your term was one less than mine before becoming a candidate)
				if req.Term <= r.currentTerm {
					e.respCh <- &pb.RequestVoteResponse{
						Term:        r.currentTerm,
						VoteGranted: false,
					}
					break
				}

				// If:
				// 1. I've never had a leader;
				// 2. Actually you are my current leader;
				// I'll consider letting you be my leader.
				if r.votedFor == "" || r.votedFor == req.CandidateId {
					// In addition your LastLogIndex must be larger than mine.
					if req.LastLogIndex > uint64(len(r.logs)) {
						r.votedFor = req.CandidateId
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
				// Just tell clients who is the current leader.
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: r.exe.MakeRedirectionResponse(r.votedFor),
				}
			}
		}
	}
}
