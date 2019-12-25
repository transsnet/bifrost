package mqtt

import (
	"flag"

	"log"
)

type publish struct {
	cli *client

	topic   string
	qos     uint
	payload string
	retain  bool
}

func NewPublishCommand(cli *client, args []string) Command {
	c := &publish{}

	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	fs.StringVar(&c.topic, "topic", "/fperf/topic", "topic to publish")
	fs.StringVar(&c.payload, "payload", "hello world", "what you want to publish")
	fs.BoolVar(&c.retain, "retain", false, "retain messgae")
	fs.UintVar(&c.qos, "qos", 1, "qos should be 0, 1, 2")

	setOpt(fs, &cli.opt)

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}

	c.cli = cli
	return c
}

func (c *publish) Exec() error {
	qos := c.qos
	payload := c.payload
	topic := c.topic
	mq := c.cli.cli
	retain := c.retain

	if token := mq.Publish(topic, byte(qos), retain, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
