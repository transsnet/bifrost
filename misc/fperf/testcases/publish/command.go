package pub

type Command interface {
	Exec() error
}

var SubCommands map[string]func(*client, []string) Command

func init() {
	SubCommands = make(map[string]func(*client, []string) Command)
	SubCommands["publish"] = NewPublishCommand
}
