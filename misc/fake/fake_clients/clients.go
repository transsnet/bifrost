package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	clients map[MQTT.Client]string
	total   int64
	clock   sync.Mutex
	idc     int
	count   int
)

type thord struct {
	address   string
	topic     string
	username  string
	password  string
	clientid  string
	qos       int
	keepAlive time.Duration
	id        func() int
	wg        sync.WaitGroup
}

func handler() MQTT.MessageHandler {
	handler := func(c MQTT.Client, m MQTT.Message) {
		atomic.AddInt64(&total, 1)
	}
	return handler
}

func idgen() func() int {
	i := int32(idc)
	return func() int {
		return int(atomic.AddInt32(&i, 1))
	}
}

func (t *thord) mqttConnect() {

	id := t.id()
	idstr := fmt.Sprintf("%d", id)

	opts := MQTT.NewClientOptions().AddBroker(t.address)
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetAutoReconnect(true)
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(t.clientid + "-" + idstr)
	opts.SetUsername(t.username)
	opts.SetPassword(t.password)
	opts.SetCleanSession(true)
	opts.SetProtocolVersion(4)
	opts.SetKeepAlive(time.Minute * 10)
	opts.SetWriteTimeout(time.Minute * 10)
	opts.SetPingTimeout(time.Minute * 10)
	opts.SetConnectTimeout(time.Minute * 10)

	opts.SetOnConnectHandler(func(client MQTT.Client) { // 设置OnConnect回调函数，建立连接或者自动重连时触发。
		if token := client.Subscribe(t.topic, byte(t.qos), handler()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
			log.Printf("Subscribe failed, %s", token.Error())
		}
	})

	mc := MQTT.NewClient(opts)
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	// if token := mc.Subscribe(t.topic, byte(t.qos), handler()); token.Wait() && token.Error() != nil {
	// log.Fatalln(token.Error())
	// }

	clock.Lock()
	clients[mc] = idstr
	clock.Unlock()

	t.wg.Done()
}

func main() {
	clients = make(map[MQTT.Client]string, 60000)
	var t thord
	flag.StringVar(&t.address, "server", "tcp://127.0.0.1:8000", "address of mqtt broker")
	flag.IntVar(&t.qos, "qos", 1, "publish of qos which is range of 0-2")
	flag.IntVar(&count, "count", 1, "the num of connectiong ")
	flag.StringVar(&t.username, "username", "test", "username")
	flag.StringVar(&t.password, "password", "test", "password")
	flag.StringVar(&t.clientid, "clientid", "fperf-clientid", "prefix of clientid")
	flag.StringVar(&t.topic, "topic", "/fperf/topic", "topic to subscribe")
	flag.DurationVar(&t.keepAlive, "keep", 0, "service of keeplive time")
	flag.IntVar(&idc, "idc", 0, "id begin")
	flag.Parse()
	t.id = idgen()
	log.Printf("%#v %#v %#v \n", t, count, idc)

	t.wg.Add(count)
	for i := 0; i < count; i++ {
		go t.mqttConnect()
		time.Sleep(time.Millisecond)
	}
	t.wg.Wait()
	for {
		time.Sleep(10 * time.Second)
		fmt.Println(total)
	}
}
