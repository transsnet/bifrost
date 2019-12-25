package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type disconnect struct {
	c *client

	clientID     string
	address      string
	cleanSession bool
}

func NewDisconnectCommand(c *client, args []string) Command {
	d := &disconnect{}
	fs := flag.NewFlagSet("disconnect", flag.ExitOnError)

	fs.StringVar(&d.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.StringVar(&d.address, "address", "127.0.0.1:2345", "Address of connd")
	fs.BoolVar(&d.cleanSession, "cleansession", false, "Clean session true or false")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	d.c = c
	return d
}

func (d *disconnect) Exec() error {
	cli := d.c.cli
	ctx := d.c.ctx

	temp := idgen()

	req := &push.DisconnectReq{
		ClientID:     d.clientID + "-" + temp,
		GrpcAddress:  d.address,
		CleanSession: d.cleanSession,
		Lost:         false,
		TraceID:      "mid" + temp,
	}

	if _, err := cli.Disconnect(ctx, req); err != nil {
		return err
	}
	return nil
}
