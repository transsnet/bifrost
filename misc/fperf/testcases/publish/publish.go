package pub

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/meitu/bifrost/grpc/publish"
)

type pub struct {
	c *client

	topic         string
	qos           int
	payload       string
	count         int
	bretain       bool
	noneDowngrade bool
	tars          []*publish.Target
}

func NewPublishCommand(cli *client, args []string) Command {
	p := &pub{}
	fs := flag.NewFlagSet("publish", flag.ExitOnError)

	fs.StringVar(&p.topic, "topic", "/fperf/topic", "Topic to publish")
	fs.IntVar(&p.qos, "qos", 1, "Qos of message")
	fs.StringVar(&p.payload, "payload", "hello world", "Content of message")
	fs.IntVar(&p.count, "count", 1, "Number of messages fetch one time")
	fs.BoolVar(&p.bretain, "retain", false, "retain if true or false")
	fs.BoolVar(&p.noneDowngrade, "nonedowngrade", true, "retain if true or false")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	p.c = cli
	p.tars = p.buildmessage()
	return p
}

func (p *pub) Exec() error {
	cli := p.c.cli
	ctx := p.c.ctx
	temp := time.Now().String()

	req := &publish.PublishRequest{
		MessageID: []byte(temp),
		Targets:   p.tars,
		Payload:   []byte(p.payload),
		Weight:    0,
	}
	_, err := cli.Publish(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (p *pub) buildmessage() []*publish.Target {
	count := p.count
	tars := make([]*publish.Target, count)
	for i := 0; i < count; i++ {
		tars[i] = &publish.Target{
			Qos:           int32(p.qos),
			Topic:         p.topic + "/" + fmt.Sprintf("%d", i),
			IsRetain:      p.bretain,
			NoneDowngrade: p.noneDowngrade,
		}
	}
	return tars
}
