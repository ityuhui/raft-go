package main

import (
	"context"
	"log"
	"raft-go/common"
	"raft-go/raft_rpc"
	"time"

	"google.golang.org/grpc"
)

type Client struct {
	headerAddr *common.Address
	command    *Command
}

var ins *Client = nil

func NewClientInstance(header string, get string, set string) *Client {

	ins = &Client{
		headerAddr: common.ParseAddress(header),
		command:    ParseCommand(get, set),
	}
	return ins
}

func (c *Client) Run() {
	c.executeCommand()
}

func (client *Client) executeCommand() (bool, int64) {
	log.Printf("Begin to execute the command: %v", client.command.ToString())
	// Set up a connection to the server.
	conn, err := grpc.Dial(client.headerAddr.Name+":"+client.headerAddr.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := raft_rpc.NewRaftServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ExecuteCommand(ctx, &raft_rpc.ExecuteCommandRequest{
		Mode: client.command.Mode.ToString(),
		Text: client.command.Text,
	})
	if err != nil {
		log.Fatalf("The raft deamon cannot execute the command from client: %v", err)
	}
	log.Printf("Execute command result: %v", r.GetSuccess())
	return r.GetSuccess(), r.GetValue()
}
