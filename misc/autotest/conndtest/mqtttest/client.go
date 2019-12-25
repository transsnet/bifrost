package mqtttest

import (
	"errors"
	"fmt"
	"time"

	pub "github.com/meitu/bifrost/grpc/publish"
)

var (
	t1   = "t1"
	t2   = "t2"
	cid1 = "c1"
	cid2 = "c2"
	qos  = 1
)

const (
	ErrMqttConnect    = "mqtt connect failed ,%s"
	ErrPublishConnect = "publish connect failed, %s"
	ErrPublishRequest = "publish Request failed, %s"
	ErrMqttMsg        = "mqtt message lost"
	ErrPayload        = "expect payload %s,actual payload %s"
	ErrMsgAbnormal    = "msg is not zero"
)

type Client struct {
	mqttsAddr string
	pubsAddr  string
	// Pubc      *PublishClient
	// Mqttc     []*MQTTClient
}

func NewClient(pubsAddr, mqttsAddr string) *Client {
	return &Client{
		mqttsAddr: mqttsAddr,
		pubsAddr:  pubsAddr,
	}
}

func (cli *Client) IMLiveCase() error {
	if err := cli.clearRetain(t1); err != nil {
		return err
	}

	c1, err := NewMQTTClient(cid1, t1, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2, err := NewMQTTClient(cid2, t1, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	p1, err := NewPublishClient(cli.pubsAddr, []byte(Payload()))
	if err != nil {
		return errors.New(fmt.Sprintf(ErrPublishConnect, err))
	}
	tar1 := Target(t1, false, false)
	if err := p1.Request([]*pub.Target{tar1}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	if err := EqualByte(c1, p1); err != nil {
		return err
	}
	if err := EqualByte(c2, p1); err != nil {
		return err
	}
	c1.Disconnect()
	if err := p1.Request([]*pub.Target{tar1}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	if err := EqualByte(c2, p1); err != nil {
		return err
	}

	c2.Disconnect()
	if err := p1.Request([]*pub.Target{tar1}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	if c1.MsgLen() != 0 || c2.MsgLen() != 0 {
		return errors.New(ErrMsgAbnormal)
	}
	return nil
}

func (cli *Client) PushCase() error {
	if err := cli.clearRetain(t1); err != nil {
		return err
	}

	if err := cli.clearRetain(t2); err != nil {
		return err
	}
	c1, err := NewMQTTClient(cid1, t1, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2, err := NewMQTTClient(cid2, t2, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	p1, err := NewPublishClient(cli.pubsAddr, []byte(Payload()))
	if err != nil {
		return errors.New(fmt.Sprintf(ErrPublishConnect, err))
	}

	tar1 := Target(t1, false, false)
	tar2 := Target(t2, false, false)
	if err := p1.Request([]*pub.Target{tar1, tar2}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	if err := EqualByte(c1, p1); err != nil {
		return err
	}
	if err := EqualByte(c2, p1); err != nil {
		return err
	}
	c1.Disconnect()
	c2.Disconnect()
	if err := p1.Request([]*pub.Target{tar1, tar2}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	if c1.MsgLen() != 0 || c2.MsgLen() != 0 {
		return errors.New(ErrMsgAbnormal)
	}
	return nil
}

func (cli *Client) Retain() error {
	if err := cli.clearRetain(t1); err != nil {
		return err
	}
	if err := cli.retainCaseOne(); err != nil {
		return err
	}

	if err := cli.retainCaseTwo(); err != nil {
		return err
	}
	return nil
}

func (cli *Client) WillMsg() error {
	if err := cli.clearRetain(t1); err != nil {
		return err
	}
	if err := cli.willCaseOne(); err != nil {
		return err
	}
	if err := cli.willCaseTwo(); err != nil {
		return err
	}
	return nil
}

func (cli *Client) CleanSession() error {
	if err := cli.clearRetain(t1); err != nil {
		return err
	}
	c2, err := NewMQTTClient(cid2, t1, qos, cli.mqttsAddr, false)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2.Disconnect()
	time.Sleep(time.Millisecond * 10)

	p1, err := NewPublishClient(cli.pubsAddr, []byte(Payload()))
	if err != nil {
		return errors.New(fmt.Sprintf(ErrPublishConnect, err))
	}
	tar1 := Target(t1, false, false)
	if err := p1.Request([]*pub.Target{tar1}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}

	c2, err = NewMQTTClient(cid2, t1, qos, cli.mqttsAddr, false)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	if err := EqualByte(c2, p1); err != nil {
		return err
	}

	c2.Disconnect()
	return nil
}

func (cli Client) retainCaseOne() error {
	p1, err := NewPublishClient(cli.pubsAddr, []byte(Payload()))
	if err != nil {
		return errors.New(fmt.Sprintf(ErrPublishConnect, err))
	}
	tar1 := Target(t1, true, false)
	if err := p1.Request([]*pub.Target{tar1}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	c1, err := NewMQTTClient(cid1, t1, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2, err := NewMQTTClient(cid2, t2, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	if err := EqualByte(c1, p1); err != nil {
		return err
	}
	c1.Disconnect()
	c2.Disconnect()
	if c1.MsgLen() == 0 && c2.MsgLen() == 0 {
		return nil
	}
	return errors.New(ErrMsgAbnormal)

}

func (cli *Client) retainCaseTwo() error {
	p1, err := NewPublishClient(cli.pubsAddr, []byte(""))
	if err != nil {
		return errors.New(fmt.Sprintf(ErrPublishConnect, err))
	}
	tar1 := Target(t1, true, false)
	tar2 := Target(t2, true, false)
	if err := p1.Request([]*pub.Target{tar1, tar2}); err != nil {
		return errors.New(fmt.Sprintf(ErrPublishRequest, err))
	}
	time.Sleep(time.Second)

	c1, err := NewMQTTClient(cid1, t1, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2, err := NewMQTTClient(cid2, t2, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c1.Disconnect()
	c2.Disconnect()
	if c1.MsgLen() == 0 && c2.MsgLen() == 0 {
		return nil
	}
	fmt.Println(c1.MsgLen(), c2.MsgLen())
	return errors.New(ErrMsgAbnormal)
}

func (cli *Client) willCaseOne() error {
	c1, err := NewWillClient(cid1, t2, qos, cli.mqttsAddr, false, "hello")
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2, err := NewMQTTClient(cid2, t2, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c1.Disconnect()
	if err := EqualPay(c2, []byte("hello")); err != nil {
		return err
	}
	return nil
}

func (cli *Client) willCaseTwo() error {
	c1, err := NewWillClient(cid1, t2, qos, cli.mqttsAddr, true, "hello")
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c2, err := NewMQTTClient(cid2, t2, qos, cli.mqttsAddr, true)
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	c1.Disconnect()
	if err := EqualPay(c2, []byte("hello")); err != nil {
		return err
	}

	c1, err = NewWillClient(cid1, t2, qos, cli.mqttsAddr, true, "hello")
	if err != nil {
		return errors.New(fmt.Sprintf(ErrMqttConnect, err))
	}
	if err := EqualWillPay(c1, []byte("hello")); err != nil {
		return err
	}
	return nil
}

func (cli *Client) clearRetain(t string) error {
	p1, err := NewPublishClient(cli.pubsAddr, []byte(""))
	if err != nil {
		return errors.New(fmt.Sprintf("clean retain new pub failed,%s", err))
	}
	tar1 := Target(t, true, false)
	if err := p1.Request([]*pub.Target{tar1}); err != nil {
		return errors.New(fmt.Sprintf("clean retain request failed,%s", err))
	}
	return nil
}
