package pushd

import (
	"flag"
	"log"
	"strconv"

	"github.com/meitu/bifrost/grpc/push"
)

type delunack struct {
	c *client

	clientID     string
	cleansession bool
}

func NewDelUnackCommand(cli *client, args []string) Command {
	p := &delunack{}
	fs := flag.NewFlagSet("delunack", flag.ExitOnError)
	fs.StringVar(&p.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.BoolVar(&p.cleansession, "cleansession", true, "Clean session true or false")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	p.c = cli
	return p
}

func (p *delunack) Exec() error {
	cli := p.c.cli
	ctx := p.c.ctx

	temp := idgen()
	mid, _ := strconv.ParseInt(temp, 10, 0)
	req := &push.DelUnackReq{
		ClientID:     p.clientID + "-" + temp,
		MessageID:    mid,
		CleanSession: p.cleansession,
		TraceID:      "mid" + "-" + temp,
	}
	_, err := cli.DelUnack(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
