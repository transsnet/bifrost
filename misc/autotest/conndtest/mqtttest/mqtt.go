package mqtttest

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	password = "password"
	username = "bifrost-appkey=service-1544516816-1-bb7c77b22be1ff945d3272"
)

type MQTTClient struct {
	cli          MQTT.Client
	username     string
	password     string
	cleansession bool
	cid          string

	qos     int
	topic   string
	msgchan chan MQTT.Message
}

func NewMQTTClient(cid, topic string, qos int, addr string, cs bool) (*MQTTClient, error) {
	mqttc := &MQTTClient{
		cid:          cid,
		topic:        topic,
		qos:          qos,
		msgchan:      make(chan MQTT.Message, 10),
		username:     username,
		password:     password,
		cleansession: cs,
	}
	if err := mqttc.Connect(addr); err != nil {
		return nil, err
	}
	return mqttc, nil
}

func (t *MQTTClient) handler(c MQTT.Client, m MQTT.Message) {
	t.msgchan <- m
}

func (t *MQTTClient) Connect(addr string) error {
	opts := MQTT.NewClientOptions().AddBroker(addr)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetAutoReconnect(false)
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(t.cid)
	opts.SetUsername(t.username)
	opts.SetPassword(t.password)
	opts.SetCleanSession(t.cleansession)
	opts.SetProtocolVersion(4)
	opts.SetKeepAlive(time.Minute * 10)
	opts.SetWriteTimeout(time.Minute * 10)
	opts.SetPingTimeout(time.Minute * 10)
	opts.SetConnectTimeout(time.Minute * 10)

	t.cli = MQTT.NewClient(opts)

	if token := t.cli.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Connect failed, %s", token.Error())
		return token.Error()
	}
	if token := t.cli.Subscribe(t.topic, byte(t.qos), t.handler); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		log.Printf("Subscribe failed, %s", token.Error())
		return token.Error()
	}
	t.cli.AddRoute(t.topic, t.handler)
	return nil
}

func (t *MQTTClient) UnSubscribe() error {
	if token := t.cli.Unsubscribe(t.topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *MQTTClient) Disconnect() {
	t.cli.Disconnect(1)
}

func (t *MQTTClient) Recevice() MQTT.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-t.msgchan:
			return msg
		}
	}
	return nil
}

func (t *MQTTClient) MsgLen() int {
	return len(t.msgchan)
}

func SetPassword(password string) {
	password = password
}

func SetUsername(username string) {
	username = username
}
