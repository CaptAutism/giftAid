package cmd

var commands map[string]Command

type Command interface {
	Name() string
	Exec(args []string)
}

func Register(name string, cmd Command) {
	commands[name] = cmd
}

func Run(args []string) {
	if len(args) < 2 {
		return
	}
	if cmd, ok := commands[args[1]]; ok {
		cmd.Exec(args[2:])
	}
}

type StartCmd struct{}

func (cmd *StartCmd) Name() string { return "start" }
func (cmd *StartCmd) Exec(args []string) {

}
