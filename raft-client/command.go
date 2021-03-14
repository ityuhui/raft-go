package main

import "raft-go/common"

//Command : data structure of command
type Command struct {
	Mode common.CommandMode
	Text string
}

//ParseCommand : parse string to command
func ParseCommand(get string, set string) *Command {
	cmd := &Command{
		Mode: common.COMMANDMODE_UNKNOWN,
		Text: "",
	}
	if get != "" {
		cmd.Mode = common.COMMANDMODE_GET
		cmd.Text = get
	} else if set != "" {
		cmd.Mode = common.COMMANDMODE_SET
		cmd.Text = set
	}
	return cmd
}

//ToString : convert command mode to string
func (cmd *Command) ToString() string {
	return cmd.Mode.ToString() + " " + cmd.Text
}
