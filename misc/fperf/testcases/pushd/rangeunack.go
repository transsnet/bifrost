package pushd

import (
	"encoding/binary"
	"flag"
	"log"

	"github.com/meitu/bifrost/grpc/push"
)

type rangeunack struct {
	c *client

	clientID string
	limit    int
	offset   int
}

func NewRangeUnackCommand(c *client, args []string) Command {
	s := &rangeunack{}
	fs := flag.NewFlagSet("rangeunack", flag.ExitOnError)

	fs.StringVar(&s.clientID, "clientid", "fperf-clientid", "Client ID of MQTT Client")
	fs.IntVar(&s.offset, "index", 0, "index message ID")
	fs.IntVar(&s.limit, "count", 100, "Number of messages fetch one time")

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}
	s.c = c
	return s
}
func (s *rangeunack) Exec() error {
	cli := s.c.cli
	ctx := s.c.ctx
	offsetBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(offsetBytes, uint64(s.offset))

	temp := idgen()
	req := &push.RangeUnackReq{
		ClientID: s.clientID + "-" + temp,
		Offset:   offsetBytes,
		Limit:    int64(s.limit),
	}

	_, err := cli.RangeUnack(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
