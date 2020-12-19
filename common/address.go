package common

import (
	"log"
	"strings"
)

type Address struct {
	name string
	port string
}

func ParseAddress(addrString string) *Address {
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

func (addr *Address) GenerateUName() string {
	return addr.name + ":" + addr.port
}
