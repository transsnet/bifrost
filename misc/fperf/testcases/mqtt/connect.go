package mqtt

import (
	"flag"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type connect struct {
	c   *client
	cli []MQTT.Client
}

func NewConnectCommand(cli *client, args []string) Command {
	c := &connect{}
	fs := flag.NewFlagSet("connect", flag.ExitOnError)
	setOpt(fs, &cli.opt)
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	c.c = cli
	return c
}

func (c *connect) Exec() error {
	cli := c.c
	cc, err := mqttConnect(cli.addr, cli.opt)
	if err != nil {
		return err
	}
	cc.Disconnect(0)
	return nil
}
