package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type pull struct {
	c *client

	topic  string
	limit  int
	offset int
}

func NewPullCommand(cli *client, args []string) Command {
	p := &pull{}

	fs := flag.NewFlagSet("pull", flag.ExitOnError)
	fs.StringVar(&p.topic, "topic", "fperf-topic", "Topic to pull")
	fs.IntVar(&p.offset, "index", 0, "index message ID")
	fs.IntVar(&p.limit, "count", 100, "Number of messages fetch one time")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	p.c = cli
	return p
}

func (p *pull) Exec() error {
	cli := p.c.cli
	ctx := p.c.ctx
	req := &push.PullReq{
		Topic: p.topic + "-" + idgen(),
		// Offset: int64(p.offset),
		Limit: int64(p.limit),
	}

	cli.Pull(ctx, req)
	return nil
}
