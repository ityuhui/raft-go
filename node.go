package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"raft-go/raft_rpc"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
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

func (n *Node) Run() {
	go startRaftServer()
	mainLoop()
}

func (n *Node) mainLoop() {
	for {
		fmt.Println("node is running...")
		time.Sleep(time.Second)
	}
}

// server is used to implement raft_rpc.RaftServiceServer.
type server struct {
	raft_rpc.UnimplementedGreeterServer
}

func (s *server) TellMyHeartBeatToFollower(ctx context.Context, in *raft_rpc.HeartBeatRequest) (*raft_rpc.HeartBeatReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &raft_rpc.HeartBeatReply{Message: "Hello " + in.GetName()}, nil
}

func (n *Node) startRaftServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	raft_rpc.RegisterRaftServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
