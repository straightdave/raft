package raft

// Role ...
type Role uint

const (
	// Follower is the Raft role of Follower.
	Follower Role = iota
	// Candidate is the Raft role of Candidate.
	Candidate
	// Leader is the Raft role of Leader.
	Leader
)
