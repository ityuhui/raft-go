package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help bool

	header string
	set    string
	get    string
)

func init() {
	flag.BoolVar(&help, "help", false, "Show the usage")
	flag.StringVar(&header, "header", "", "header name:port")
	flag.StringVar(&set, "set", "", "key=value")
	flag.StringVar(&get, "get", "", "key")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage:
raft-client --header=ip:port --set key=value
raft-client --header=ip:port --get key
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if help || header == "" || (get == "" && set == "") || (get != "" && set != "") {
		flag.Usage()
		return
	}

	c := NewClientInstance(header, get, set)
	c.Run()
}
