package main

type Command struct {
	Mode string
	Text string
}

func ParseCommand(set string, get string) *Command {
	cmd := &Command{
		Mode: "",
		Text: "",
	}
	return cmd
}

func (cmd *Command) ToString() string {
	return cmd.Mode + " " + cmd.Text
}
