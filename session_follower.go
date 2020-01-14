package main

import (
	"log"

	"github.com/straightdave/raft/pb"
)

func (n *Node) asFollower() {
	log.Printf("> as follower")
	n.setRole(FOLLOWER)

	for {
		// TODO: make heartbeat detect timeout configurable
		timeout := randomTimeout150300()

		select {
		case <-timeout:
			go n.asCandidate()
			return

		/* handle other raft calls */
		case req := <-n.appendEntriesCalls:
			go n.handleAppendEntriesCallAsFollower(req)
		case req := <-n.requestVoteCalls:
			go n.handleRequestVoteCallAsFollower(req)

			/* handle external calls */
			// TODO: implement. mainly just redirection
		}
	}
}

func (n *Node) handleAppendEntriesCallAsFollower(req *pb.AppendEntriesRequest) {

}

func (n *Node) handleRequestVoteCallAsFollower(req *pb.RequestVoteRequest) {
	term := n.getTerm()
	if term > req.Term {
		n.requestVoteResps <- &pb.RequestVoteResponse{
			Term:        term,
			VoteGranted: false,
		}
		return
	}

	if n.votedFor == "" || req.CandidateId {

	}

}
