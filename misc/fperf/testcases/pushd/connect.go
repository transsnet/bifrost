package pushd

import (
	"flag"
	"log"
	"strconv"

	"github.com/meitu/bifrost/grpc/push"
)

type connect struct {
	c *client

	clientID     string
	conndAddr    string
	cleanSession bool
}

func NewConnectCommand(cli *client, args []string) Command {
	c := &connect{}
	fs := flag.NewFlagSet("connect", flag.ExitOnError)
	fs.StringVar(&c.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.StringVar(&c.conndAddr, "conndaddr", "127.0.0.1:2345", "conndaddr")
	fs.BoolVar(&c.cleanSession, "cleansession", true, "Clean session true or false")
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	c.c = cli
	return c
}

func (c *connect) Exec() error {
	cli := c.c.cli
	ctx := c.c.ctx
	temp := idgen()
	id, _ := strconv.Atoi(temp)

	req := &push.ConnectReq{
		ClientID:      c.clientID + "-" + temp,
		Username:      "username",
		Password:      []byte("password"),
		CleanSession:  c.cleanSession,
		GrpcAddress:   c.conndAddr,
		ClientAddress: temp,
		TraceID:       "mid",
		ConnectionID:  int64(id),
	}
	_, err := cli.Connect(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
