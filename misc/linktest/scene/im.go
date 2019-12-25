package scene

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	pub "github.com/meitu/bifrost/grpc/publish"
)

type IM struct {
	Cli       *Client
	Stat      *Statistics
	Clis      *Clients
	KeepAlive time.Duration
	Interval  time.Duration
	Payload   string
	Exit      bool

	recvMsgs map[MQTT.Client]MQTT.Message
	mlock    sync.Mutex

	Pubc    *Pubcli
	PubAddr string
	wg      sync.WaitGroup
}

func NewIM() *IM {
	return &IM{
		recvMsgs: make(map[MQTT.Client]MQTT.Message),
		Clis:     NewClients(),
		Cli:      NewClient(),
		Stat:     NewStat(),
	}
}

func (t *IM) Start() {
	t.wg.Add(int(t.Stat.Count))
	for i := 0; i < t.Stat.Count; i++ {
		go func(id int) {
			defer t.wg.Done()
			Cli := t.Cli.MqttConnect(id, t.handler)
			t.Clis.Add(id, Cli)
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

func (t *IM) handler(Cli MQTT.Client, msg MQTT.Message) {
	t.mlock.Lock()
	if t.Exit {
		if _, ok := t.recvMsgs[Cli]; ok {
			panic(msg.Topic())
		}
	}
	t.recvMsgs[Cli] = msg
	t.mlock.Unlock()

	if t.Stat.Calculate() && t.Exit {
		t.Stat.Print(t.Exit)
		t.Detail()
		os.Exit(1)
	}
}

func (t *IM) Publish() {
	t.Stat.StartSendMsg = time.Now()
	if t.Payload == "" {
		t.Payload = time.Now().String()
	}
	tar := &pub.Target{}
	tar.Topic = t.Cli.Topic
	tar.Qos = int32(t.Cli.QoS)
	for {
		err := t.Pubc.Request([]byte(t.Payload), []*pub.Target{tar})
		if err == nil {
			break
		}
	}
}

func (t *IM) Tick() {
	if t.KeepAlive <= 0 {
		return
	}
	t.Publish()
	<-time.Tick(t.KeepAlive)
	t.Stat.Print(t.Exit)
	t.Detail()
	os.Exit(1)
}

func (t *IM) Detail() {
	t.mlock.Lock()
	defer t.mlock.Unlock()
	handler := func(id int, Cli MQTT.Client) {
		m, ok := t.recvMsgs[Cli]
		if !ok {
			log.Printf("the %d of Client discard message", id)
			return
		}
		if !strings.EqualFold(t.Payload, string(m.Payload())) {
			log.Printf("expect message %s,real message %s \n", t.Payload, m.Payload())
		}
	}
	t.Clis.Range(handler)
}
