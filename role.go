package main

// Role ...
type Role uint

const (
	// INIT is default value of Role (zero value)
	// which is meaningless in Raft
	INIT Role = iota
	// FOLLOWER ...
	FOLLOWER
	// CANDIDATE ...
	CANDIDATE
	// LEADER ...
	LEADER
)
