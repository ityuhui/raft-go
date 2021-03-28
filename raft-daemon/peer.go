package main

import (
	"log"
	"raft-go/common"
	"strings"
)

//Peer : The data structure for peer node
type Peer struct {
	address    *common.Address
	nextIndex  int64
	matchIndex int64
}

//InitPeers : Init peer nodes for a node
func InitPeers(addrStrings string) []*Peer {
	var peers []*Peer
	splited := strings.Split(addrStrings, ",")
	if len(splited) > 1 {
		for _, addr := range splited {
			peer := &Peer{
				address:    common.ParseAddress(addr),
				nextIndex:  0,
				matchIndex: 0,
			}
			peers = append(peers, peer)
		}
	} else {
		log.Fatalf("More than 1 host for Peers is required.")
	}

	return peers
}

//GetAddress : Get the address of peer node
func (p *Peer) GetAddress() *common.Address {
	return p.address
}

//GetNextIndex : Get nextIndex of peer node
func (p *Peer) GetNextIndex() int64 {
	return p.nextIndex
}

//SetNextIndex : Set nextIndex of peer node
func (p *Peer) SetNextIndex(ni int64) {
	p.nextIndex = ni
}

//GetMatchIndex : Get matchIndex of peer node
func (p *Peer) GetMatchIndex() int64 {
	return p.matchIndex
}

//SetMatchIndex : Set matchIndex of peer node
func (p *Peer) SetMatchIndex(mi int64) {
	p.matchIndex = mi
}

//DecreaseNextIndex : decrease nextIndex by 1
func (p *Peer) DecreaseNextIndex() {
	p.nextIndex--
}

//UpdateNextIndexAndMatchIndex : update nextIndex and matchIndex
func (p *Peer) UpdateNextIndexAndMatchIndex() {
	p.matchIndex = p.nextIndex
	p.nextIndex++
}
