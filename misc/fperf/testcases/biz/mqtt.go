package link

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meitu/bifrost/grpc/publish"
	"google.golang.org/grpc"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fperf/fperf"
)

var (
	idgen      func() string
	expectMsg  []string
	expectLock sync.RWMutex
	equal      bool
)

func init() {
	fperf.Register("biz", NewMQTTClient, "biz fperf")
	idgen = idgenerator()
}

type client struct {
	cli publish.PublishServiceClient
	cmd Command
	opt mqttOpt
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
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	cli := publish.NewPublishServiceClient(conn)
	c.cli = cli
	return nil
}

func (c *client) Request() error {
	return c.cmd.Exec()
}

type mqttOpt struct {
	clientID string
	clean    bool
	topic    string
	qos      int
	addr     string

	same bool
}

func setOpt(fs *flag.FlagSet, opt *mqttOpt) {
	fs.StringVar(&opt.clientID, "clientid", "fperf-clientid", "ID of this client, this should be uniq")
	fs.StringVar(&opt.topic, "topic", "fperf-topic", "client subscribe topic")
	fs.StringVar(&opt.addr, "addr", "127.0.0.1:1883", "mqtt server")
	fs.BoolVar(&opt.clean, "cleansession", true, "Set cleansession flag")
	fs.BoolVar(&opt.same, "same", false, "all client same topic")
	fs.IntVar(&opt.qos, "qos", 1, "topic is ops")
}

func idgenerator() func() string {
	var i int32
	return func() string {
		count := atomic.AddInt32(&i, 1)
		return fmt.Sprintf("%d", count)
	}
}

func mqttConnect(opt mqttOpt) (string, error) {
	id := idgen()
	opts := MQTT.NewClientOptions().AddBroker(opt.addr)
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
		return "", token.Error()
	}

	var topic string
	if opt.same {
		topic = opt.topic
		if token := cli.Subscribe(topic, byte(opt.qos), handler); token.Wait() && token.Error() != nil {
			return "", token.Error()
		}
	} else {
		topic = opt.topic + "-" + id
		if token := cli.Subscribe(topic, byte(opt.qos), handler); token.Wait() && token.Error() != nil {
			return "", token.Error()
		}
	}
	return topic, nil
}

func handler(client MQTT.Client, msg MQTT.Message) {
	if equal {
		var count, err = strconv.Atoi(string(msg.Payload()))
		if err != nil {
			fmt.Errorf(" atoi %s", err)
		}

		expectLock.RLock()
		defer expectLock.Unlock()
		if !bytes.Equal(msg.Payload(), []byte(expectMsg[count])) {
			fmt.Errorf("expect msg %s,autual msg %s", expectMsg[count], string(msg.Payload()))
		}
	}
}
