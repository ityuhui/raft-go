package main

type NodePole int32

const (
	NodeRole_Leader    NodePole = 1
	NodeRole_Candidate NodePole = 2
	NodeRole_Follower  NodePole = 3
)

func (r NodePole) ToString() string {
	switch r {
	case NodeRole_Leader:
		return "leader"
	case NodeRole_Candidate:
		return "candidate"
	case NodeRole_Follower:
		return "follower"
	default:
		return "UNKNOWN"
	}
}
