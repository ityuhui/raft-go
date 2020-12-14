package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help bool

	I     string
	peers string
)

func init() {
	flag.BoolVar(&help, "help", false, "Show the usage")
	flag.StringVar(&I, "I", "", "My name:port")
	flag.StringVar(&peers, "Peers", "", "The name:port of my peers")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage:
raft-daemon --I=myip:port --Peers=host2:port,host3:port
`)
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if help || I == "" || peers == "" {
		flag.Usage()
		return
	}

	n := NewNodeInstance(I, peers)
	n.Run()
}
