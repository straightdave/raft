package main

// Role ...
type Role uint

// server roles ...
const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)
