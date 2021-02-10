package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"raft-go/common"
	"raft-go/raft_rpc"
	"strconv"
	"strings"
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

	// node log
	nodeLog []*LogEntry
	// counter for node log
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
		commitIndex:     0,
		lastApplied:     0,
		stateMachine:    NewStateMachineInstance(),
	}
	return nodeInstance
}

func (n *Node) getNodeLogEntryTermByIndex(index int64) (int64, error) {
	logEntry := n.GetNodeLog()[index]
	if logEntry != nil {
		return logEntry.term, nil
	} else {
		return 0, errors.New("The log entry does not exist.")
	}
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

func (n *Node) GetLastApplied() int64 {
	return n.lastApplied
}

func (n *Node) GetCommitIndex() int64 {
	return n.commitIndex
}

func (n *Node) GetNodeLog() []*LogEntry {
	return n.nodeLog
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
		go n.sendHeartBeatOrAppendLogToFollower(peer, 0, 0)
	}
}

func (n *Node) AppendLogToFollowers(prevLogIndex int64, prevLogTerm int64) {
	for _, peer := range n.peers {
		go n.sendHeartBeatOrAppendLogToFollower(peer, prevLogIndex, prevLogTerm)
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

func (n *Node) addToNodeLog(log string) int64 {
	entry := &LogEntry{
		term: n.GetCurrentTerm(),
		text: log,
	}
	n.nodeLog = append(n.nodeLog, entry)
	return int64(len(n.nodeLog))
}

func (n *Node) applyNodeLogToStateMachine() {
	for n.commitIndex > n.lastApplied {
		n.lastApplied++
		logentry := n.nodeLog[n.lastApplied]
		n.executeStateMachineCommand(logentry.text)
	}
}

func (n *Node) executeStateMachineCommand(cmd string) {
	cmds := strings.Split(cmd, "=")
	key := cmds[0]
	value, _ := strconv.ParseInt(cmds[1], 10, 64)
	n.stateMachine.Set(key, value)
}

type server struct {
	raft_rpc.UnimplementedRaftServiceServer
}

func (s *server) AppendEntries(ctx context.Context, in *raft_rpc.AppendRequest) (*raft_rpc.AppendReply, error) {
	success := false
	var message string

	candinateTerm := in.GetTerm()
	if candinateTerm >= GetNodeInstance().GetCurrentTerm() {
		if in.GetLogEntries() == nil {
			GetNodeInstance().ResetElectionTimeout()
			GetNodeInstance().SetRole(NodeRole_Follower)
			GetNodeInstance().SetCurrentTerm(candinateTerm)
			log.Printf("I [%v] received heart beat from leader: %v, term %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetLeaderId(), in.GetTerm())
		} else {
			log.Printf("I [%v] am required to append log entry from leader: %v, term %v", GetNodeInstance().GetMyAddress().GenerateUName(), in.GetLeaderId(), in.GetTerm())

		}
		success = true
		message = "[" + GetNodeInstance().GetMyAddress().GenerateUName() + "] accepted the append from leader " + in.GetLeaderId()
	} else {
		message = "[" + GetNodeInstance().GetMyAddress().GenerateUName() + "] have refused the append request from " + in.GetLeaderId()
		success = false
	}
	log.Printf("I %v", message)

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
			message = "The command is executed."
		}
	} else if in.GetMode() == common.COMMANDMODE_SET.ToString() {
		node := GetNodeInstance()
		prevLogIndex := node.addToNodeLog(in.GetText())
		prevLogTerm, rc := node.getNodeLogEntryTermByIndex(prevLogIndex)
		if rc == nil {
			node.AppendLogToFollowers(prevLogIndex, prevLogTerm)
		}
		node.applyNodeLogToStateMachine()
		if node.GetLastApplied() == node.GetCommitIndex() {
			success = true
			message = "The command is executed."
		}
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

func (n *Node) prepareNodeLogToAppend(peer *Peer) []*LogEntry {
	return []
}

func (n *Node) sendHeartBeatOrAppendLogToFollower(peer *Peer, prevLogIndex int64, prevLogTerm int64) {
	// Set up a connection to the server.
	addr := peer.GetAddress()

	conn, err := grpc.Dial(addr.Name+":"+addr.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := raft_rpc.NewRaftServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AppendEntries(ctx, &raft_rpc.AppendRequest{
		Term:         n.GetCurrentTerm(),
		LeaderId:     n.GetMyAddress().GenerateUName(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		LogEntries:   n.prepareNodeLogToAppend(peer),
		LeaderCommit: n.GetCommitIndex(),
	})
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
