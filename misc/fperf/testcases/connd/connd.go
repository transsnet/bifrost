package connd

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/fperf/fperf"
	"github.com/meitu/bifrost/grpc/conn"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var idgen func() string

func init() {
	fperf.Register("connd", NewClient, "connd grpc service")
	idgen = idgenerator()
}

type client struct {
	cli conn.ConnServiceClient
	ctx context.Context
	cmd Command
}

func NewClient(flag *fperf.FlagSet) fperf.Client {
	//subcommands: connect, subscribe, publish ...
	c := &client{}
	flag.Parse()
	if flag.NArg() < 1 {
		log.Println("subcommand invalid")
		fmt.Println("Avaliable subcommands list:")
		for name, _ := range SubCommands {
			fmt.Println("  ", name)
		}
		os.Exit(-1)
	}
	name := flag.Arg(0)
	cmdf, found := SubCommands[name]
	if !found {
		log.Fatalln("command not found:", name)
	}
	cmd := cmdf(c, flag.Args())
	c.cmd = cmd
	c.ctx = context.Background()
	return c
}

func (c *client) Dial(address string) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	cli := conn.NewConnServiceClient(cc)
	c.cli = cli
	return nil
}

func (c *client) Request() error {
	return c.cmd.Exec()
}

func idgenerator() func() string {
	var i int32
	return func() string {
		id := atomic.AddInt32(&i, 1)
		return fmt.Sprintf("%d", id)
	}
}
