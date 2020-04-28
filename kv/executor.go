package kv

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/straightdave/raft"
	"github.com/straightdave/raft/pb"
)

type impl struct {
	store map[string]string
}

// NewExecutor ...
func NewExecutor() raft.Executor {
	return &impl{
		store: make(map[string]string),
	}
}

// MakeRedirectionResponse ...
func (e *impl) MakeRedirectionResponse(leader string) string {
	return fmt.Sprintf("REDIRECT %s", leader)
}

// ToRaftLog ...
func (e *impl) ToRaftLog(term uint64, request proto.Message) (raft.LogEntry, error) {
	req, ok := request.(*pb.CommandRequest)
	if ok {
		return &LogEntry{
			term: term,
			verb: strings.ToUpper(req.Entry.Command),
			args: req.Entry.Args,
		}, nil
	}
	return nil, fmt.Errorf("unknown message")
}

// Apply ...
func (e *impl) Apply(entry raft.LogEntry) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	l, ok := entry.(*LogEntry)
	if ok {
		switch l.verb {
		case "SET":
			if len(l.args) < 2 {
				err = fmt.Errorf("SET command needs 2 args (%v)", l.args)
				return
			}
			e.store[l.args[0]] = l.args[1]
			result = "OK"

		case "GET":
			if len(l.args) < 1 {
				err = fmt.Errorf("GET command needs 1 args (%v)", l.args)
				return
			}
			result, _ = e.store[l.args[0]]

		default:
			err = fmt.Errorf("unknown command %s", l.verb)
		}
		return
	}
	return "", fmt.Errorf("unknown log type")
}
