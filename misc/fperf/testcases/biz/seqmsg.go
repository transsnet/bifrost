package link

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/meitu/bifrost/grpc/publish"
)

type seq struct {
	c *client

	tars          []*publish.Target
	count         int
	noneDowngrade bool
	bretain       bool
	idgen         func() string
}

func NewSeqCommand(cli *client, args []string) Command {
	p := &seq{}
	fs := flag.NewFlagSet("seq", flag.ExitOnError)
	fs.BoolVar(&p.noneDowngrade, "nonedowngrade", true, "nonedowngrade if true or false")
	fs.BoolVar(&p.bretain, "bretain", false, "retain if true or false")
	fs.BoolVar(&equal, "equal", true, "the message order is charming ")
	fs.IntVar(&p.count, "count", 1, "the num of client")
	setOpt(fs, &cli.opt)

	if err := fs.Parse(args[1:]); err != nil {
		log.Fatal(err)
	}

	p.c = cli
	p.idgen = idgenerator()
	p.initClient()
	return p
}

func (p *seq) Exec() error {
	cli := p.c.cli
	ctx := context.TODO()
	tars := p.tars
	msg := p.idgen()
	expectLock.Lock()
	expectMsg = append(expectMsg, msg)
	expectLock.Unlock()

	req := &publish.PublishRequest{
		MessageID: []byte(msg),
		Targets:   tars,
		Payload:   []byte(msg),
	}
	_, err := cli.Publish(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (p *seq) initClient() {
	count := p.count
	if p.c.opt.same {
		tars := make([]*publish.Target, 1)
		for i := 0; i < count; i++ {
			_, err := mqttConnect(p.c.opt)
			if err != nil {
				panic(err.Error())
			}
		}
		tars[0] = &publish.Target{
			Qos:           int32(1),
			Topic:         p.c.opt.topic,
			IsRetain:      false,
			NoneDowngrade: p.noneDowngrade,
		}
		p.tars = tars
	} else {
		tars := make([]*publish.Target, count)
		for i := 0; i < count; i++ {
			topic, err := mqttConnect(p.c.opt)
			if err != nil {
				fmt.Errorf("init client failed")
			}
			tars[i] = &publish.Target{
				Qos:           int32(1),
				Topic:         topic,
				IsRetain:      false,
				NoneDowngrade: p.noneDowngrade,
			}
		}
		p.tars = tars
	}
}
