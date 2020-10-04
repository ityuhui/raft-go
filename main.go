package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help bool

	Iam   string
	peers string
)

func init() {
	flag.BoolVar(&help, "help", false, "Show the usage")
	flag.StringVar(&Iam, "Iam", "", "My name:port")
	flag.StringVar(&peers, "Peers", "", "The name:port of my peers")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage:
raft-go --Iam=myip:port --Peers=host2:port,host3:port
`)
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if help || Iam == "" || peers == "" {
		flag.Usage()
		return
	}

	n := NewNodeInstance(Iam, peers)
	n.Run()
}
