package mqtttest

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	pub "github.com/meitu/bifrost/grpc/publish"
)

func EqualByte(c *MQTTClient, p *PublishClient) error {
	msg := c.Recevice()
	if msg == nil {
		return errors.New(ErrMqttMsg)
	}
	if !bytes.Equal(msg.Payload(), p.Payload()) {
		return errors.New(fmt.Sprintf(ErrPayload, msg.Payload(), p.Payload()))
	}
	return nil
}

func EqualPay(c *MQTTClient, payload []byte) error {
	msg := c.Recevice()
	if msg == nil {
		return errors.New(ErrMqttMsg)
	}
	if !bytes.Equal(msg.Payload(), payload) {
		return errors.New(fmt.Sprintf(ErrPayload, msg.Payload(), payload))
	}
	return nil
}

func EqualWillPay(c *WillClient, payload []byte) error {
	msg := c.Recevice()
	if msg == nil {
		return errors.New(ErrMqttMsg)
	}
	if !bytes.Equal(msg.Payload(), payload) {
		return errors.New(fmt.Sprintf(ErrPayload, msg.Payload(), payload))
	}
	return nil
}

func Payload() []byte {
	return []byte(time.Now().String())
}

func Target(t string, retain, ndg bool) *pub.Target {
	return &pub.Target{
		Topic:         t,
		Qos:           1,
		IsRetain:      retain,
		NoneDowngrade: ndg,
	}
}
