package raft

import (
	"math/rand"
	"net"
	"time"
)

func randomTimeoutInMSRange(min, max int) <-chan time.Time {
	delta := max - min
	return time.After(time.Duration(rand.Intn(min)+delta) * time.Millisecond)
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
