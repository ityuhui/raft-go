package main

type NodePole int32

const (
	NodeRole_Leader    NodePole = 1
	NodeRole_Candidate NodePole = 2
	NodeRole_Follower  NodePole = 3
)

func (r NodePole) String() string {
	switch r {
	case NodeRole_Leader:
		return "Leader"
	case NodeRole_Candidate:
		return "Candidate"
	case NodeRole_Follower:
		return "Follower"
	default:
		return "UNKNOWN"
	}
}
