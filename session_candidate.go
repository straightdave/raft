package raft

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

func (r *Raft) asCandidate() {
	r.sessionLock.Lock()
	defer r.sessionLock.Unlock()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var (
		voteCount       = 1
		quotum          = len(r.peers) / 2
		validVotesCh    = make(chan string, len(r.peers))
		stateTransCh    = make(chan Role, len(r.peers))
		electionTimeout = randomTimeoutInMSRange(150, 300)
	)

	r.role = Candidate
	r.currentTerm++
	r.votedFor = r.selfID // vote for self
	log.Printf("%s becomes CANDIDATE {term=%d}", r.selfID, r.currentTerm)

	r.requestVotes(ctx, validVotesCh, stateTransCh)

	for {
		select {
		case e := <-r.events:
			switch req := e.req.(type) {
			case *pb.RequestVoteRequest:
				e.respCh <- &pb.RequestVoteResponse{
					Term:        r.currentTerm,
					VoteGranted: false,
				}

			case *pb.AppendEntriesRequest:
				e.respCh <- &pb.AppendEntriesResponse{
					Term:    r.currentTerm,
					Success: false,
				}

				if req.Term > r.currentTerm {
					r.currentTerm = req.Term
					go r.asFollower()
					return
				}

			case *pb.CommandRequest:
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: r.exe.MakeRedirectionResponse(r.selfID),
				}
			}

		case <-validVotesCh:
			voteCount++
			if voteCount > quotum {
				go r.asLeader()
				return
			}

		case role := <-stateTransCh:
			switch role {
			case Follower:
				go r.asFollower()
				return
			}

		case <-electionTimeout:
			go r.asLeader()
			return
		}
	}
}

func (r *Raft) requestVotes(ctx context.Context, validVotesCh chan<- string, stateTransCh chan<- Role) {
	for _, addr := range r.peers {
		go r.requestVote(ctx, addr, validVotesCh, stateTransCh)
	}
}

func (r *Raft) requestVote(ctx context.Context, addr string, validVotesCh chan<- string, stateTransCh chan<- Role) {
	defer rescue(func(err error) {
		log.Printf("rescue: [%s] RPC RequestVote: %v", addr, err)
	})

	cc, err := grpc.Dial(addr)
	if err != nil {
		log.Printf("error: [%s] dialing failed: %v", addr, err)
		return
	}
	defer cc.Close()

	lastIndex := uint64(len(r.logs) - 1)
	c := pb.NewRaftClient(cc)
	resp, err := c.RequestVote(ctx, &pb.RequestVoteRequest{
		Term:         r.currentTerm,
		CandidateId:  r.selfID,
		LastLogIndex: lastIndex,
		LastLogTerm:  r.logs[lastIndex].Term(),
	})
	if err != nil {
		log.Printf("error: [%s] RPC RequestVote failed: %v", addr, err)
		return
	}

	if resp.Term > r.currentTerm {
		stateTransCh <- Follower
		return
	}

	if resp.VoteGranted {
		validVotesCh <- addr
	}
}
