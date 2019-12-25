package pushd

import (
	"flag"
	"log"
	"time"

	"github.com/meitu/bifrost/grpc/push"
)

type unsubscribe struct {
	c *client

	clientID     string
	topics       []string
	cleanSession bool
	count        int
}

func NewUnsubscribeCommand(cli *client, args []string) Command {
	u := &unsubscribe{}
	var topic string
	fs := flag.NewFlagSet("unsubscribe", flag.ExitOnError)
	fs.StringVar(&u.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.StringVar(&topic, "topic", "fperf-topic", "Topic to unsubscribe")
	fs.BoolVar(&u.cleanSession, "cleansession", true, "Clean session true or false")
	fs.IntVar(&u.count, "count", 1, "The num of sending message")
	u.topics = u.buildmessage(topic)
	u.c = cli
	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	return u
}

func (u *unsubscribe) Exec() error {
	cli := u.c.cli
	ctx := u.c.ctx

	req := &push.UnsubscribeReq{
		ClientID:     idgen() + "-" + u.clientID,
		Topics:       u.topics,
		CleanSession: u.cleanSession,
		TraceID:      time.Now().String(),
	}

	_, err := cli.Unsubscribe(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *unsubscribe) buildmessage(topic string) []string {

	count := s.count
	topics := make([]string, count)

	for i := 0; i < count; i++ {
		topics[i] = topic + "-" + idgen()
	}
	return topics
}
