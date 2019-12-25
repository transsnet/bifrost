package scene

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	Address []string

	Topic    string
	Clientid string
	Username string
	Password string
	QoS      int
	IDPoint  int
	SubSame  bool
}

func NewClient() *Client {
	return &Client{}
}

func (cli *Client) MqttConnect(id int, handler MQTT.MessageHandler) MQTT.Client {
	idstr := fmt.Sprintf("%d", cli.IDPoint+id)

	opts := MQTT.NewClientOptions().AddBroker(cli.getAddr())
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetAutoReconnect(false)
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(cli.Clientid + "-" + idstr)
	opts.SetUsername(cli.Username)
	opts.SetPassword(cli.Password)
	opts.SetCleanSession(true)
	opts.SetProtocolVersion(4)
	opts.SetKeepAlive(time.Minute * 10)
	opts.SetWriteTimeout(time.Minute * 10)
	opts.SetPingTimeout(time.Minute * 10)
	opts.SetConnectTimeout(time.Minute * 10)
	/*
		opts.SetOnConnectHandler(func(client MQTT.Client) { // 设置OnConnect回调函数，建立连接或者自动重连时触发。
			// if token := client.Subscribe(topic, byte(t.qos), t.handler()); token.WaitTimeout(5*time.Second) && token.Error() != nil {
			// log.Printf("Subscribe failed, %s", token.Error())
			// }
		})
	*/

	mc := MQTT.NewClient(opts)
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	if cli.SubSame {
		if token := mc.Subscribe(cli.Topic, byte(cli.QoS), handler); token.Wait() && token.Error() != nil {
			log.Fatalf("%#v , %#v", token.Error(), cli.Topic)
		}
		return mc
	}

	topic := cli.Topic + "-" + idstr
	if token := mc.Subscribe(topic, byte(cli.QoS), handler); token.Wait() && token.Error() != nil {
		log.Fatalf("%#v , %#v", token.Error(), topic)
	}
	return mc
}

func (cli *Client) getAddr() string {
	id := rand.Intn(len(cli.Address))
	return cli.Address[id]
}
