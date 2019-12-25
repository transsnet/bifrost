package connd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/conn"
)

type notify struct {
	c     *client
	topic string
	index int64
	none  bool
}

func NewNotifyCommand(cli *client, args []string) Command {
	c := &notify{}
	fs := flag.NewFlagSet("connect", flag.ExitOnError)
	fs.StringVar(&c.topic, "topic", "fperf-topic", "Topic that it is notified to operate")
	fs.Int64Var(&c.index, "index", 0, "The postion of message in topic ")
	fs.BoolVar(&c.none, "nonedowngrade", false, "msg is discard in topic")
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	c.c = cli
	return c
}
func (c *notify) Exec() error {

	cli := c.c.cli
	ctx := c.c.ctx
	topic := c.topic
	// index := c.index

	req := &conn.NotifyReq{
		Topic: "1" + "-" + topic,
		// Index:         index,
		NoneDowngrade: c.none,
	}

	_, err := cli.Notify(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
