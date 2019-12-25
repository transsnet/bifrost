package connd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/conn"
)

type disconnect struct {
	c        *client
	clientID string
}

func NewDisconnectCommand(c *client, args []string) Command {
	d := &disconnect{}
	fs := flag.NewFlagSet("disconnect", flag.ExitOnError)
	fs.StringVar(&d.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	d.c = c
	return d
}

func (d *disconnect) Exec() error {
	cli := d.c.cli
	ctx := d.c.ctx

	req := &conn.DisconnectReq{
		ClientID: idgen() + "-" + d.clientID,
	}

	if _, err := cli.Disconnect(ctx, req); err != nil {
		return err
	}
	return nil
}
