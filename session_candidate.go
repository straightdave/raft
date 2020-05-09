package raft

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

type fallbackCmd struct {
	LeaderAddr string
	Term       uint64
}

func (r *Raft) asCandidate(reasons ...string) {
	r.sessionLock.Lock()
	defer r.sessionLock.Unlock()

	log.Printf("%s becomes CANDIDATE {term=%d} due to reasons=%v", r.selfID, r.currentTerm, reasons)
	r.role = Candidate
	r.votedFor = r.selfID // vote for self

	var (
		voteCount       = 1
		quotum          = len(r.peers) / 2
		validVotesCh    = make(chan string, len(r.peers))
		fallbackCh      = make(chan fallbackCmd, len(r.peers))
		electionTimeout = randomTimeoutInMSRange(150, 300)
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	r.requestVotes(ctx, validVotesCh, fallbackCh)

	for {
		select {

		case <-electionTimeout:
			// Especially in single-server (no peer) mode, requests for vote won't get
			// any response, so the case `<-validateVotesCh` won't happen. In this case,
			// we let the single server be the leader directly.
			if len(r.peers) == 0 {
				go r.asLeader("single node", "election timeout")
				return
			}

			// Otherwise, it should start a new candidate session.
			go r.asCandidate("election timeout")
			return

		case <-validVotesCh:
			// Collecting granted vote count from responses of vote requests.
			voteCount++
			if voteCount > quotum {
				// WIN!
				go r.asLeader("win")
				return
			}

		case fb := <-fallbackCh:
			// Responses of vote requests require a fallback.
			r.currentTerm = fb.Term
			r.votedFor = fb.LeaderAddr
			go r.asFollower("fallback", "leader:"+fb.LeaderAddr)
			return

		case e := <-r.events:
			switch req := e.req.(type) {
			case *pb.RequestVoteRequest:
				// Case: receiving vote requests from other candidates.
				// If the incoming requests are sent by candidates with higher or same terms,
				// it gives up candidate session and fallbacks to be FOLLOWER.
				if req.Term >= r.currentTerm {
					r.currentTerm = req.Term
					r.votedFor = req.CandidateId
					e.respCh <- &pb.RequestVoteResponse{
						Term:        r.currentTerm,
						VoteGranted: true,
					}
					go r.asFollower("other candiadate wins")
					return
				}

				e.respCh <- &pb.RequestVoteResponse{
					Term:        r.currentTerm,
					VoteGranted: false,
				}

			case *pb.AppendEntriesRequest:
				// Case: receiving append log requests from other leaders.
				// If the incoming requests are sent by leaders with higher terms,
				// fallback as their followers.
				// NOTE:
				// In this case, we just fallback as the follower without appending logs but just
				// responding with false. So then the leader would send another append-log request.

				// No matter fallback or not, respond with false.
				e.respCh <- &pb.AppendEntriesResponse{
					Term:    r.currentTerm,
					Success: false,
				}

				if req.Term > r.currentTerm {
					r.currentTerm = req.Term
					r.votedFor = req.LeaderId
					go r.asFollower("valid leader requests")
					return
				}

			case *pb.CommandRequest:
				// Case: receiving external command requests.
				// This case is simpler, just respond with current leader address (voteFor).
				// NOTE: actually leader address is set as self address in the beginning of session.
				// Client requests may get the address that is not the final leader. So client just
				// retry and finally will get correct one.
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: r.exe.MakeRedirectionResponse(r.votedFor),
				}
			}
		}
	}
}

func (r *Raft) requestVotes(ctx context.Context, validVotesCh chan<- string, fallbackCh chan<- fallbackCmd) {
	for _, addr := range r.peers {
		go r.requestVote(ctx, addr, validVotesCh, fallbackCh)
	}
}

func (r *Raft) requestVote(ctx context.Context, addr string, validVotesCh chan<- string, fallbackCh chan<- fallbackCmd) {
	defer rescue(func(err error) {
		log.Printf("rescue: [%s] RPC RequestVote: %v\nStack: %s", addr, err, debug.Stack())
	})

	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("error: [%s] dialing failed: %v", addr, err)
		return
	}
	defer cc.Close()
	c := pb.NewRaftClient(cc)

	lastIndex := uint64(len(r.logs) - 1)
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

	// The response shows the target is higher than the candidate, so fallback.
	if resp.Term > r.currentTerm {
		fallbackCh <- fallbackCmd{LeaderAddr: addr, Term: resp.Term}
		return
	}

	if resp.VoteGranted {
		validVotesCh <- addr
	}
}
