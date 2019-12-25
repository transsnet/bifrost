package mqtttest

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type WillClient struct {
	cli      MQTT.Client
	username string
	password string
	cid      string

	qos     int
	topic   string
	retain  bool
	payload string
	msgchan chan MQTT.Message
}

func NewWillClient(cid, topic string, qos int, addr string, retain bool, payload string) (*WillClient, error) {
	mqttc := &WillClient{
		cid:      cid,
		topic:    topic,
		qos:      qos,
		msgchan:  make(chan MQTT.Message, 10),
		username: username,
		password: password,
		retain:   retain,
		payload:  payload,
	}
	if err := mqttc.Connect(addr); err != nil {
		return nil, err
	}
	return mqttc, nil
}

func (t *WillClient) handler(c MQTT.Client, m MQTT.Message) {
	t.msgchan <- m
}

func (t *WillClient) Connect(addr string) error {
	opts := MQTT.NewClientOptions().AddBroker(addr)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetAutoReconnect(false)
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(t.cid)
	opts.SetUsername(t.username)
	opts.SetPassword(t.password)
	opts.SetProtocolVersion(4)
	opts.SetKeepAlive(time.Minute * 10)
	opts.SetWriteTimeout(time.Minute * 10)
	opts.SetPingTimeout(time.Minute * 10)
	opts.SetConnectTimeout(time.Minute * 10)
	opts.SetWill(t.topic, t.payload, byte(t.qos), t.retain)

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

func (t *WillClient) Disconnect() {
	t.cli.Disconnect(1)
}

func (t *WillClient) Recevice() MQTT.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-t.msgchan:
			return msg
		}
	}
}
