package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

func (s *Server) asCandidate() {
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var (
		voteCount       = 1
		quotum          = len(s.peers) / 2
		validVotesCh    = make(chan string, len(s.peers))
		stateTransCh    = make(chan role, len(s.peers))
		electionTimeout = randomTimeout150300()
	)

	s.role = candidate
	s.currentTerm++
	s.votedFor = s.selfID // vote for self
	log.Printf("%s becomes CANDIDATE {term=%d}", s.selfID, s.currentTerm)

	s.requestVotes(ctx, validVotesCh, stateTransCh)

	// through the life time as candidate:
	// this for-select loop doesn't deal with requests from outside.
	for {
		select {
		case <-validVotesCh:
			voteCount++
			if voteCount > quotum {
				go s.asLeader()
				return
			}

		case role := <-stateTransCh:
			switch role {
			case follower:
				go s.asFollower()
				return
			}

		case <-electionTimeout:
			go s.asCandidate()
			return
		}
	}
}

func (s *Server) requestVotes(ctx context.Context, validVotesCh chan<- string, stateTransCh chan<- role) {
	for _, addr := range s.peers {
		go s.requestVote(ctx, addr, validVotesCh, stateTransCh)
	}
}

func (s *Server) requestVote(ctx context.Context, addr string, validVotesCh chan<- string, stateTransCh chan<- role) {
	defer rescue(func(err error) {
		log.Printf("rescue: [%s] RPC RequestVote: %v", addr, err)
	})

	cc, err := grpc.Dial(addr)
	if err != nil {
		log.Printf("error: [%s] dialing failed: %v", addr, err)
		return
	}
	defer cc.Close()

	c := pb.NewRaftClient(cc)
	resp, err := c.RequestVote(ctx, &pb.RequestVoteRequest{
		Term:        s.currentTerm,
		CandidateId: s.selfID,
	})
	if err != nil {
		log.Printf("error: [%s] RPC RequestVote failed: %v", addr, err)
		return
	}

	if resp.Term > s.currentTerm {
		stateTransCh <- follower
		return
	}

	if resp.VoteGranted {
		validVotesCh <- addr
	}
}
