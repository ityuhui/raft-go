package main

import "raft-go/common"

type Command struct {
	Mode common.CommandMode
	Text string
}

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

func (cmd *Command) ToString() string {
	return cmd.Mode.ToString() + " " + cmd.Text
}
