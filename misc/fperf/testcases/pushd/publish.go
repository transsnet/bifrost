package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type publish struct {
	c *client

	clientID string
	topic    string
	payload  string
	qos      int
	retain   bool
}

func NewPublishCommand(cli *client, args []string) Command {
	p := &publish{}
	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	fs.StringVar(&p.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.StringVar(&p.topic, "topic", "fperf-topic", "Topic to publish")
	fs.StringVar(&p.payload, "payload", "hello world", "Message content")
	fs.BoolVar(&p.retain, "retain", false, "retain message")
	fs.IntVar(&p.qos, "qos", 1, "Qos of message")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	p.c = cli
	return p
}

func (p *publish) Exec() error {
	cli := p.c.cli
	ctx := p.c.ctx
	temp := idgen()

	msg := &push.Message{
		Topic:   "1" + "-" + p.topic,
		Qos:     int32(p.qos),
		Payload: []byte(p.payload),
		TraceID: temp,
		Retain:  p.retain,
		BizID:   []byte(temp),
	}

	req := &push.PublishReq{
		ClientID: p.clientID + "-" + temp,
		Message:  msg,
	}

	_, err := cli.MQTTPublish(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
