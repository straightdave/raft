package kv

// LogEntry ..
type LogEntry struct {
	term uint64
	verb string
	args []string
}

// Term ...
func (l *LogEntry) Term() uint64 {
	return l.term
}
