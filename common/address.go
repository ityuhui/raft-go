package common

import (
	"log"
	"strings"
)

//Address : the address of a node
type Address struct {
	Name string
	Port string
}

//ParseAddress : parase the string to address data structure
func ParseAddress(addrString string) *Address {
	splited := strings.Split(addrString, ":")
	if 2 == len(splited) && splited[0] != "" && splited[1] != "" {
		return &Address{
			Name: splited[0],
			Port: splited[1],
		}
	}
	log.Fatalf("The address string is invalid.")
	return nil
}

//GenerateUName : generate an unique name for a node
func (addr *Address) GenerateUName() string {
	return addr.Name + ":" + addr.Port
}
