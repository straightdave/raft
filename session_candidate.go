package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

func (s *Server) asCandidate() {
	ctx := context.Background()
	cctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	var (
		voteNum = 1 // start with 1 vote from oneself
		others  = s.others.Snapshot()
	)

	s.role = CANDIDATE
	s.incrTerm()
	s.votedFor = s.ip // vote for self

	log.Printf("Becomes CANDIDIATE {%d}", s.currentTerm)

	validVotes, stateTrans := s.requestVotes(cctx, others)
	for {
		// create/reset the election timeout for a new election.
		// TODO: make election timeout configurable
		electionTimeout := randomTimeout150300()

		select {
		case <-validVotes:
			voteNum++
			if voteNum > len(others)/2 { // majority
				go s.asLeader()
				return
			}

		case role := <-stateTrans:
			// unexpected results during requesting votes,
			// forcing candidate to change role.
			// normally just change to FOLLOWER.

			switch role {
			case FOLLOWER:
				go s.asFollower()
				return
			}

		case <-electionTimeout:
			go s.asCandidate() // start another round of election
			return

		case req := <-s.appendEntriesCalls:
			if req.Term < s.currentTerm {
				// "no matter who you are, you are out-of-date."
				// ignore it and continue the election
				break
			}

			// handle this request and return to FOLLOWER
			go s.handleAppendEntriesCallAsFollower(req)
			go s.asFollower()
			return

		case req := <-s.requestVoteCalls:
			// NOTE: the paper doesn't talk about this case.

			// whenever receiving other candidates' vote requests:
			// => trans to FOLLOWER if req.Term is bigger;
			// => others, ignore & carry on (no return)
			if req.Term > s.currentTerm {
				go s.asFollower()
				return
			}
		}
	}
}

// candidate requests votes and starts one goroutine for each node to track responses.
// this is only called by node in the state of candidate.
// one return value is the 'valid votes', value is the addr of the responder;
// another return value is the 'state trans' indicater.
func (s *Server) requestVotes(ctx context.Context, others []string) (<-chan string, <-chan Role) {
	// initialize two indicator channels whenever begins to request votes
	validVotes := make(chan string, len(others))
	stateTrans := make(chan Role, 1)

	for _, addr := range s.others.Snapshot() {
		go s.requestVote(ctx, validVotes, stateTrans, addr)
	}
	return validVotes, stateTrans
}

// start a goroutine to call one node's RequestVote and track.
// if success, push addr into the valid votes channel.
func (s *Server) requestVote(ctx context.Context, validVotes chan<- string, stateTrans chan<- Role, addr string) {
	cc, err := grpc.Dial(addr)
	if err != nil {
		return
	}
	defer cc.Close()

	c := pb.NewRaftClient(cc)
	resp, err := c.RequestVote(ctx, &pb.RequestVoteRequest{
		Term:        s.currentTerm, // sender's current term
		CandidateId: s.ip,          // sender's ID (using IP here)
	})
	if err != nil {
		return
	}

	if resp.Term > s.currentTerm {
		// whenever it gets a response with a bigger term, it should
		// return to the FOLLOWER role.
		stateTrans <- FOLLOWER
		return
	}

	if resp.VoteGranted {
		validVotes <- addr
	}
}
