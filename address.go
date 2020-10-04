package main

import (
	"log"
	"strings"
)

type Address struct {
	name string
	port string
}

func parseAddress(addrString string) *Address {
	splited := strings.Split(addrString, ":")
	if 2 == len(splited) && splited[0] != "" && splited[1] != "" {
		return &Address{
			name: splited[0],
			port: splited[1],
		}
	} else {
		log.Fatalf("The address string is invalid.")
		return nil
	}
}

func parseAddresses(addrStrings string) []*Address {
	var addresses []*Address
	splited := strings.Split(addrStrings, ",")
	if len(splited) > 1 {
		for _, addr := range splited {
			addresses = append(addresses, parseAddress(addr))
		}
	} else {
		log.Fatalf("More than 1 host for Peers is required.")
	}

	return addresses
}
