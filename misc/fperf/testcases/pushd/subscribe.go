package pushd

import (
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type subscribe struct {
	c *client

	clientID     string
	address      string
	topics       []string
	qoss         []int32
	cleanSession bool
	count        int
}

func NewSubscribeCommand(c *client, args []string) Command {
	s := &subscribe{}
	fs := flag.NewFlagSet("subscribe", flag.ExitOnError)

	var topic string
	var qos int
	fs.StringVar(&s.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.StringVar(&s.address, "address", "127.0.0.1:2345", "Address of connd")
	fs.StringVar(&topic, "topic", "fperf-topic", "Topic to subscribe")
	fs.IntVar(&qos, "qos", 1, "Qos of message")
	fs.IntVar(&s.count, "count", 1, "The num of sending message")
	fs.BoolVar(&s.cleanSession, "cleansession", true, "Clean session true or false")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	s.c = c
	s.topics, s.qoss = s.buildmessage(topic, int32(qos))
	return s
}
func (s *subscribe) Exec() error {
	cli := s.c.cli
	ctx := s.c.ctx

	temp := idgen()
	req := &push.SubscribeReq{
		ClientID:     temp + "-" + s.clientID,
		CleanSession: s.cleanSession,
		Topics:       s.topics,
		TraceID:      temp,
		Qoss:         s.qoss,
	}

	_, err := cli.Subscribe(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *subscribe) buildmessage(topic string, qos int32) ([]string, []int32) {
	count := s.count
	if count == 0 {
		panic("count is less than one")
	}
	topics := make([]string, count)
	qoss := make([]int32, count)
	for i := 0; i < count; i++ {
		topics[i] = topic + "-" + idgen()
		qoss[i] = qos
	}
	return topics, qoss
}
