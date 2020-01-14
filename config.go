package main

import (
	"sync"
)

type nodeList struct {
	sync.RWMutex
	m map[string]int
}

func newNodeList(nodes []string) *nodeList {
	res := &nodeList{m: make(map[string]int)}
	for _, node := range nodes {
		res.m[node] = 1
	}
	return res
}

func (l *nodeList) add(nodes ...string) {
	l.Lock()
	defer l.Unlock()
	for _, node := range nodes {
		l.m[node] = 1
	}
}

func (l *nodeList) remove(nodes ...string) {
	l.Lock()
	defer l.Unlock()
	for _, node := range nodes {
		delete(l.m, node)
	}
}

func (l *nodeList) snapshot() []string {
	l.RLock()
	defer l.RUnlock()
	var res []string
	for k := range l.m {
		res = append(res, k)
	}
	return res
}
