package main

import "sync"

type Client struct {
}

var ins *Client = nil
var lock sync.Mutex

func NewClientInstance() *Client {

	ins = &Client{}
	return ins
}

func (c *Client) Get(cmd string) {

}

func (c *Client) Set(cmd string) {

}
