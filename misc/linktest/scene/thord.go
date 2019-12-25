package scene

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	pub "github.com/meitu/bifrost/grpc/publish"
)

type Thord struct {
	Cli       *Client
	Stat      *Statistics
	Clis      *Clients
	KeepAlive time.Duration
	Interval  time.Duration
	Payload   string
	Exit      bool

	recvMsgs map[string]MQTT.Message
	mlock    sync.Mutex

	Pubc     *Pubcli
	PubAddr  string
	sendMsgs map[string]string
	slock    sync.Mutex
	wg       sync.WaitGroup
}

func NewThord() *Thord {
	return &Thord{
		recvMsgs: make(map[string]MQTT.Message),
		sendMsgs: make(map[string]string),
		Clis:     NewClients(),
		Cli:      NewClient(),
		Stat:     NewStat(),
	}
}

func (t *Thord) Start() {
	t.wg.Add(int(t.Stat.Count))
	for i := 0; i < t.Stat.Count; i++ {
		go func(id int) {
			defer t.wg.Done()
			cli := t.Cli.MqttConnect(id, t.handler)
			t.Clis.Add(id, cli)
		}(i)
	}
	t.wg.Wait()

	go t.Tick()
	go t.Clis.RandromDropClient(t.handler)

	for {
		if t.KeepAlive == 0 {
			t.Publish()
			if t.Exit {
				t.KeepAlive = -1
			}
		}
		time.Sleep(t.Interval)
		t.Stat.Print(!t.Exit)
	}
}

func (t *Thord) handler(cli MQTT.Client, msg MQTT.Message) {
	t.mlock.Lock()
	if t.Exit {
		if _, ok := t.recvMsgs[msg.Topic()]; ok {
			panic(msg.Topic())
		}
	}
	t.recvMsgs[msg.Topic()] = msg
	t.mlock.Unlock()

	if t.Stat.Calculate() && t.Exit {
		t.Stat.Print(t.Exit)
		t.Detail()
		os.Exit(1)
	}
}

func (t *Thord) Publish() {
	t.Stat.StartSendMsg = time.Now()
	var wg sync.WaitGroup
	wg.Add(t.Stat.Count)
	for i := 0; i < t.Stat.Count; i++ {
		go func(i int) {
			defer wg.Done()
			idstr := fmt.Sprintf("%d", i)
			load := t.Payload
			if load == "" {
				load = time.Now().String()
				load = load + "-" + idstr
			}
			tar := &pub.Target{}
			tar.Topic = t.Cli.Topic + "-" + idstr
			tar.Qos = int32(t.Cli.QoS)
			for {
				err := t.Pubc.Request([]byte(load), []*pub.Target{tar})
				if err == nil {
					break
				}
			}
			t.slock.Lock()
			t.sendMsgs[tar.Topic] = load
			t.slock.Unlock()
		}(i)
	}
	wg.Wait()
}

func (t *Thord) Tick() {
	if t.KeepAlive <= 0 {
		return
	}
	t.Publish()
	<-time.Tick(t.KeepAlive)
	t.Stat.Print(t.Exit)
	t.Detail()
	os.Exit(1)
}

func (t *Thord) Detail() {
	if t.Clis.Vaild() {
		return
	}
	t.mlock.Lock()
	defer t.mlock.Unlock()
	t.slock.Lock()
	defer t.slock.Unlock()
	for topic, msg := range t.sendMsgs {
		m, ok := t.recvMsgs[topic]
		if !ok {
			log.Printf("the %s of client discard message", topic)
			continue
		}
		if !strings.EqualFold(msg, string(m.Payload())) {
			log.Printf("expect message %s,real message %s \n", msg, m.Payload())
			continue
		}
	}
}
