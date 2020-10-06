package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"raft-go/raft_rpc"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	ELECTION_TIMEOUT     = 5
	INC_ELECTION_TIMEOUT = 1
)

type Node struct {
	role            NodePole
	currentTerm     int
	electionTimeout int
	myAddr          *Address
	peers           []*Address
}

var ins *Node = nil
var lock sync.Mutex

func NewNodeInstance(Iam string, peers string) *Node {

	ins = &Node{
		role:            NodeRole_Follower,
		currentTerm:     0,
		electionTimeout: 0,
		myAddr:          parseAddress(Iam),
		peers:           parseAddresses(peers),
	}
	return ins
}

func GetNodeInstance() *Node {
	return ins
}

func (n *Node) resetElectionTimeout() {
	lock.Lock()
	n.electionTimeout = 0
	lock.Unlock()
}

func (n *Node) incElectionTimeout() {
	lock.Lock()
	n.electionTimeout += INC_ELECTION_TIMEOUT
	lock.Unlock()
}

func (n *Node) setRole(r NodePole) {
	lock.Lock()
	n.role = r
	lock.Unlock()
}

func (n *Node) incCurrentTerm() {
	n.currentTerm++
}

func (n *Node) Run() {
	go n.startRaftServer()
	n.mainLoop()
}

func (n *Node) mainLoop() {
	for {
		fmt.Printf("I [%s:%s] am a %s...\n", n.myAddr.name, n.myAddr.port, n.role.ToString())
		time.Sleep(time.Second)
		if n.electionTimeout > ELECTION_TIMEOUT {
			n.resetElectionTimeout()
			n.gotoElectionPeriod()
		} else {
			n.incElectionTimeout()
		}
	}
}

func (n *Node) gotoElectionPeriod() {
	fmt.Printf("I [%s:%s] starts to electe ...\n", n.myAddr.name, n.myAddr.port)
	n.incCurrentTerm()
	n.setRole(NodeRole_Candidate)
}

type server struct {
	raft_rpc.UnimplementedRaftServiceServer
}

func (s *server) TellMyHeartBeatToFollower(ctx context.Context, in *raft_rpc.HeartBeatRequest) (*raft_rpc.HeartBeatReply, error) {
	log.Printf("Received: %v", in.GetName())
	GetNodeInstance().resetElectionTimeout()
	GetNodeInstance().setRole(NodeRole_Follower)
	return &raft_rpc.HeartBeatReply{Message: "Received the heart beat of leader: " + in.GetName()}, nil
}

func (n *Node) startRaftServer() {
	lis, err := net.Listen("tcp", ":"+n.myAddr.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	raft_rpc.RegisterRaftServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
