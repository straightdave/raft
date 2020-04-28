package raft

// Option ...
type Option func(*Raft)

// WithPeers ...
func WithPeers(peers []string) Option {
	return func(r *Raft) {
		r.peers = peers
	}
}

// WithExecutor ...
func WithExecutor(executor Executor) Option {
	return func(r *Raft) {
		r.exe = executor
	}
}
