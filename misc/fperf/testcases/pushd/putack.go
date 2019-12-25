package pushd

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/meitu/bifrost/grpc/push"
)

type putunack struct {
	c *client

	clientID string

	msg []*push.UnackDesc
}

func NewPutUnackCommand(cli *client, args []string) Command {
	c := &putunack{}
	var topic string
	var count int
	fs := flag.NewFlagSet("putack", flag.ExitOnError)
	fs.StringVar(&c.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.StringVar(&topic, "topic", "fperf-topic", "Topic to subscribe")
	fs.IntVar(&count, "count", 1, "The num of sending message")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	c.c = cli
	c.msg = c.buildmessage(topic, count)
	return c
}

func (c *putunack) Exec() error {
	cli := c.c.cli
	ctx := c.c.ctx
	temp := idgen()
	req := &push.PutUnackReq{
		ClientID: c.clientID + "-" + temp,
	}
	_, err := cli.PutUnack(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *putunack) buildmessage(topic string, count int) []*push.UnackDesc {
	msgs := make([]*push.UnackDesc, count)
	for i, _ := range msgs {
		msgs[i] = &push.UnackDesc{
			Topic: "1" + "-" + topic,
			// Index:     0,
			MessageID: int64(time.Now().Nanosecond()),
			BizID:     []byte(fmt.Sprintf("%d", i)),
		}
	}
	return msgs
}
