package main

import (
	"errors"
	"log"
	"sync"
)

//StateMachine : state machine on a node
type StateMachine struct {
	data map[string]int64
}

var stateMachineInstance *StateMachine = nil
var stateMachineLock sync.Mutex

//NewStateMachineInstance : Create a new instance of state machine
func NewStateMachineInstance() *StateMachine {

	stateMachineInstance = &StateMachine{
		data: make(map[string]int64),
	}
	return stateMachineInstance
}

//GetStateMachineInstance : Get the instance of state machine
func GetStateMachineInstance() *StateMachine {
	return stateMachineInstance
}

//Set : Add or update (key,value) in state machine
func (sm *StateMachine) Set(_key string, _newValue int64) {
	_curVal, ok := sm.data[_key]
	if ok {
		sm.data[_key] = _newValue
		log.Printf("The value for the key [%v] changes from [%v] to [%v].", _key, _curVal, _newValue)
	} else {
		sm.data[_key] = _newValue
		log.Printf("The value for the key [%v] is set to [%v].", _key, _newValue)
	}
}

//Get : Get value by key in state chine
func (sm *StateMachine) Get(_key string) (int64, error) {
	_curVal, ok := sm.data[_key]
	if ok {
		return _curVal, nil
	}
	errMsg := "The value for the key [" + _key + "] does not exist in state machine."
	log.Printf(errMsg)
	return 0, errors.New(errMsg)
}
