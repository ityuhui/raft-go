package common

import (
	"log"
	"strings"
)

type Address struct {
	Name string
	Port string
}

func ParseAddress(addrString string) *Address {
	splited := strings.Split(addrString, ":")
	if 2 == len(splited) && splited[0] != "" && splited[1] != "" {
		return &Address{
			Name: splited[0],
			Port: splited[1],
		}
	} else {
		log.Fatalf("The address string is invalid.")
		return nil
	}
}

func (addr *Address) GenerateUName() string {
	return addr.Name + ":" + addr.Port
}
