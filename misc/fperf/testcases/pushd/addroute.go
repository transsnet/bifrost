package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type addroute struct {
	c         *client
	conndAddr string
	topic     string
}

func NewAddrouteCommand(c *client, args []string) Command {
	d := &addroute{}
	fs := flag.NewFlagSet("addroute", flag.ExitOnError)
	fs.StringVar(&d.conndAddr, "conndaddr", "127.0.0.1:2345", "conndaddr")
	fs.StringVar(&d.topic, "topic", "fperf-topic", "Topic to subscribe")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	d.c = c
	return d
}

func (d *addroute) Exec() error {
	cli := d.c.cli
	ctx := d.c.ctx
	req := &push.AddRouteReq{
		Topic:       d.topic + "-" + idgen(),
		GrpcAddress: d.conndAddr,
	}

	if _, err := cli.AddRoute(ctx, req); err != nil {
		return err
	}
	return nil
}
