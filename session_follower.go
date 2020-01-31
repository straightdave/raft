package main

import (
	"log"
)

func (s *Server) asFollower() {
	// mutex here: different sessions are mutually exclusive.
	// each session should have mutex with full-function scope
	// as the critical section.
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	s.propertyLock.Lock()
	s.role = FOLLOWER
	s.propertyLock.Unlock()
	log.Printf("Becomes FOLLOWER")

	for {
		// reset on each round of looping
		timeout := randomTimeout150300()

		select {
		case <-timeout:
			// exit FOLLOWER session and start CANDIDIATE session
			defer func() { go s.asCandidate() }()
			return

		}
	}
}
