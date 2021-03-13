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
	electionTimeout    = 5
	incElectionTimeout = 1
)

//Node : data structure for node
type Node struct {
	// eclection
	role            NodeRole
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

//NewNodeInstance : create an instance of node
func NewNodeInstance(I string, peers string) *Node {

	nodeInstance = &Node{
		role:            NodeRoleFollower,
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

func getNodeInstance() *Node {
	return nodeInstance
}

func (n *Node) resetElectionTimeout() {
	nodeLock.Lock()
	n.electionTimeout = 0
	nodeLock.Unlock()
}

func (n *Node) incElectionTimeout() {
	nodeLock.Lock()
	n.electionTimeout += incElectionTimeout
	nodeLock.Unlock()
}

func (n *Node) setRole(r NodeRole) {
	nodeLock.Lock()
	n.role = r
	nodeLock.Unlock()
}

func (n *Node) incCurrentTerm() {
	n.currentTerm++
}

func (n *Node) getCurrentTerm() int64 {
	return n.currentTerm
}

func (n *Node) setCurrentTerm(term int64) {
	n.currentTerm = term
}

func (n *Node) getMyAddress() *common.Address {
	return n.myAddr
}

func (n *Node) getVotedFor() string {
	return n.votedFor
}

func (n *Node) setVotedFor(votedFor string) {
	n.votedFor = votedFor
}

func (n *Node) getLastApplied() int64 {
	return n.lastApplied
}

func (n *Node) getCommitIndex() int64 {
	return n.commitIndex
}

func (n *Node) setCommitIndex(index int64) {
	n.commitIndex = index
}

func (n *Node) getNodeLog() []*LogEntry {
	return n.nodeLog
}

//Run : node runs
func (n *Node) Run() {
	go n.startRaftServer()
	n.mainLoop()
}

func (n *Node) mainLoop() {
	for {
		fmt.Printf("I [%s] am a %s...\n", n.getMyAddress().GenerateUName(), n.role.ToString())
		time.Sleep(time.Second)

		switch n.role {
		case NodeRoleLeader:
			n.sendHeartBeatToFollowers()
		case NodeRoleFollower:
			if n.electionTimeout > electionTimeout {
				n.resetElectionTimeout()
				n.gotoElectionPeriod()
			} else {
				n.incElectionTimeout()
			}
		}

		n.applyNodeLogToStateMachine()
	}
}

func (n *Node) sendHeartBeatToFollowers() {
	for _, peer := range n.peers {
		go n.sendHeartBeatOrAppendLogToFollower(peer, 0, 0)
	}
}

func (n *Node) appendLogToFollower(peer *Peer, prevLogIndex int64, prevLogTerm int64) {
	var err error
	for ; err != nil; err = n.sendHeartBeatOrAppendLogToFollower(peer, prevLogIndex, prevLogTerm) {
		peer.DecreaseNextIndex()
	}
	peer.UpdateNextIndexAndMatchIndex()
}

func (n *Node) appendLogToFollowers(prevLogIndex int64, prevLogTerm int64) {
	for _, peer := range n.peers {
		go n.appendLogToFollower(peer, prevLogIndex, prevLogTerm)
	}
}

func (n *Node) gotoElectionPeriod() {
	fmt.Printf("I [%s:%s] starts to electe ...\n", n.myAddr.Name, n.myAddr.Port)
	n.incCurrentTerm()
	n.setRole(NodeRoleCandidate)
	n.setVotedFor(n.getMyAddress().GenerateUName())
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
		n.setRole(NodeRoleLeader)
		n.initPeersNextIndexAndMatchIndex()
	}
}

func (n *Node) initPeersNextIndexAndMatchIndex() {
	for _, peer := range n.peers {
		peer.SetNextIndex(n.getLastNodeLogEntryIndex() + 1)
		peer.SetMatchIndex(0)
	}
}

// node log operations
func (n *Node) getNodeLogLength() int64 {
	return int64(len(n.getNodeLog()))
}

func (n *Node) getNodeLogEntryTermByIndex(index int64) (int64, error) {
	logEntry := n.getNodeLog()[index]
	if logEntry != nil {
		return logEntry.Term, nil
	}
	return 0, errors.New("The log entry does not exist")
}

func (n *Node) getLastNodeLogEntryIndex() int64 {
	return n.getNodeLogLength()
}

func (n *Node) addToNodeLog(entry *LogEntry) int64 {
	n.nodeLog = append(n.nodeLog, entry)
	return n.getNodeLogLength()
}

func (n *Node) addCmdToNodeLog(log string) int64 {
	entry := &LogEntry{
		Term: n.getCurrentTerm(),
		Text: log,
	}
	n.nodeLog = append(n.nodeLog, entry)
	return int64(len(n.nodeLog))
}

//applyNodeLogToStateMachine : execute the command from node log in state machine
func (n *Node) applyNodeLogToStateMachine() {
	for n.getCommitIndex() > n.lastApplied {
		n.lastApplied++
		logentry := n.nodeLog[n.lastApplied]
		n.executeStateMachineCommand(logentry.Text)
	}
}

func (n *Node) isMajorityMatchIndexGreaterThanN(newCI int64) bool {
	return true
}

//UpdateMyCommitIndexWhenIamLeader : update commitIndex when node is a leader
func (n *Node) UpdateMyCommitIndexWhenIamLeader() {
	for newCI := n.getCommitIndex() + 1; n.isMajorityMatchIndexGreaterThanN(newCI); newCI++ {
		term, err := n.getNodeLogEntryTermByIndex(newCI)
		if err == nil && term == n.getCurrentTerm() {
			n.setCommitIndex(newCI)
		}
	}
}

//UpdateMyCommitIndexWhenIamFollower : update commitIndex when node is a follower
func (n *Node) UpdateMyCommitIndexWhenIamFollower(leaderCommit int64) {
	if leaderCommit > n.getCommitIndex() {
		lastIndex := n.getLastNodeLogEntryIndex()
		if leaderCommit <= lastIndex {
			n.setCommitIndex(leaderCommit)
		} else {
			n.setCommitIndex(lastIndex)
		}
	}
}

func (n *Node) deleteLogEntryAndItsFollowerInNodeLog(index int64) error {
	if index > int64(n.getNodeLogLength()) {
		return errors.New("Cannot find the log entry with the index: " + strconv.FormatInt(index, 10))
	}
	var i int64
	nodeLog := n.getNodeLog()
	for i = index; i < n.getNodeLogLength(); i++ {
		nodeLog[i] = nil
	}
	return nil
}

func (n *Node) appendEntryFromLeaderToMyNodeLog(logEntries []*raft_rpc.LogEntry) {
	var i int
	for i = 0; i < len(logEntries); i++ {
		logEntry := &LogEntry{
			Term: logEntries[i].Term,
			Text: logEntries[i].Text,
		}
		n.addToNodeLog(logEntry)
	}
}

func (n *Node) prepareNodeLogToAppend(peer *Peer) []*raft_rpc.LogEntry {
	myLastLogIndex := n.getLastNodeLogEntryIndex()
	peerNextIndex := peer.GetNextIndex()
	if myLastLogIndex >= peerNextIndex {
		toAppend := []*raft_rpc.LogEntry{}
		var i int64
		for i = 0; i < peerNextIndex; i++ {
			logEntry := &raft_rpc.LogEntry{
				Term: n.getNodeLog()[i].Term,
				Text: n.getNodeLog()[i].Text,
			}
			toAppend = append(toAppend, logEntry)
		}
		return toAppend
	}
	return nil
}

// append log
func (n *Node) sendHeartBeatOrAppendLogToFollower(peer *Peer, prevLogIndex int64, prevLogTerm int64) error {
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
		Term:         n.getCurrentTerm(),
		LeaderId:     n.getMyAddress().GenerateUName(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		LogEntries:   n.prepareNodeLogToAppend(peer),
		LeaderCommit: n.getCommitIndex(),
	})
	if err != nil {
		log.Fatalf("could not tell my heart beat: %v", err)
	}
	log.Printf("Telling: %s", r.GetMessage())
	return nil
}

// state machine operations
func (n *Node) executeStateMachineCommand(cmd string) {
	cmds := strings.Split(cmd, "=")
	key := cmds[0]
	value, _ := strconv.ParseInt(cmds[1], 10, 64)
	n.stateMachine.Set(key, value)
}

// vote
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
	r, err := c.RequestVote(ctx, &raft_rpc.VoteRequest{CandidateId: n.getMyAddress().GenerateUName(),
		Term: n.getCurrentTerm()})
	if err != nil {
		log.Fatalf("could not request to vote: %v", err)
	}
	log.Printf("Get voteGranted: %v", r.GetVoteGranted())
	return r.GetVoteGranted()
}

// raft server
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

type server struct {
	raft_rpc.UnimplementedRaftServiceServer
}

func (s *server) AppendEntries(ctx context.Context, in *raft_rpc.AppendRequest) (*raft_rpc.AppendReply, error) {
	success := false
	var message string
	var err error = nil

	leaderID := in.GetLeaderId()
	leaderTerm := in.GetTerm()
	leaderCommit := in.GetLeaderCommit()
	myName := getNodeInstance().getMyAddress().GenerateUName()

	if leaderTerm >= getNodeInstance().getCurrentTerm() {
		logEntries := in.GetLogEntries()
		if logEntries == nil {
			log.Printf("I [%v] received heart beat from leader: %v, term %v", myName, leaderID, leaderTerm)
			getNodeInstance().resetElectionTimeout()
			getNodeInstance().setRole(NodeRoleFollower)
			getNodeInstance().setCurrentTerm(leaderTerm)
			success = true
			message = "[" + myName + "] accepted the heart beat from leader " + leaderID
		} else {
			log.Printf("I [%v] am required to append log entry from leader: %v, term %v", myName, leaderID, leaderTerm)
			logEntryIndex := in.GetPrevLogIndex()
			logEntryTerm, rc := getNodeInstance().getNodeLogEntryTermByIndex(logEntryIndex)
			if rc != nil {
				if logEntryTerm != in.GetPrevLogTerm() {
					message = "[" + myName + "] has a log entry with the index [" + fmt.Sprint(logEntryIndex) + "], but its term is [" + fmt.Sprint(logEntryTerm) + "], not [" + fmt.Sprint(in.GetPrevLogTerm()) + "] from append request."
				}
				err = getNodeInstance().deleteLogEntryAndItsFollowerInNodeLog(logEntryIndex)
				if err != nil {
					success = false
					message = err.Error()
				} else {
					getNodeInstance().appendEntryFromLeaderToMyNodeLog(logEntries)
					getNodeInstance().UpdateMyCommitIndexWhenIamFollower(leaderCommit)
					success = true
					message = "[" + myName + "] have appended the log entries from " + leaderID + " successfully."
				}
			} else {
				success = false
				message = "[" + myName + "] does not have the log entry with index [" + fmt.Sprint(logEntryIndex) + "]."
				err = errors.New(message)
			}
		}
	} else {
		message = "[" + myName + "] have refused the append request from " + leaderID
		success = false
		err = errors.New(message)
	}
	log.Printf("I %v", message)

	return &raft_rpc.AppendReply{Term: getNodeInstance().getCurrentTerm(), Success: success, Message: message}, err
}

func (s *server) RequestVote(ctx context.Context, in *raft_rpc.VoteRequest) (*raft_rpc.VoteReply, error) {
	log.Printf("I [%v] received vote request from candinate: %v", getNodeInstance().getMyAddress().GenerateUName(), in.GetCandidateId())
	candinateTerm := in.GetTerm()
	candinateID := in.GetCandidateId()
	agree := false

	if candinateTerm > getNodeInstance().getCurrentTerm() &&
		(getNodeInstance().getVotedFor() == "" ||
			candinateID == getNodeInstance().getVotedFor()) {
		getNodeInstance().setVotedFor(candinateID)
		getNodeInstance().setCurrentTerm(candinateTerm)
		getNodeInstance().resetElectionTimeout()
		agree = true
		log.Printf("I [%v] agreed the vote request from candinate: %v", getNodeInstance().getMyAddress().GenerateUName(), in.GetCandidateId())
	} else {
		agree = false
		log.Printf("I [%v] disagreed the vote request from candinate: %v", getNodeInstance().getMyAddress().GenerateUName(), in.GetCandidateId())
	}

	return &raft_rpc.VoteReply{Term: candinateTerm, VoteGranted: agree}, nil
}

func (s *server) ExecuteCommand(ctx context.Context, in *raft_rpc.ExecuteCommandRequest) (*raft_rpc.ExecuteCommandReply, error) {
	log.Printf("I [%v] am requested to execute command <%v> from client.", getNodeInstance().getMyAddress().GenerateUName(), in.GetMode()+in.GetText())
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
		node := getNodeInstance()
		prevLogIndex := node.addCmdToNodeLog(in.GetText())
		prevLogTerm, rc := node.getNodeLogEntryTermByIndex(prevLogIndex)
		if rc == nil {
			node.appendLogToFollowers(prevLogIndex, prevLogTerm)
		}
		node.UpdateMyCommitIndexWhenIamLeader()
		node.applyNodeLogToStateMachine()
		if node.getLastApplied() == node.getCommitIndex() {
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
