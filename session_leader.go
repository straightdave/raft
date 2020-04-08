package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

const (
	heartbeatInterval = 100 * time.Millisecond
	retryInterval     = 100 * time.Millisecond
)

func (s *Server) asLeader() {
	// the session lock makes sure that new session must start after
	// the previous session ends.
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	// some acts (as sub-goroutines) of leader are infinite, so here
	// we need to cancel sub goroutines when leader session ends.
	bctx := context.Background()
	ctx, cancelfunc := context.WithCancel(bctx)
	defer cancelfunc()

	// it's safe to change role here;
	// role is only changed during session transition.
	// role is read-only in other cases which is safe,
	// since during session transforming, no serving, no reading.
	s.role = leader
	s.resetIndices()
	log.Printf("%s becomes LEADER {term=%d}", s.selfID, s.currentTerm)

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	// listen to signals (payload is the term) to fallback as follower
	fallBackCh := make(chan uint64, len(s.peers))

	// forever loop as being a leader
	for {
		select {

		// handle the fallback signals and quit
		case term := <-fallBackCh:
			s.currentTerm = term
			go s.asFollower()
			// since session is locked until asLeader session ends,
			// so here `go s.asFollower()` will start a goroutine but block
			// until all deferred actions are done in leader session.
			return

		// periodically send heartbeat
		case <-heartbeatTicker.C:
			go s.callPeers(ctx, nil, fallBackCh)

		// handle events (mainly internal raft or external requests)
		// each event payload contains:
		// * request data
		// * response channel
		case e := <-s.events:
			switch req := e.req.(type) {

			// External command requests
			case *pb.CommandRequest:
				logs, _ := s.exe.ToRaftLogs(s.currentTerm, req.Entry)
				s.logs = append(s.logs, logs...)
				go s.callPeers(ctx, []*pb.CommandEntry{req.Entry}, fallBackCh)
				e.respCh <- &pb.CommandResponse{
					Cid:    req.Cid,
					Result: fmt.Sprintf("%d", len(logs)),
				}

			// Raft requests: RequestVote (from other one who's running an election)
			// when receiving this, the leader will only concern about
			// the term recorded in the payload. if that term is newer than
			// current term, then fallback
			case *pb.RequestVoteRequest:
				if req.Term > s.currentTerm {
					// "OK you're boss. I'll fallback."
					e.respCh <- &pb.RequestVoteResponse{
						Term:        s.currentTerm,
						VoteGranted: true,
					}
					fallBackCh <- req.Term
				} else {
					// "I am Boss. You old beggar."
					e.respCh <- &pb.RequestVoteResponse{
						Term:        s.currentTerm,
						VoteGranted: false,
					}
				}

			// Raft requests: AppendEntries (from other one who thinks it's the leader)
			case *pb.AppendEntriesRequest:
				if req.Term > s.currentTerm {
					// "You are the boss. I'll fallback and please send data again."
					// CAUTION: will this response cause another AppendEntries call?
					e.respCh <- &pb.AppendEntriesResponse{
						Term:    s.currentTerm,
						Success: false,
					}
					fallBackCh <- req.Term
				} else {
					// "I am Boss. You old beggar."
					e.respCh <- &pb.RequestVoteResponse{
						Term:        s.currentTerm,
						VoteGranted: false,
					}
				}
			}
		}
	}
}

func (s *Server) resetIndices() {
	lastLogIndex := uint64(len(s.logs) - 1)
	for _, peer := range s.peers {
		s.nextIndex[peer] = lastLogIndex
		s.matchIndex[peer] = 0
	}
}

func (s *Server) callPeers(ctx context.Context, entries []*pb.CommandEntry, needFallback chan<- uint64) {
	var (
		responses = make(chan *pb.AppendEntriesResponse, len(s.peers))
		quorum    = len(s.peers) / 2
		count     = 0
	)

	// sending the requests to each peer forever until it's success,
	// or the leader session ends from outside (ctx is done).
	for _, addr := range s.peers {
		go func(ctx context.Context, addr string) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				resp, err := s.appendEntries(addr, entries)
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

func (s *Server) appendEntries(addr string, entries []*pb.CommandEntry) (resp *pb.AppendEntriesResponse, err error) {
	defer rescue(func(e error) {
		err = e
	})

	var (
		lastIndex = uint64(len(s.logs) - 1)
		cc        *grpc.ClientConn
		ctx       = context.Background()
	)

	// using the cached grpc connection, or
	// lazy-initialize it for the address

	s.connGuard.RLock()
	cc, ok := s.grpcConnections[addr]
	s.connGuard.RUnlock()

	if !ok {
		cc, err = grpc.Dial(addr)
		if err != nil {
			return
		}

		s.connGuard.Lock()
		s.grpcConnections[addr] = cc
		s.connGuard.Unlock()
	}

	c := pb.NewRaftClient(cc)
	return c.AppendEntries(ctx, &pb.AppendEntriesRequest{
		Term:         s.currentTerm,
		LeaderId:     s.selfID,
		PrevLogIndex: lastIndex,
		PrevLogTerm:  s.logs[lastIndex].Term(),
		Entries:      entries,
		LeaderCommit: s.commitIndex,
	})
}
