package main

type CommandMode int32

const (
	COMMANDMODE_UNKNOWN CommandMode = 0
	COMMANDMODE_GET     CommandMode = 1
	COMMANDMODE_SET     CommandMode = 2
)

func (cm CommandMode) ToString() string {
	switch cm {
	case COMMANDMODE_GET:
		return "GET"
	case COMMANDMODE_SET:
		return "SET"
	default:
		return "UNKNOWN"
	}
}
