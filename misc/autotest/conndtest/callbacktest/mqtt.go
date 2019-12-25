package callbacktest

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	cli      MQTT.Client
	username string
	password string
	cid      string

	msgchan chan MQTT.Message
}

func NewMQTTClient(addr string, username, password string) (*MQTTClient, error) {
	mqttc := &MQTTClient{
		cid:      "clientid",
		msgchan:  make(chan MQTT.Message, 100),
		username: username,
		password: password,
	}
	if err := mqttc.NewClient(addr); err != nil {
		return nil, err
	}
	return mqttc, nil
}

func (t *MQTTClient) handler(c MQTT.Client, m MQTT.Message) {
	t.msgchan <- m
}

func (t *MQTTClient) NewClient(addr string) error {
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
	t.cli = MQTT.NewClient(opts)
	return nil
}

func (t *MQTTClient) Connect() error {
	if token := t.cli.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Connect failed, %s", token.Error())
		return token.Error()
	}
	return nil
}

func (t *MQTTClient) Subscribe(topic string) error {
	if token := t.cli.Subscribe(topic, byte(1), t.handler); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		log.Printf("Subscribe failed, %s", token.Error())
		return token.Error()
	}
	return nil
}

func (t *MQTTClient) UnSubscribe(topic string) error {
	if token := t.cli.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println("unsub", token)
		return token.Error()
	}
	return nil
}
func (t *MQTTClient) Publish(topic string) error {
	if token := t.cli.Publish(topic, byte(1), false, "hello"); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *MQTTClient) Disconnect() {
	t.cli.Disconnect(1)
}

func (t *MQTTClient) Recevice() MQTT.Message {
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
	return nil
}

func (t *MQTTClient) MsgLen() int {
	return len(t.msgchan)
}
