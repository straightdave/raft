package main

import (
	"errors"

	"github.com/straightdave/raft/pb"
)

var (
	// ErrNotSupported means such command has not been implemented yet.
	ErrNotSupported = errors.New("not supported")
)

// LogEntry ...
type LogEntry struct {
	Term    uint64
	Command string
	Args    []string
}

// Executor handles real command execution, working as the
// so called 'replicated state machine' in each node describbed in Raft.
type Executor struct {
}

func (e *Executor) cmd2logs(term uint64, entry *pb.CommandEntry) []LogEntry {
	// currently cmd:log = 1:1
	// sooner it will support cmd:log = 1:
	l := LogEntry{
		Term:    term,
		Command: entry.Command,
	}
	copy(l.Args, entry.Args)
	return []LogEntry{l}
}

// Apply logs
func (e *Executor) Apply(entries ...LogEntry) (string, error) {

	return "", nil
}
