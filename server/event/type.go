package event

// Type ...
type Type uint

// event types
const (
	Unknown Type = iota
	AppendEntriesCall
	RequestVoteCall
)
