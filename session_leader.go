package main

import (
	"log"
)

func (n *Node) asLeader() {
	log.Printf("> as leader")

	n.setRole(LEADER)
}
