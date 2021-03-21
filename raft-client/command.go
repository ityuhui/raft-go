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
		Mode: common.CommandModeUnknown,
		Text: "",
	}
	if get != "" {
		cmd.Mode = common.CommandModeGet
		cmd.Text = get
	} else if set != "" {
		cmd.Mode = common.CommandModeSet
		cmd.Text = set
	}
	return cmd
}

//ToString : convert command mode to string
func (cmd *Command) ToString() string {
	return cmd.Mode.ToString() + " " + cmd.Text
}
