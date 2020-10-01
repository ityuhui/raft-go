package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help bool

	host string
	port string
)

func init() {
	flag.BoolVar(&help, "help", false, "Show the usage")
	flag.StringVar(&host, "host", "", "host name or host ip")
	flag.StringVar(&port, "port", "", "server port")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage:
raft-go --host= --port=
`)
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if help || host == "" || port == "" {
		flag.Usage()
		return
	}

	n := NewNode(host, port)
	n.Run()
}
