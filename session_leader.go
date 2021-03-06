package raft

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

const (
	heartbeatInterval = 100 * time.Millisecond
	retryInterval     = 100 * time.Millisecond
)

func (r *Raft) asLeader(reasons ...string) {
	// the session lock makes sure that new session must start after
	// the previous session endr.
	r.sessionLock.Lock()
	defer r.sessionLock.Unlock()

	// some acts (as sub-goroutines) of leader are infinite, so here
	// we need to cancel sub goroutines when leader session endr.
	bctx := context.Background()
	ctx, cancelfunc := context.WithCancel(bctx)
	defer cancelfunc()

	// it's safe to change role here;
	// role is only changed during session transition.
	// role is read-only in other cases which is safe,
	// since during session transforming, no serving, no reading.
	r.role = Leader
	r.resetIndices()
	log.Printf("%s becomes LEADER {term=%d} due to reasons=%v", r.selfID, r.currentTerm, reasons)

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	// listen to signals (payload is the term) to fallback as follower
	fallBackCh := make(chan uint64, len(r.peers))

	// forever loop as being a leader
	for {
		select {

		// periodically send heartbeat
		case <-heartbeatTicker.C:
			go r.callPeers(ctx, nil, fallBackCh)

		// handle the fallback signals and quit
		case term := <-fallBackCh:
			r.currentTerm = term
			go r.asFollower("fallback required")
			// since session is locked until asLeader session ends,
			// so here `go r.asFollower()` will start a goroutine but block
			// until all deferred actions are done in leader session.
			return

		// handle events (mainly internal raft or external requests)
		// each event payload contains:
		// * request data
		// * response channel
		case e := <-r.events:
			switch req := e.req.(type) {

			// External command requests
			case *pb.CommandRequest:
				log, _ := r.exe.ToRaftLog(r.currentTerm, req)
				r.logs = append(r.logs, log)
				go r.callPeers(ctx, []*pb.CommandEntry{req.Entry}, fallBackCh)
				result, err := r.exe.Apply(log)
				if err != nil {
					result = err.Error()
				}
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: result,
				}

			// Raft requests: RequestVote (from other one who's running an election)
			// when receiving this, the leader will only concern about
			// the term recorded in the payload. if that term is newer than
			// current term, then fallback
			case *pb.RequestVoteRequest:
				if req.Term > r.currentTerm {
					// "OK you're bosr. I'll fallback."
					e.respCh <- &pb.RequestVoteResponse{
						Term:        r.currentTerm,
						VoteGranted: true,
					}
					fallBackCh <- req.Term
				} else {
					// "I am Boss. You old beggar."
					e.respCh <- &pb.RequestVoteResponse{
						Term:        r.currentTerm,
						VoteGranted: false,
					}
				}

			// Raft requests: AppendEntries (from other one who thinks it's the leader)
			case *pb.AppendEntriesRequest:
				if req.Term > r.currentTerm {
					// "You are the boss. I'll fallback and please send data again."
					// CAUTION: will this response cause another AppendEntries call?
					e.respCh <- &pb.AppendEntriesResponse{
						Term:    r.currentTerm,
						Success: false,
					}
					fallBackCh <- req.Term
				} else {
					// "I am Boss. You old beggar."
					e.respCh <- &pb.RequestVoteResponse{
						Term:        r.currentTerm,
						VoteGranted: false,
					}
				}
			}
		}
	}
}

func (r *Raft) resetIndices() {
	lastLogIndex := uint64(len(r.logs) - 1)
	for _, peer := range r.peers {
		r.nextIndex[peer] = lastLogIndex
		r.matchIndex[peer] = 0
	}
}

func (r *Raft) callPeers(ctx context.Context, entries []*pb.CommandEntry, needFallback chan<- uint64) {
	var (
		responses = make(chan *pb.AppendEntriesResponse, len(r.peers))
		quorum    = len(r.peers) / 2
		count     = 0
	)

	// sending the requests to each peer forever until success,
	// or the leader session ends from outside (ctx is done).
	for _, addr := range r.peers {
		go func(ctx context.Context, addr string) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				resp, err := r.appendEntries(ctx, addr, entries)
				if err == nil {
					responses <- resp
					return
				}

				log.Printf("[append-entries] err: %v", err)
				time.Sleep(retryInterval)
			}
		}(ctx, addr)
	}

	// blocks here.
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-responses:
			if resp.Success {
				count++
				if count > quorum {

					// CAUTION! to be fixed, data racing

					return
				}
			}
		}
	}
}

func (r *Raft) appendEntries(ctx context.Context, addr string, entries []*pb.CommandEntry) (resp *pb.AppendEntriesResponse, err error) {
	defer rescue(func(e error) {
		err = e
		log.Printf("rescue: [%s] RPC RequestVote: %v\nStack: %s", addr, err, debug.Stack())
	})

	var (
		lastIndex = uint64(len(r.logs) - 1)
		cc        *grpc.ClientConn
	)

	cc, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("error: [%s] dialing failed: %v", addr, err)
		return
	}
	defer cc.Close()
	c := pb.NewRaftClient(cc)

	return c.AppendEntries(ctx, &pb.AppendEntriesRequest{
		Term:         r.currentTerm,
		LeaderId:     r.selfID,
		PrevLogIndex: lastIndex,
		PrevLogTerm:  r.logs[lastIndex].Term(),
		Entries:      entries,
		LeaderCommit: r.commitIndex,
	})
}
