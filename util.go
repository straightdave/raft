package raft

import (
	"math/rand"
	"net"
	"time"
)

func randomTimeout150300() <-chan time.Time {
	return time.After(time.Duration(rand.Intn(150)+150) * time.Millisecond)
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

type actWithErr func(error)

func rescue(acts ...actWithErr) {
	if r := recover(); r != nil {
		e := r.(error)
		for _, act := range acts {
			act(e)
		}
	}
}
