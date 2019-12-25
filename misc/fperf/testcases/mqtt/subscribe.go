package mqtt

import (
	"flag"
	"log"
)

type subscribe struct {
	cli *client

	topic string
	qos   uint
	idgen func() string
	unsub bool
}

func NewSubscribeCommand(cli *client, args []string) Command {
	c := &subscribe{}
	fs := flag.NewFlagSet("subscribe", flag.ExitOnError)
	setOpt(fs, &cli.opt)
	fs.StringVar(&c.topic, "topic", "/fperf/topic", "Topic to subscribe")
	fs.UintVar(&c.qos, "qos", 1, "Qos should be 0, 1, 2")
	fs.BoolVar(&c.unsub, "unsub", false, "Unsub after subscribe")
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	c.cli = cli
	c.idgen = idgenerator()
	return c
}

func (c *subscribe) Exec() error {
	qos := c.qos
	topic := c.topic
	cli := c.cli

	mq := cli.cli
	idstr := c.idgen()
	if token := mq.Subscribe(topic+"/"+idstr, byte(qos), nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if c.unsub {
		if token := mq.Unsubscribe(topic + "/" + idstr); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}
