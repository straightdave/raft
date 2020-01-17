package main

// CommandEntry is the local version of pb.CommandEntry
type CommandEntry struct {
	Command string
	Args    []string
}

// Log ...
type Log struct {
	Term     uint64
	Commands []CommandEntry
}
