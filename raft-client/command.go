package main

type Command struct {
	Mode CommandMode
	Text string
}

func ParseCommand(set string, get string) *Command {
	cmd := &Command{
		Mode: COMMANDMODE_UNKNOWN,
		Text: "",
	}
	if get != "" {
		cmd.Mode = COMMANDMODE_GET
		cmd.Text = get
	} else if set != "" {
		cmd.Mode = COMMANDMODE_SET
		cmd.Text = set
	}
	return cmd
}

func (cmd *Command) ToString() string {
	return cmd.Mode.ToString() + " " + cmd.Text
}
