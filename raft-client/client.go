package main

import (
	"raft-go/common"
	"fmt"
	"os"
)

type Client struct {
	headerAddr *Address
	set string
	get string
}

var ins *Client = nil

func NewClientInstance(header string, set string, get string) *Client {

	ins = &Client{
		headerAddr:          ParseAddress(header),
		set : set,
		get : get,
	}
	return ins
}

func (c *Client) Run() {
	if c.get != "" {
		c.Get(get)
	} else if c.set != "" {
		c.Set(set)
	} else {
		log.Printf("No command to execute.")
	}
}

func (c *Client) Set() bool {
	log.Printf("Begin to send vote request to: %v", addr.GenerateUName())
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.name+":"+addr.port, grpc.WithInsecure(), grpc.WithBlock())
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
