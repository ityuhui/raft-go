package main

type Peer struct {
	address    *Address
	nextIndex  int64
	matchIndex int64
}

func initPeers(peers string) []*Peer {
	addresses := parseAddresses(peers)

}

func (p *Peer) GetAddress() *Address {
	return p.address
}

func (p *Peer) GetNextIndex() int64 {
	return p.nextIndex
}

func (p *Peer) GetMatchIndex() int64 {
	return p.matchIndex
}
