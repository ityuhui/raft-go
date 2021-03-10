package main

//NodeRole : role of node
type NodeRole int32

const (
	//NodeRoleLeader : leader
	NodeRoleLeader NodeRole = 1
	//NodeRoleCandidate : candidate
	NodeRoleCandidate NodeRole = 2
	//NodeRoleFollower : follower
	NodeRoleFollower NodeRole = 3
)

// ToString : convert NodeRole to string
func (r NodeRole) ToString() string {
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
