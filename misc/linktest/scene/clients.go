package scene

import (
	"math/rand"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Clients struct {
	clients     map[int]MQTT.Client
	clock       sync.Mutex
	DropCount   int
	MqttConnect func(id int, handler MQTT.MessageHandler) MQTT.Client
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[int]MQTT.Client),
	}
}

func (c *Clients) SetHandler(handler func(id int, handler MQTT.MessageHandler) MQTT.Client) {
	c.MqttConnect = handler
}

func (c *Clients) Add(id int, cli MQTT.Client) {
	c.clock.Lock()
	c.clients[id] = cli
	c.clock.Unlock()
}

func (c *Clients) Delete(id int) {
	c.clock.Lock()
	delete(c.clients, id)
	c.clock.Unlock()
}

func (c *Clients) Get(id int) MQTT.Client {
	c.clock.Lock()
	cli := c.clients[id]
	c.clock.Unlock()
	return cli
}

func (c *Clients) RandromDropClient(handler MQTT.MessageHandler) {
	for c.DropCount != 0 {
		time.Sleep(time.Second)
		c.clock.Lock()
		for i := 0; i < c.DropCount; i++ {
			id := rand.Intn(len(c.clients))
			if cli, ok := c.clients[id]; ok {
				cli.Disconnect(1)
				c.clients[id] = c.MqttConnect(id, handler)
				continue
			}

		}
		c.clock.Unlock()
	}
}

func (c *Clients) Length() int {
	c.clock.Lock()
	defer c.clock.Unlock()
	return len(c.clients)
}

func (c *Clients) Vaild() bool {
	if c.DropCount == 0 {
		return false
	}
	return true
}

func (c *Clients) Range(handler func(id int, cli MQTT.Client)) {
	if c.Vaild() {
		return
	}
	for k, v := range c.clients {
		handler(k, v)
	}
}
