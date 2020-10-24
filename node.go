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
	currentTerm     int64
	electionTimeout int
	myAddr          *Address
	peers           []*Address
	votedFor        string
}

var ins *Node = nil
var lock sync.Mutex

func NewNodeInstance(I string, peers string) *Node {

	ins = &Node{
		role:            NodeRole_Follower,
		currentTerm:     0,
		electionTimeout: 0,
		myAddr:          parseAddress(I),
		peers:           parseAddresses(peers),
		votedFor:        "",
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

func (n *Node) getCurrentTerm() int64 {
	return n.currentTerm
}

func (n *Node) getMyAddress() *Address {
	return n.myAddr
}

func (n *Node) getVotedFor() string {
	return n.votedFor
}

func (n *Node) setVotedFor(votedFor string) {
	n.votedFor = votedFor
}

func (n *Node) Run() {
	go n.startRaftServer()
	n.mainLoop()
}

func (n *Node) mainLoop() {
	for {
		fmt.Printf("I [%s] am a %s...\n", n.getMyAddress().generateUName(), n.role.ToString())
		time.Sleep(time.Second)

		switch n.role {
		case NodeRole_Leader:
			n.sendHeartBeatToFollowers()
		case NodeRole_Follower:
			if n.electionTimeout > ELECTION_TIMEOUT {
				n.resetElectionTimeout()
				n.gotoElectionPeriod()
			} else {
				n.incElectionTimeout()
			}
		}
	}
}

func (n *Node) sendHeartBeatToFollowers() {
	for _, peer := range n.peers {
		n.sendHeartBeatToFollower(peer)
	}
}

func (n *Node) gotoElectionPeriod() {
	fmt.Printf("I [%s:%s] starts to electe ...\n", n.myAddr.name, n.myAddr.port)
	n.incCurrentTerm()
	n.setRole(NodeRole_Candidate)

	var agreeMap map[string]bool
	for _, peer := range n.peers {
		agreeMap[peer.generateUName()] = n.sendVoteRequest(peer)
	}

	numOfHalfPeers := float64(len(n.peers)) / 2.0

	numOfAgree := 1.0 // vote to myself
	for _, v := range agreeMap {
		if v == true {
			numOfAgree++
		}
	}
	if numOfAgree > numOfHalfPeers {
		n.setRole(NodeRole_Leader)
	}
}

type server struct {
	raft_rpc.UnimplementedRaftServiceServer
}

func (s *server) TellMyHeartBeatToFollower(ctx context.Context, in *raft_rpc.HeartBeatRequest) (*raft_rpc.HeartBeatReply, error) {
	log.Printf("Received heart beat from leader: %v", in.GetName())
	GetNodeInstance().resetElectionTimeout()
	GetNodeInstance().setRole(NodeRole_Follower)
	return &raft_rpc.HeartBeatReply{Message: GetNodeInstance().getMyAddress().generateUName() + " received the heart beat."}, nil
}

func (s *server) RequestToVote(ctx context.Context, in *raft_rpc.VoteRequest) (*raft_rpc.VoteReply, error) {
	log.Printf("Received vote request from candinate: %v", in.GetCandidateId())
	candinateTerm := in.GetTerm()
	candinateID := in.GetCandidateId()
	agree := false

	if candinateTerm < GetNodeInstance().getCurrentTerm() {
		agree = false
	} else if GetNodeInstance().getVotedFor() == "" || candinateID == GetNodeInstance().getVotedFor() {
		agree = true
		GetNodeInstance().setVotedFor(candinateID)
	}

	return &raft_rpc.VoteReply{Term: candinateTerm, VoteGranted: agree}, nil
}

func (n *Node) sendVoteRequest(addr *Address) bool {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.name+":"+addr.port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := raft_rpc.NewRaftServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RequestToVote(ctx, &raft_rpc.VoteRequest{CandidateId: n.myAddr.generateUName(),
		Term: n.getCurrentTerm()})
	if err != nil {
		log.Fatalf("could not request to vote: %v", err)
	}
	log.Printf("Get voteGranted: %v", r.GetVoteGranted())
	return r.GetVoteGranted()
}

func (n *Node) sendHeartBeatToFollower(addr *Address) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.name+":"+addr.port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := raft_rpc.NewRaftServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.TellMyHeartBeatToFollower(ctx, &raft_rpc.HeartBeatRequest{Name: n.myAddr.generateUName()})
	if err != nil {
		log.Fatalf("could not tell my heart beat: %v", err)
	}
	log.Printf("Telling: %s", r.GetMessage())
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
