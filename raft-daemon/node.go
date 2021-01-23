package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"raft-go/common"
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
	myAddr          *common.Address
	peers           []*Peer
	votedFor        string

	// log state machine
	log []*LogEntry

	// log commit
	commitIndex int64
	lastApplied int64

	// state machine
	stateMachine *StateMachine
}

var nodeInstance *Node = nil
var nodeLock sync.Mutex

func NewNodeInstance(I string, peers string) *Node {

	nodeInstance = &Node{
		role:            NodeRole_Follower,
		currentTerm:     0,
		electionTimeout: 0,
		myAddr:          common.ParseAddress(I),
		peers:           InitPeers(peers),
		votedFor:        "",
		stateMachine:    NewStateMachineInstance(),
	}
	return nodeInstance
}

func GetNodeInstance() *Node {
	return nodeInstance
}

func (n *Node) ResetElectionTimeout() {
	nodeLock.Lock()
	n.electionTimeout = 0
	nodeLock.Unlock()
}

func (n *Node) IncElectionTimeout() {
	nodeLock.Lock()
	n.electionTimeout += INC_ELECTION_TIMEOUT
	nodeLock.Unlock()
}

func (n *Node) SetRole(r NodePole) {
	nodeLock.Lock()
	n.role = r
	nodeLock.Unlock()
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

func (n *Node) GetMyAddress() *common.Address {
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
	fmt.Printf("I [%s:%s] starts to electe ...\n", n.myAddr.Name, n.myAddr.Port)
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
	success := false
	var message string
	if in.GetLogEntries() == nil {
		log.Printf("I [%v] received heart beat from leader: %v, term %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetLeaderId(), in.GetTerm())
		candinateTerm := in.GetTerm()
		if candinateTerm >= GetNodeInstance().GetCurrentTerm() {
			GetNodeInstance().ResetElectionTimeout()
			GetNodeInstance().SetRole(NodeRole_Follower)
			GetNodeInstance().SetCurrentTerm(candinateTerm)
			success = true
			message = "[" + GetNodeInstance().GetMyAddress().GenerateUName() + "] accepted the heart beat from leader " + in.GetLeaderId()
		} else {
			message = "[" + GetNodeInstance().GetMyAddress().GenerateUName() + "] have refused the heart beat from " + in.GetLeaderId()
			success = false
		}
		log.Printf("I %v", message)
	} else {
		success = true
		log.Printf("I [%v] am required to append log entry from leader: %v, term %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetLeaderId(), in.GetTerm())
	}
	return &raft_rpc.AppendReply{Term: GetNodeInstance().GetCurrentTerm(), Success: success, Message: message}, nil
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

func (s *server) ExecuteCommand(ctx context.Context, in *raft_rpc.ExecuteCommandRequest) (*raft_rpc.ExecuteCommandReply, error) {
	log.Printf("I [%v] am requested to execute command <%v> from client.", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetMode()+in.GetText())
	success := false
	var value int64 = 0
	var rc error = nil
	message := ""
	if in.GetMode() == common.COMMANDMODE_GET.ToString() {
		value, rc = GetStateMachineInstance().Get(in.GetText())
		if rc == nil {
			success = true
		}
	} else if in.GetMode() == common.COMMANDMODE_SET.ToString() {

	} else {
		message = "The command is unknown."
		rc = errors.New(message)
	}
	return &raft_rpc.ExecuteCommandReply{
		Success: success,
		Value:   value,
		Message: message,
	}, rc
}

func (n *Node) sendVoteRequest(addr *common.Address) bool {
	log.Printf("Begin to send vote request to: %v", addr.GenerateUName())
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.Name+":"+addr.Port, grpc.WithInsecure(), grpc.WithBlock())
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

func (n *Node) sendHeartBeatToFollower(addr *common.Address) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.Name+":"+addr.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := raft_rpc.NewRaftServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AppendEntries(ctx, &raft_rpc.AppendRequest{LeaderId: n.GetMyAddress().GenerateUName(), Term: n.GetCurrentTerm()})
	if err != nil {
		log.Fatalf("could not tell my heart beat: %v", err)
	}
	log.Printf("Telling: %s", r.GetMessage())
}

func (n *Node) startRaftServer() {
	lis, err := net.Listen("tcp", ":"+n.myAddr.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	raft_rpc.RegisterRaftServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
