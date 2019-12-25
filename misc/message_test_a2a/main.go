package main

import (
	"crypto/tls"
	"flag"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Setting is the setting option of this test case.
type Setting struct {
	Addr          string
	ClientID      string
	Username      string
	Password      string
	CleanSession  bool
	AutoReconnect bool
}

func main() {
	var setting Setting

	flag.StringVar(&setting.Addr, "h", "tcp://127.0.0.1:1883", "the MQTT service address")
	flag.StringVar(&setting.ClientID, "c", "ssl-sample-sender", "the client ID")
	flag.StringVar(&setting.Username, "u", "username", "the username")
	flag.StringVar(&setting.Password, "p", "password", "the password")
	flag.BoolVar(&setting.CleanSession, "s", true, "true for enable clean session")
	flag.BoolVar(&setting.AutoReconnect, "r", false, "true for enable auto reconnect")
	flag.Parse()

	opts := MQTT.NewClientOptions()
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)
	opts.AddBroker(setting.Addr)
	opts.SetClientID(setting.ClientID).SetUsername(setting.Username).SetPassword(setting.Password).SetProtocolVersion(4)
	opts.SetCleanSession(setting.CleanSession)
	opts.SetAutoReconnect(setting.AutoReconnect)

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
		if string(m.Payload()) == messages[readCount] {
			readCount++
			fmt.Println(readCount, string(m.Payload()))
			if readCount == len(messages) {
				close(signal)
			}
		}
	}

	opts.SetDefaultPublishHandler(onMessage)

	cli := MQTT.NewClient(opts)

	fmt.Println("connect with server")
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("connect failed, %s\n", token.Error())
		return
	}

	fmt.Println("subscribe topic")
	if token := cli.Subscribe("/test/topic", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("subscriber failed, %s\n", token.Error())
		return
	}

	fmt.Println("publish some message")
	for i := range messages {
		if token := cli.Publish("/test/topic", 1, false, []byte(messages[i])); token.Wait() && token.Error() != nil {
			fmt.Printf("publish failed, %s\n", token.Error())
			return
		}
	}

	fmt.Println("wait for the published message")
	<-signal

	fmt.Println("unsubscribe topic")
	if token := cli.Unsubscribe("/test/topic"); token.Wait() && token.Error() != nil {
		fmt.Printf("Unsubscribe failed, %s\n", token.Error())
		return
	}

	// Reset signal message
	signal = make(chan interface{})
	readCount = 0

	fmt.Println("subscribe topic again")
	if token := cli.Subscribe("/test/topic", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("subscriber failed, %s\n", token.Error())
		return
	}

	fmt.Println("publish some message again")
	for i := range messages {
		if token := cli.Publish("/test/topic", 1, false, []byte(messages[i])); token.Wait() && token.Error() != nil {
			fmt.Printf("publish failed, %s\n", token.Error())
			return
		}
	}

	fmt.Println("wait for the published message the second times")
	<-signal

	fmt.Println("unsubscribe topic again")
	if token := cli.Unsubscribe("/test/topic"); token.Wait() && token.Error() != nil {
		fmt.Printf("Unsubscribe failed, %s\n", token.Error())
		return
	}

	fmt.Println("test pass!!!!")

	cli.Disconnect(0)
}
