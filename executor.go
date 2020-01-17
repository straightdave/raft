package main

import (
	"errors"

	"github.com/straightdave/raft/pb"
)

var (
	// ErrNotSupported means such command has not been implemented yet.
	ErrNotSupported = errors.New("not supported")
)

// Executor handles real command execution, working as the
// so called 'replicated state machine' in each node describbed in Raft.
type Executor struct {
}

// Apply ...
func (e *Executor) Apply(entries ...*pb.CommandEntry) (string, error) {

	return "", nil
}
