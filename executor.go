package main

import (
	"context"

	"github.com/straightdave/raft/pb"
)

// Executor handles real command execution, working as the
// so called 'replicated state machine' in each node describbed in Raft.
type Executor struct {
}

// Apply ...
func (e *Executor) Apply(ctx context.Context, entries ...*pb.CommandEntry) (string, error) {

	return "", nil
}
