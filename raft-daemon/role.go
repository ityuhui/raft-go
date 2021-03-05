package main

type NodePole int32

const (
	NodeRoleLeader    NodePole = 1
	NodeRoleCandidate NodePole = 2
	NodeRoleFollower  NodePole = 3
)

func (r NodePole) ToString() string {
	switch r {
	case NodeRoleLeader:
		return "leader"
	case NodeRoleCandidate:
		return "candidate"
	case NodeRoleFollower:
		return "follower"
	default:
		return "UNKNOWN"
	}
}
