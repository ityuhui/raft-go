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
	// eclection
	role            NodePole
	currentTerm     int64
	electionTimeout int
	myAddr          *Address
	peers           []*Peer
	votedFor        string

	// log state machine
	log []*LogEntry
}

var ins *Node = nil
var lock sync.Mutex

func NewNodeInstance(I string, peers string) *Node {

	ins = &Node{
		role:            NodeRole_Follower,
		currentTerm:     0,
		electionTimeout: 0,
		myAddr:          ParseAddress(I),
		peers:           InitPeers(peers),
		votedFor:        "",
	}
	return ins
}

func GetNodeInstance() *Node {
	return ins
}

func (n *Node) ResetElectionTimeout() {
	lock.Lock()
	n.electionTimeout = 0
	lock.Unlock()
}

func (n *Node) IncElectionTimeout() {
	lock.Lock()
	n.electionTimeout += INC_ELECTION_TIMEOUT
	lock.Unlock()
}

func (n *Node) SetRole(r NodePole) {
	lock.Lock()
	n.role = r
	lock.Unlock()
}

func (n *Node) IncCurrentTerm() {
	n.currentTerm++
}

func (n *Node) GetCurrentTerm() int64 {
	return n.currentTerm
}

func (n *Node) SetCurrentTerm(term int64) {
	n.currentTerm = term
}

func (n *Node) GetMyAddress() *Address {
	return n.myAddr
}

func (n *Node) GetVotedFor() string {
	return n.votedFor
}

func (n *Node) SetVotedFor(votedFor string) {
	n.votedFor = votedFor
}

func (n *Node) Run() {
	go n.startRaftServer()
	n.MainLoop()
}

func (n *Node) MainLoop() {
	for {
		fmt.Printf("I [%s] am a %s...\n", n.GetMyAddress().GenerateUName(), n.role.ToString())
		time.Sleep(time.Second)

		switch n.role {
		case NodeRole_Leader:
			n.SendHeartBeatToFollowers()
		case NodeRole_Follower:
			if n.electionTimeout > ELECTION_TIMEOUT {
				n.ResetElectionTimeout()
				n.GotoElectionPeriod()
			} else {
				n.IncElectionTimeout()
			}
		}
	}
}

func (n *Node) SendHeartBeatToFollowers() {
	for _, peer := range n.peers {
		n.sendHeartBeatToFollower(peer.GetAddress())
	}
}

func (n *Node) GotoElectionPeriod() {
	fmt.Printf("I [%s:%s] starts to electe ...\n", n.myAddr.name, n.myAddr.port)
	n.IncCurrentTerm()
	n.SetRole(NodeRole_Candidate)
	n.SetVotedFor(n.GetMyAddress().GenerateUName())
	halfNumOfNodes := (float64(len(n.peers)) + 1.0) / 2.0
	numOfAgree := 1.0 // vote to myself

	agreeMap := make(map[string]bool)
	for _, peer := range n.peers {
		agree := n.sendVoteRequest(peer.GetAddress())
		peerName := peer.GetAddress().GenerateUName()
		agreeMap[peerName] = agree
	}

	for _, v := range agreeMap {
		if v == true {
			numOfAgree++
		}
	}
	if numOfAgree > halfNumOfNodes {
		n.SetRole(NodeRole_Leader)
	}
}

type server struct {
	raft_rpc.UnimplementedRaftServiceServer
}

func (s *server) AppendEntries(ctx context.Context, in *raft_rpc.AppendRequest) (*raft_rpc.AppendReply, error) {
	log.Printf("I [%v] received heart beat from leader: %v, term %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetName(), in.GetTerm())
	var message string
	if in.GetTerm() >= GetNodeInstance().GetCurrentTerm() {
		GetNodeInstance().ResetElectionTimeout()
		GetNodeInstance().SetRole(NodeRole_Follower)
		message = "[" + GetNodeInstance().GetMyAddress().GenerateUName() + "] accepted the heart beat from leader " + in.GetName()
	} else {
		message = "[" + GetNodeInstance().GetMyAddress().GenerateUName() + "] have refused the heart beat from " + in.GetName()
	}
	log.Printf("I %v", message)
	return &raft_rpc.AppendReply{Message: message}, nil
}

func (s *server) RequestVote(ctx context.Context, in *raft_rpc.VoteRequest) (*raft_rpc.VoteReply, error) {
	log.Printf("I [%v] received vote request from candinate: %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetCandidateId())
	candinateTerm := in.GetTerm()
	candinateID := in.GetCandidateId()
	agree := false

	if candinateTerm > GetNodeInstance().GetCurrentTerm() &&
		(GetNodeInstance().GetVotedFor() == "" ||
			candinateID == GetNodeInstance().GetVotedFor()) {
		GetNodeInstance().SetVotedFor(candinateID)
		GetNodeInstance().SetCurrentTerm(candinateTerm)
		GetNodeInstance().ResetElectionTimeout()
		agree = true
		log.Printf("I [%v] agreed the vote request from candinate: %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetCandidateId())
	} else {
		agree = false
		log.Printf("I [%v] disagreed the vote request from candinate: %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetCandidateId())
	}

	return &raft_rpc.VoteReply{Term: candinateTerm, VoteGranted: agree}, nil
}

func (n *Node) sendVoteRequest(addr *Address) bool {
	log.Printf("Begin to send vote request to: %v", addr.GenerateUName())
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
	r, err := c.RequestVote(ctx, &raft_rpc.VoteRequest{CandidateId: n.GetMyAddress().GenerateUName(),
		Term: n.GetCurrentTerm()})
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
	r, err := c.AppendEntries(ctx, &raft_rpc.AppendRequest{Name: n.GetMyAddress().GenerateUName(), Term: n.GetCurrentTerm()})
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
