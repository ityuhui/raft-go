package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"raft-go/raft_rpc"
	"time"

	"google.golang.org/grpc"
)

type Node struct {
	role            NodePole
	electionTimeout int
	name            string
	port            string
}

func NewNode(host string, port string) *Node {
	return &Node{
		role:            NodeRole_Follower,
		electionTimeout: 0,
		name:            host,
		port:            port,
	}
}

func (n *Node) resetElectionTimeout() {
	n.electionTimeout = 0
}

func (n *Node) Run() {
	go n.startRaftServer()
	n.mainLoop()
}

func (n *Node) mainLoop() {
	for {
		fmt.Printf("node is running on %s:%s ...\n", n.name, n.port)
		time.Sleep(time.Second)
	}
}

type server struct {
	raft_rpc.UnimplementedRaftServiceServer
}

func (s *server) TellMyHeartBeatToFollower(ctx context.Context, in *raft_rpc.HeartBeatRequest) (*raft_rpc.HeartBeatReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &raft_rpc.HeartBeatReply{Message: "Hello " + in.GetName()}, nil
}

func (n *Node) startRaftServer() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	raft_rpc.RegisterRaftServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
