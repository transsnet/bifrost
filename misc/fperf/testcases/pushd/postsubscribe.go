package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type postsubscribe struct {
	c *client

	clientID string
}

func NewPostsubscribeCommand(c *client, args []string) Command {
	d := &postsubscribe{}
	fs := flag.NewFlagSet("postsubscribe", flag.ExitOnError)

	fs.StringVar(&d.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	d.c = c
	return d
}

func (d *postsubscribe) Exec() error {
	cli := d.c.cli
	ctx := d.c.ctx

	temp := idgen()

	req := &push.PostSubscribeReq{
		ClientID: d.clientID + "-" + temp,
	}

	if _, err := cli.PostSubscribe(ctx, req); err != nil {
		return err
	}
	return nil
}
