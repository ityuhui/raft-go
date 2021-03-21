package common

type CommandMode int32

const (
	//CommandModeUnknown : unknown command
	CommandModeUnknown CommandMode = 0
	//CommandModeGet : "get" command
	CommandModeGet CommandMode = 1
	//CommandModeSet : "set" command
	CommandModeSet CommandMode = 2
)

func (cm CommandMode) ToString() string {
	switch cm {
	case CommandModeGet:
		return "GET"
	case CommandModeSet:
		return "SET"
	default:
		return "UNKNOWN"
	}
}
