package mqtt

type Command interface {
	Exec() error
}

var SubCommands map[string]func(*client, []string) Command

func init() {
	SubCommands = make(map[string]func(*client, []string) Command)
	SubCommands["connect"] = NewConnectCommand
	SubCommands["publish"] = NewPublishCommand
	SubCommands["subscribe"] = NewSubscribeCommand
}
