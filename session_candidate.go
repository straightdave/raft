package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

func (s *Server) asCandidate() {
	log.Printf("> as candidate")
	s.role = CANDIDATE
	ctx := context.Background()
	cctx, cancelFunc := context.WithCancel(ctx)

	// when this function returns (ends), all sub-goroutines started with
	// the context:cctx would receive cancel signals to stop working.
	defer cancelFunc()

	var (
		voteNum = 0
		others  = s.others.Snapshot()
	)

	// increase current term before sending vote requests
	s.incrTerm()
	validVotes, stateTrans := s.requestVotes(cctx, others)

	for {
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
			switch role {
			case FOLLOWER:
				go s.asFollower()
			case CANDIDATE:
				go s.asCandidate()
			case LEADER:
				go s.asLeader()
			}
			// NOTE: no matter what role it would be, it would returs.
			// (actually only FOLLOWER)
			return

		case <-electionTimeout:
			go s.asCandidate() // start another round of election
			return

		/* other raft endpoint calls */

		case req := <-s.appendEntriesCalls:
			// whenever receiving appendEntries calls, returns to FOLLOWER role
			// (handle this request as a follower, don't waste it)
			go s.handleAppendEntriesCallAsFollower(req)
			go s.asFollower()
			return

		case req := <-s.requestVoteCalls:
			// whenever receiving other candidates' vote requests:
			// => trans to FOLLOWER if req.Term is bigger;
			// => others, ignore & carry on (no return)

			if req.Term > s.currentTerm {
				go s.asFollower()
				return
			}

			// candidate don't handle external command calls
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
		go s.requestVoteFrom(ctx, validVotes, stateTrans, addr)
	}
	return validVotes, stateTrans
}

// start a goroutine to call one node's RequestVote and track.
// if success, push addr into the valid votes channel.
func (s *Server) requestVoteFrom(ctx context.Context, validVotes chan<- string, stateTrans chan<- Role, addr string) {
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
