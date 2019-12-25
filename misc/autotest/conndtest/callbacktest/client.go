package callbacktest

import (
	"errors"
	"fmt"
	"time"

	"github.com/meitu/bifrost/misc/fake/callback"
)

type Client struct {
	mqttsAddr string
	pubsAddr  string
	etcdAddr  []string
}

//TODO username == appkey
func NewClient(pubsAddr, mqttsAddr string, etcdAddr []string) *Client {
	return &Client{
		mqttsAddr: mqttsAddr,
		pubsAddr:  pubsAddr,
		etcdAddr:  etcdAddr,
	}
}

func (cli *Client) Service(service string) error {
	if err := cli.serviceCase(service); err != nil {
		return err
	}
	return nil
}

func (cli *Client) serviceCase(service string) error {
	// username = "bifrost-appkey=service-1544516816-1-bb7c77b22be1ff945d3272"
	c1, err := NewMQTTClient(cli.mqttsAddr, "bifrost-appkey=s1-1544766957-1-4dee5844225633c0527479", "p")
	if err != nil {
		return err
	}
	fnames := "OnConnect,PostSubscribe,OnDisconnect,OnSubscribe,OnUnsubscribe,OnACK,OnPublish,OnOffline"

	call, err := callback.NewCallbackServer(cli.etcdAddr, service, "s1", fnames)
	if err != nil {
		return err
	}

	go func() {
		call.Start("localhost:7744")
	}()
	time.Sleep(time.Second)

	fnames = "PostSubscribe,OnDisconnect,OnSubscribe"
	call2, err := callback.NewCallbackServer(cli.etcdAddr, service, "s1", fnames)
	if err != nil {
		return err
	}

	go func() {
		call2.Start("localhost:12323")
	}()
	call2.Stop()

	if err := c1.Connect(); err != nil {
		return fmt.Errorf("connect %s", err.Error())
	}

	if err := CallSubscribe(c1, call, 1); err != nil {
		return fmt.Errorf("subscribe %s", err.Error())
	}
	if err := CallUnsubscrib(c1, call, 0); err != nil {
		return fmt.Errorf("unsubscribe %s", err.Error())
	}

	if err := CallDisconnect(c1, call, 1); err != nil {
		return fmt.Errorf("disconnect %s", err.Error())
	}

	call.Stop()
	if call.RecvLen() != 0 {
		return errors.New("recv is no equal zero")
	}
	return nil
}

func (cli *Client) CookieEuqal(service string) error {
	c1, err := NewMQTTClient(cli.mqttsAddr, "bifrost-appkey=s1-1544766957-1-4dee5844225633c0527479", "p")
	call, err := callback.NewConndCookieServer(cli.etcdAddr, service, "s1")
	if err != nil {
		return err
	}
	go func() {
		call.Start("localhost:23231")
	}()

	time.Sleep(50 * time.Millisecond)
	if err := c1.Connect(); err != nil {
		return fmt.Errorf("connect %s", err.Error())
	}
	if err := c1.Subscribe("t"); err != nil {
		return fmt.Errorf("subscribe %s", err.Error())
	}
	if err := c1.Publish("t"); err != nil {
		return fmt.Errorf("publish %s", err.Error())
	}
	if err := c1.UnSubscribe("t"); err != nil {
		return fmt.Errorf("unsubscribe %s", err.Error())
	}
	c1.Disconnect()
	call.Stop()
	return nil
}

func CallConnect(c1 *MQTTClient, call *callback.Call, expect int) error {
	if err := c1.Connect(); err != nil {
		return err
	}
	if err := EuqalRecv(call, expect); err != nil {
		return err
	}
	return nil
}

func CallSubscribe(c1 *MQTTClient, call *callback.Call, expect int) error {
	if err := c1.Subscribe("t"); err != nil {
		return err
	}
	if err := EuqalRecv(call, expect); err != nil {
		return err
	}

	if err := EuqalRecv(call, expect); err != nil {
		return err
	}
	return nil
}

func CallPublish(c1 *MQTTClient, call *callback.Call, expect int) error {
	if err := c1.Publish("t"); err != nil {
		return err
	}
	if err := EuqalRecv(call, expect); err != nil {
		return err
	}
	if err := EuqalRecv(call, expect); err != nil {
		return err
	}

	return nil
}

func CallUnsubscrib(c1 *MQTTClient, call *callback.Call, expect int) error {
	if err := c1.UnSubscribe("t"); err != nil {
		return err
	}

	if err := EuqalRecv(call, expect); err != nil {
		return err
	}
	return nil
}

func CallDisconnect(c1 *MQTTClient, call *callback.Call, expect int) error {
	c1.Disconnect()
	if err := EuqalRecv(call, expect); err != nil {
		return err
	}
	return nil
}

func EuqalRecv(call *callback.Call, num int) error {
	enum := call.Recv()
	if enum != num {
		return errors.New(fmt.Sprintf("expect recv msg %d,autual recv msg %d", num, enum))
	}
	return nil
}
