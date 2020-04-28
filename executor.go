package raft

import "github.com/golang/protobuf/proto"

// Executor ...
type Executor interface {
	// MakeRedirectionResponse turns one request (especially CommandRequest) into
	// another form of proto message containing redirection info as responses.
	MakeRedirectionResponse(leader string) string
	// ToRaftLog converts one gRPC payload to LogEntry
	ToRaftLog(term uint64, request proto.Message) (LogEntry, error)
	// Apply applies LogEntry to the system
	Apply(LogEntry) (string, error)
}

// DummyExecutor ...
type DummyExecutor struct{}

// MakeRedirectionResponse ...
func (DummyExecutor) MakeRedirectionResponse(leader string) string {
	return "redirect"
}

// ToRaftLog ...
func (DummyExecutor) ToRaftLog(term uint64, request proto.Message) (LogEntry, error) {
	return NewEmptyLogEntryWithTerm(term), nil
}

// Apply ...
func (DummyExecutor) Apply(log LogEntry) (string, error) {
	return "", nil
}
