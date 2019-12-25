package main

import (
	"flag"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Setting struct {
	address       string
	sender        string
	sUsername     string
	sPassword     string
	receiver      string
	rUsername     string
	rPassword     string
	autoReconnect bool
}

func main() {
	var setting Setting
	flag.StringVar(&setting.address, "a", "tcp://127.0.0.1:1883", "the MQTT broker address")
	flag.StringVar(&setting.sender, "sid", "mqtt-sender", "the MQTT sender client ID of this test case")
	flag.StringVar(&setting.sUsername, "su", "sender-username", "the MQTT sender username of this test case")
	flag.StringVar(&setting.sPassword, "sp", "sender-password", "the MQTT sender password of this test case")
	flag.StringVar(&setting.receiver, "rid", "mqtt-receiver", "the MQTT recevier client ID of this test case")
	flag.StringVar(&setting.rUsername, "ru", "receiver-username", "the MQTT receiver username of this test case")
	flag.StringVar(&setting.rPassword, "rp", "recevier-password", "the MQTT receiver password of this test case")
	flag.BoolVar(&setting.autoReconnect, "r", false, "true for enabling auto reconnect")
	flag.Parse()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(setting.address)
	opts.SetClientID(setting.receiver).SetUsername(setting.rUsername).SetPassword(setting.rPassword).SetProtocolVersion(4)
	opts.SetCleanSession(false)
	opts.SetAutoReconnect(setting.autoReconnect)

	signal := make(chan interface{})
	messages := []string{
		"hello1",
		"hello2",
		"hello3",
		"hello4",
		"hello5",
		"hello6",
	}
	readCount := 0

	onMessage := func(cli MQTT.Client, m MQTT.Message) {
		if readCount < 6 && string(m.Payload()) == messages[readCount] {
			readCount++
			if readCount == len(messages) {
				close(signal)
			}
		}
	}

	opts.SetDefaultPublishHandler(onMessage)

	fmt.Printf("sub connect to server with clean session flag = false\n")
	sub := MQTT.NewClient(opts)
	if token := sub.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("connect failed, %s\n", token.Error())
		return
	}

	fmt.Printf("sub subscribe a topic with qos = 1\n")
	if token := sub.Subscribe("/test/ack", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("subscribe failed, %s\n", token.Error())
		return
	}

	fmt.Printf("sub disconnect\n")
	sub.Disconnect(50)

	opts.SetClientID(setting.sender).SetUsername(setting.sUsername).SetPassword(setting.sPassword).SetProtocolVersion(4)
	pub := MQTT.NewClient(opts)
	fmt.Printf("pub connect to server\n")
	if token := pub.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("pub connect failed, %s\n", token.Error())
		return
	}

	fmt.Printf("pub publish some message to the topic sub subscribed\n")
	for i := range messages {
		if token := pub.Publish("/test/ack", 1, false, []byte(messages[i])); token.Wait() && token.Error() != nil {
			fmt.Printf("publish failed, %s\n", token.Error())
		}
	}

	fmt.Printf("pub disconnect\n")
	pub.Disconnect(50)

	opts.SetClientID(setting.receiver).SetUsername(setting.rUsername).SetPassword(setting.rPassword).SetProtocolVersion(4)
	subl := MQTT.NewClient(opts)
	fmt.Printf("sub reconnect to server with clean session flag = false\n")
	if token := subl.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("sub reconnect failed, %s\n", token.Error())
		return
	}

	fmt.Printf("waiting for offline message...\n")
	<-signal
	fmt.Printf("received offline message, test pass...\n")

	if token := subl.Unsubscribe("/test/ack"); token.Wait() && token.Error() != nil {
		fmt.Printf("subl unsubscribe failed, %s\n", token.Error())
		return
	}
}
