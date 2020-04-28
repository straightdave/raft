package raft

// LogEntry ...
type LogEntry interface {
	Term() uint64
}

// EmptyLogEntry is used to fill as the first log
type EmptyLogEntry struct {
	term uint64
}

// NewEmptyLogEntryWithTerm ...
func NewEmptyLogEntryWithTerm(t uint64) *EmptyLogEntry {
	return &EmptyLogEntry{term: t}
}

// Term ...
func (e *EmptyLogEntry) Term() uint64 {
	return e.term
}
