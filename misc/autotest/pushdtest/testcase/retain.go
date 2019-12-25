package testcase

import (
	"errors"

	"github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"
)

/*
## Retain 消息处理
* pub service 发布一个topic retain 消息
* c2 建立连接
* c2 增加路由信息
* c2 订阅这个topic
* c2 校验信息
* c1 建立连接
* c1 订阅这个 topic
* c1 取消订阅 topic
* pub service  发布一个清理topic 信息
* c1 重新订阅topic
* c1 断开连接
* c2 取消订阅 topic
* c2 断开连接
* c2 清理有路信息
*/

func (client *Client) Retain(publishService string) error {
	pubsrv, err := pushcli.NewPublishcli(client.etcdAddrs, publishService)
	if err != nil {
		return err
	}
	_, err = pubsrv.Publish("retain-topic", []byte("retain-topic"), true, 1)
	if err != nil {
		return err
	}

	for _, srv := range client.srvs {
		if srv.NotifyRespLen() != 0 {
			return errors.New("the num of notify is not zero")
		}
	}

	p2, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetClientID("c2"), pushcli.SetCleanSession())
	if err != nil {
		return err
	}

	//链接
	c, err := p2.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, false, 0); err != nil {
		return err
	}

	//路由信息
	_, err = p2.AddRoute("retain-topic", client.addrs[0])
	if err != nil {
		return err
	}

	//订阅
	s, err := p2.Subscribe([]string{"retain-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s, 1, []byte("")); err != nil {
		return err
	}

	p1, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetClientID("c1"), pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	//链接
	c, err = p1.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, false, 0); err != nil {
		return err
	}

	//订阅
	s, err = p1.Subscribe([]string{"retain-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s, 1, []byte("")); err != nil {
		return err
	}

	_, err = p1.Unsubscribe([]string{"retain-topic"})
	if err != nil {
		return err
	}

	_, err = pubsrv.Publish("retain-topic", []byte(""), true, 1)
	if err != nil {
		return err
	}

	s, err = p1.Subscribe([]string{"retain-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s, 0, []byte("")); err != nil {
		return err
	}

	_, err = p1.Unsubscribe([]string{"retain-topic"})
	if err != nil {
		return err
	}

	_, err = p1.Disconnect("")
	if err != nil {
		return err
	}

	_, err = p2.Unsubscribe([]string{"retain-topic"})
	if err != nil {
		return err
	}

	_, err = p2.RemoveRoute("retain-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = p2.Disconnect("")
	if err != nil {
		return err
	}

	for _, srv := range client.srvs {
		if srv.NotifyRespLen() != 0 {
			return errors.New("the num of notify is not zero")
		}
	}
	return nil
}
