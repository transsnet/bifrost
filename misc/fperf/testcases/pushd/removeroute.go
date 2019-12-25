package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type removeroute struct {
	c         *client
	conndAddr string
	topic     string
}

func NewRemoverouteCommand(c *client, args []string) Command {
	d := &removeroute{}
	fs := flag.NewFlagSet("removeroute", flag.ExitOnError)
	fs.StringVar(&d.conndAddr, "conndaddr", "127.0.0.1:2345", "conndaddr")
	fs.StringVar(&d.topic, "topic", "fperf-topic", "Topic to subscribe")
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	d.c = c
	return d
}

func (d *removeroute) Exec() error {
	cli := d.c.cli
	ctx := d.c.ctx
	req := &push.RemoveRouteReq{
		Topic:       d.topic + "-" + idgen(),
		GrpcAddress: d.conndAddr,
	}

	if _, err := cli.RemoveRoute(ctx, req); err != nil {
		return err
	}
	return nil
}
