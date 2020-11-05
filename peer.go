package main

import (
	"log"
	"strings"
)

type Peer struct {
	address    *Address
	nextIndex  int64
	matchIndex int64
}

func InitPeers(addrStrings string) []*Peer {
	var peers []*Peer
	splited := strings.Split(addrStrings, ",")
	if len(splited) > 1 {
		for _, addr := range splited {
			peer := &Peer{
				address:    ParseAddress(addr),
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

func (p *Peer) GetAddress() *Address {
	return p.address
}

func (p *Peer) GetNextIndex() int64 {
	return p.nextIndex
}

func (p *Peer) GetMatchIndex() int64 {
	return p.matchIndex
}
