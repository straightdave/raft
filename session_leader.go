package main

import (
	"context"
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
	// role is only updated when trans-session.
	// role is read-only in other cases which is safe,
	// since during session transforming, no serving, no reading.
	s.role = leader
	s.resetIndices()
	log.Printf("%s becomes LEADER {term=%d}", s.selfID, s.currentTerm)

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	var (
		// signal with term for the leader to fallback as follower;
		fallBackCh = make(chan uint64, len(s.peers))
	)

	for {
		select {
		case term := <-fallBackCh:
			s.currentTerm = term
			go s.asFollower()
			return

		case <-heartbeatTicker.C:
			go s.callPeers(ctx, nil, fallBackCh)

		case e := <-s.events:
			switch req := e.req.(type) {
			case *pb.CommandRequest:
				logs := s.exe.cmd2logs(s.currentTerm, req.Entry)
				s.logs = append(s.logs, logs...)

				//
				//
				//
				//
				//
				//
				//

			case *pb.RequestVoteRequest:
				if req.Term > s.currentTerm {
					fallBackCh <- req.Term
				}

			case *pb.AppendEntriesRequest:
				if req.Term > s.currentTerm {
					fallBackCh <- req.Term
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

	cc, err = grpc.Dial(addr)
	if err != nil {
		return
	}
	defer cc.Close()

	c := pb.NewRaftClient(cc)
	return c.AppendEntries(ctx, &pb.AppendEntriesRequest{
		Term:         s.currentTerm,
		LeaderId:     s.selfID,
		PrevLogIndex: lastIndex,
		PrevLogTerm:  s.logs[lastIndex].Term,
		Entries:      entries,
		LeaderCommit: s.commitIndex,
	})
}
