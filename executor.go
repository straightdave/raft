package main

import (
	"github.com/straightdave/raft/pb"
)

// LogEntry ...
type LogEntry interface {
	Term() uint64
}

type emptyLogEntry struct{}

// Term ...
func (e emptyLogEntry) Term() uint64 {
	return 0
}

// Executor ...
type Executor interface {
	ToRaftLogs(term uint64, payload *pb.CommandEntry) ([]LogEntry, error)
	Apply(logs LogEntry) error
}

// NewExecutor ...
func NewExecutor() Executor {
	// TODO
	return nil
}
