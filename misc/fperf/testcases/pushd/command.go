package pushd

type Command interface {
	Exec() error
}

var SubCommands map[string]func(*client, []string) Command

func init() {
	SubCommands = make(map[string]func(*client, []string) Command)
	SubCommands["connect"] = NewConnectCommand
	SubCommands["disconnect"] = NewDisconnectCommand
	SubCommands["publish"] = NewPublishCommand
	SubCommands["subscribe"] = NewSubscribeCommand
	SubCommands["unsubscribe"] = NewUnsubscribeCommand
	SubCommands["delunack"] = NewDelUnackCommand
	SubCommands["rangeunack"] = NewRangeUnackCommand
	SubCommands["putunack"] = NewPutUnackCommand
	SubCommands["postsubscribe"] = NewPostsubscribeCommand
	SubCommands["pull"] = NewPullCommand
	SubCommands["addroute"] = NewAddrouteCommand
	SubCommands["removeroute"] = NewRemoverouteCommand
	//TODO Pubrec Pubrel Pubcomp
}
