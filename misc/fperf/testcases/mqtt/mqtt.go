package mqtt

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fperf/fperf"
)

var idgen func() string

func init() {
	fperf.Register("mqtt", NewMQTTClient, "mqtt client")
	idgen = idgenerator()
}

type client struct {
	cli   MQTT.Client
	cmd   Command
	opt   mqttOpt
	addr  string
	ready int32
}

func NewMQTTClient(flag *fperf.FlagSet) fperf.Client {
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
	return c
}

func (c *client) Dial(addr string) error {
	c.addr = addr
	go func() {
		var err error
		c.cli, err = mqttConnect(addr, c.opt)
		if err != nil {
			log.Fatal(err)
		}
		atomic.AddInt32(&c.ready, 1)
	}()
	return nil
}

func (c *client) Request() error {
	c.Ready()
	return c.cmd.Exec()
}

type mqttOpt struct {
	clientID string
	clean    bool
}

func setOpt(fs *flag.FlagSet, opt *mqttOpt) {
	fs.StringVar(&opt.clientID, "clientid", "fperf-clientid", "ID of this client, this should be uniq")
	fs.BoolVar(&opt.clean, "cleansession", true, "Set cleansession flag")
}

func idgenerator() func() string {
	var i int32
	return func() string {
		count := atomic.AddInt32(&i, 1)
		return fmt.Sprintf("%d", count)
	}
}

func mqttConnect(addr string, opt mqttOpt) (MQTT.Client, error) {
	id := idgen()
	opts := MQTT.NewClientOptions().AddBroker(addr)
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(opt.clientID + "-" + id)
	opts.SetUsername("username")
	opts.SetPassword("password")
	opts.SetCleanSession(opt.clean)
	opts.SetProtocolVersion(4)
	//防止bifrost处理速度过慢，导致连接丢失
	opts.SetConnectTimeout(30 * time.Minute)
	opts.SetKeepAlive(30 * time.Minute)
	opts.SetPingTimeout(30 * time.Minute)
	opts.SetWriteTimeout(30 * time.Minute)

	cli := MQTT.NewClient(opts)
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return cli, nil
}

func (c *client) Ready() {
	for {
		if atomic.LoadInt32(&c.ready) == 1 {
			return
		}
		time.Sleep(time.Second)
	}
}
