package main

import (
	"sync"
)

// ServerList ...
type ServerList struct {
	sync.RWMutex
	m map[string]int
}

// NewServerList ...
func NewServerList(servers []string) *ServerList {
	res := &ServerList{m: make(map[string]int)}
	for _, server := range servers {
		res.m[server] = 1
	}
	return res
}

// Add ...
func (s *ServerList) Add(servers ...string) {
	s.Lock()
	defer s.Unlock()
	for _, server := range servers {
		s.m[server] = 1
	}
}

// Remove ...
func (s *ServerList) Remove(servers ...string) {
	s.Lock()
	defer s.Unlock()
	for _, server := range servers {
		delete(s.m, server)
	}
}

// Snapshot ...
func (s *ServerList) Snapshot() []string {
	s.RLock()
	defer s.RUnlock()
	var res []string
	for k := range s.m {
		res = append(res, k)
	}
	return res
}
