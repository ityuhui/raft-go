package main

import (
	"fmt"
	"time"
)

type Node struct {
	role            NodePole
	electionTimeout int
}

func NewNode() *Node {
	return &Node{
		role:            NodeRole_Follower,
		electionTimeout: 0,
	}
}

func (n *Node) resetElectionTimeout() {
	n.electionTimeout = 0
}

func (n *Node) MainLoop() {
	for {
		fmt.Println("node is running...")
		time.Sleep(time.Second)
	}
}
