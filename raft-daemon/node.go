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
	nodeLogStartIndex  = 1
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
var gWg4appendingLogToFollowers sync.WaitGroup

//NewNodeInstance : create an instance of node
func NewNodeInstance(I string, peers string) *Node {

	nodeInstance = &Node{
		role:            NodeRoleFollower,
		currentTerm:     0,
		electionTimeout: 0,
		myAddr:          common.ParseAddress(I),
		peers:           InitPeers(peers),
		votedFor:        "",
		nodeLog:         initNodeLog(),
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

func initNodeLog() []*LogEntry {
	var nodeLog []*LogEntry
	/* The first log entry (index:0) is place-holder,
	 * the valid log entry will begin at the index of "1"
	 */
	entry := &LogEntry{
		Term: 0,
		Text: "",
	}
	nodeLog = append(nodeLog, entry)
	return nodeLog
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
	fname := "appendLogToFollower"
	log.Printf("%v: enter...prevLogIndex=%v, prevLogTerm=%v, ", fname, prevLogIndex, prevLogTerm)
	err := n.sendHeartBeatOrAppendLogToFollower(peer, prevLogIndex, prevLogTerm)
	for ; err != nil; err = n.sendHeartBeatOrAppendLogToFollower(peer, prevLogIndex, prevLogTerm) {
		peer.DecreaseNextIndex()
		log.Printf("%v: [%v] nextIndex=%v, ", fname, peer.GetAddress().GenerateUName(), peer.GetNextIndex())
	}
	peer.UpdateNextIndexAndMatchIndex()
	log.Printf("%v: [%v] nextIndex=%v, matchIndex=%v, ", fname, peer.GetAddress().GenerateUName(), peer.GetNextIndex(), peer.GetMatchIndex())
	gWg4appendingLogToFollowers.Done()
}

func (n *Node) appendLogToFollowers(prevLogIndex int64, prevLogTerm int64) {
	fname := "appendLogToFollowers"
	log.Printf("%v: enter...", fname)
	for _, peer := range n.peers {
		log.Printf("%v: will send to [%v]", fname, peer.GetAddress().GenerateUName())
		gWg4appendingLogToFollowers.Add(1)
		go n.appendLogToFollower(peer, prevLogIndex, prevLogTerm)
	}
	gWg4appendingLogToFollowers.Wait()
	log.Printf("%v: exit...", fname)
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
	//fname := "getNodeLogLength()"
	//log.Printf("%v: nodeLength=%v", fname, int64(len(n.getNodeLog())))
	/*for _, logitem := range n.nodeLog {
		log.Printf("%v: %v-%v", fname, logitem.Text, logitem.Term)
	}*/

	return int64(len(n.getNodeLog()))
}

func (n *Node) getNodeLogEntryTermByIndex(index int64) (int64, error) {
	if index < nodeLogStartIndex || index > n.getLastNodeLogEntryIndex() {
		return 0, errors.New("The log entry does not exist")
	}
	logEntry := n.getNodeLog()[index]
	return logEntry.Term, nil
}

func (n *Node) getLastNodeLogEntryIndex() int64 {
	return n.getNodeLogLength() - 1
}

func (n *Node) addEntryToNodeLog(entry *LogEntry) int64 {
	n.nodeLog = append(n.nodeLog, entry)
	return n.getLastNodeLogEntryIndex()
}

func (n *Node) addStringToNodeLog(log string) int64 {
	entry := &LogEntry{
		Term: n.getCurrentTerm(),
		Text: log,
	}
	n.nodeLog = append(n.nodeLog, entry)
	return n.getLastNodeLogEntryIndex()
}

//applyNodeLogToStateMachine : execute the command from node log in state machine
func (n *Node) applyNodeLogToStateMachine() {
	fname := "applyNodeLogToStateMachine()"
	log.Printf("%v: a commitIndex=%v, lastApplied=%v", fname, n.getCommitIndex(), n.lastApplied)
	for n.getCommitIndex() > n.lastApplied {
		log.Printf("%v: b commitIndex=%v, lastApplied=%v", fname, n.getCommitIndex(), n.lastApplied)
		n.lastApplied++
		logentry := n.nodeLog[n.lastApplied]
		n.executeStateMachineCommand(logentry.Text)
	}
}

func (n *Node) isMajorityMatchIndexGreaterThanOrEqualsTo(newCI int64) bool {
	numOfGreaterOrEqual := 0.0
	for _, peer := range n.peers {
		if peer.GetMatchIndex() >= newCI {
			numOfGreaterOrEqual++
		}
	}
	halfNumOfNodes := (float64(len(n.peers)) + 1.0) / 2.0
	if numOfGreaterOrEqual > halfNumOfNodes {
		return true
	}
	return false
}

//UpdateMyCommitIndexWhenIamLeader : update commitIndex when node is a leader
func (n *Node) UpdateMyCommitIndexWhenIamLeader() {
	fname := "UpdateMyCommitIndexWhenIamLeader()"

	for newCI := n.getCommitIndex() + 1; n.isMajorityMatchIndexGreaterThanOrEqualsTo(newCI); newCI++ {
		log.Printf("%v:newCI=%v", fname, newCI)
		term, err := n.getNodeLogEntryTermByIndex(newCI)
		if err == nil && term == n.getCurrentTerm() {
			n.setCommitIndex(newCI)
			break
		}
	}
	log.Printf("%v: My commitIndex=%v", fname, n.getCommitIndex())
}

//UpdateMyCommitIndexWhenIamFollower : update commitIndex when node is a follower
func (n *Node) UpdateMyCommitIndexWhenIamFollower(leaderCommit int64) {
	fname := "UpdateMyCommitIndexWhenIamFollower()"
	log.Printf("%v: leaderCommit=%v, commitIndex=%v", fname, leaderCommit, n.getCommitIndex())
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
	if index > n.getLastNodeLogEntryIndex() {
		return errors.New("Cannot find the log entry with the index: " + strconv.FormatInt(index, 10))
	}
	var i int64
	nodeLog := n.getNodeLog()
	for i = index; i <= n.getLastNodeLogEntryIndex(); i++ {
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
		n.addEntryToNodeLog(logEntry)
	}
}

func (n *Node) prepareNodeLogToAppend(peer *Peer) []*raft_rpc.LogEntry {
	fname := "prepareNodeLogToAppend()"
	myLastLogIndex := n.getLastNodeLogEntryIndex()
	peerNextIndex := peer.GetNextIndex()
	log.Printf("%v: myLastLogIndex=%v, peerNextIndex=%v", fname, myLastLogIndex, peerNextIndex)
	if myLastLogIndex >= peerNextIndex {
		toAppend := []*raft_rpc.LogEntry{}
		var i int64
		for i = nodeLogStartIndex; i <= peerNextIndex; i++ {
			logEntry := &raft_rpc.LogEntry{
				Term: n.getNodeLog()[i].Term,
				Text: n.getNodeLog()[i].Text,
			}
			toAppend = append(toAppend, logEntry)
		}
		log.Printf("%v: The length of append log=%v", fname, len(toAppend))
		return toAppend
	}
	log.Printf("%v: No log to append to [%v]", fname, peer.GetAddress().GenerateUName())
	return nil
}

// append log
func (n *Node) sendHeartBeatOrAppendLogToFollower(peer *Peer, prevLogIndex int64, prevLogTerm int64) error {
	fname := "sendHeartBeatOrAppendLogToFollower()"
	log.Printf("%v: enter...", fname)
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
		log.Fatalf("%v: could not tell my heart beat or append entry: %v", fname, err)
	}
	log.Printf("%v: Telling: %s", fname, r.GetMessage())
	return nil
}

// state machine operations
func (n *Node) executeStateMachineCommand(cmd string) {
	fname := "executeStateMachineCommand()"
	log.Printf("%v: enter", fname)
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
	fname := "ExecuteCommand()"
	log.Printf("I [%v] am requested to execute command <%v %v> from client.", getNodeInstance().getMyAddress().GenerateUName(), in.GetMode(), in.GetText())
	success := false
	var value int64 = 0
	var rc error = nil
	message := ""
	if in.GetMode() == common.CommandModeGet.ToString() {
		value, rc = GetStateMachineInstance().Get(in.GetText())
		if rc == nil {
			success = true
			message = "The command is executed."
		}
	} else if in.GetMode() == common.CommandModeSet.ToString() {
		node := getNodeInstance()
		prevLogIndex := node.addStringToNodeLog(in.GetText()) - 1
		log.Printf("%v: prevLogIndex=%v", fname, prevLogIndex)
		prevLogTerm, rc := node.getNodeLogEntryTermByIndex(prevLogIndex)
		if rc == nil {
			log.Printf("%v: prevLogTerm=%v", fname, prevLogTerm)
			node.appendLogToFollowers(prevLogIndex, prevLogTerm)
		} else if 0 == prevLogIndex {
			prevLogTerm = node.getCurrentTerm()
			log.Printf("%v: the first LogTerm=%v", fname, prevLogTerm)
			node.appendLogToFollowers(prevLogIndex, prevLogTerm)
		} else {
			log.Printf("%v: err of getNodeLogEntryTermByIndex=%v", fname, rc.Error())
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
