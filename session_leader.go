package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/straightdave/raft/pb"
)

func (s *Server) asLeader() {
	s.role = LEADER
	log.Printf("Becomes LEADER")

	for {
		select {
		case req := <-s.appendEntriesCalls:
			// NOTE: '>=' or just '>' ?
			if req.Term >= s.currentTerm {
				go s.asFollower()
				return
			}
		case req := <-s.requestVoteCalls:
			// NOTE: '>=' or just '>' ?
			if req.Term >= s.currentTerm {
				go s.asFollower()
				return
			}
		case req := <-s.commandCalls:
			go s.handleCommandsAsLeader(req)
		}
	}
}

func (s *Server) handleCommandsAsLeader(req *pb.CommandRequest) {
	// 1 - append to []log
	s.log = append(s.log, Log{
		Term: s.currentTerm,
		Commands: []*pb.CommandEntry{
			req.Entry,
		},
	})

	// 2- send AppendEntries calls to all followers
	others := s.others.Snapshot()
	okNum := 0
	validResults := s.appendEntries(req, others)

OUT:
	for {
		select {
		case <-validResults:
			okNum++
			if okNum > len(others)/2 { // majority
				break OUT
			}
		}
	}

	// 3 - commit log (execution & commitIndex++)
	res, err := s.exe.Apply(req.GetEntry())
	if err != nil {
		s.commandResps <- &pb.CommandResponse{
			Result: fmt.Sprintf("err %s", err.Error()),
		}
		return
	}

	s.incrCommitIndex()
	s.commandResps <- &pb.CommandResponse{
		Result: res,
	}
}

func (s *Server) appendEntries(req *pb.CommandRequest, others []string) <-chan string {
	validResults := make(chan string, len(others))
	for _, addr := range s.others.Snapshot() {
		go s.appendEntry(validResults, req, addr)
	}
	return validResults
}

func (s *Server) appendEntry(validResults chan<- string, req *pb.CommandRequest, addr string) {
	cc, err := grpc.Dial(addr)
	if err != nil {
		return
	}
	defer cc.Close()

	prevLogIndex := uint64(len(s.log) - 1)

	c := pb.NewRaftClient(cc)
	resp, err := c.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
		Term:         s.currentTerm,
		LeaderId:     s.ip,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  s.log[prevLogIndex].Term,
		Entries:      []*pb.CommandEntry{req.Entry},
		LeaderCommit: s.commitIndex,
	})
	if err != nil {
		return
	}

	if resp.Success {
		validResults <- addr
	}
}
