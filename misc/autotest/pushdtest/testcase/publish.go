package testcase

import (
	"errors"

	"github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"
)

/*
## Publish
* 下发消息
* 校验信息
* 客户端登录
* 增加路由信息
* c1 c2 订阅topic
* 下发消息（5 条）
* c1 客户端 拉去消息 拉去最大数量为2条 拉取两次.
* ack 两条消息.
* c2 客户端 拉去消息 拉去最大数量为2条.
* ack 两条消息.
* 删除c1 路由
* c1 客户端离线
* 下发一条消息
* c2 拉取消息 2 次 一次最多2条
* c2 删除路由
* c2 离线
* 客户端登录 * c1 客户端离线
* c1 增加路由信息
* c1 订阅topic
* 发布一条消息
* 删除路由
* c1 客户端离线
*/

func (client *Client) Publish(publishService string) error {
	pubsrv, err := pushcli.NewPublishcli(client.etcdAddrs, publishService)
	if err != nil {
		return err
	}
	pr, err := pubsrv.Publish("publish-topic", []byte(""), false, 1)
	if err := EqualPublish(pr, 0, 1, 1); err != nil {
		return err
	}

	p1, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetClientID("c1"), pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	p2, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetClientID("c2"), pushcli.SetCleanSession())
	if err != nil {
		return err
	}

	//链接
	c, err := p1.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, false, 0); err != nil {
		return err
	}

	c, err = p2.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, false, 0); err != nil {
		return err
	}

	//路由信息,由于订阅相同信息 因此c1 增加topic 不同路由效果是一样的
	for _, addr := range client.addrs {
		_, err = p1.AddRoute("publish-topic", addr)
		if err != nil {
			return err
		}
	}

	//订阅
	s1, err := p1.Subscribe([]string{"publish-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s1, 0, []byte("")); err != nil {
		return err
	}

	s2, err := p2.Subscribe([]string{"publish-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s2, 0, []byte("")); err != nil {
		return err
	}

	//下发消息
	var indexs [][]byte
	for i := 0; i < 5; i++ {
		pr, err = pubsrv.Publish("publish-topic", []byte(""), false, 1)
		if err != nil {
			return err
		}
		if err := EqualPublish(pr, 0, 1, 0); err != nil {
			return err
		}
		var index []byte
		for _, srv := range client.srvs {
			n := srv.NotifyResp()
			if err := EqualNotify(n, "publish-topic"); err != nil {
				return err
			}
			index = n.Index
		}
		indexs = append(indexs, index)
	}

	//消息拉取
	count := 2
	resp, err := p1.Pull("publish-topic", s1.Index[0], int64(count))
	if err != nil {
		return err
	}
	if err := EqualPull(resp, count, "publish-topic", false); err != nil {
		return err
	}

	resp, err = p1.Pull("publish-topic", indexs[2], int64(count))
	if err != nil {
		return err
	}
	if err := EqualPull(resp, count, "publish-topic", false); err != nil {
		return err
	}

	resp, err = p2.Pull("publish-topic", s2.Index[0], int64(count))
	if err != nil {
		return err
	}
	if err := EqualPull(resp, count, "publish-topic", false); err != nil {
		return err
	}

	//发送ack
	var i int64
	for i = 1; i < 5; i++ {
		_, err = p1.PutUnack("publish-topic", indexs[i], i+1)
		if err != nil {
			return err
		}
	}

	//回复ack
	for i = 1; i < 5; i++ {
		_, err = p1.DelUnack(i + 1)
		if err != nil {
			return err
		}
	}

	//发送ack
	for i = 1; i < 3; i++ {
		_, err = p2.PutUnack("publish-topic", indexs[i], i+1)
		if err != nil {
			return err
		}
	}

	//回复ack
	for i = 1; i < 3; i++ {
		_, err = p2.DelUnack(i + 1)
		if err != nil {
			return err
		}
	}

	//p1 离线 删除路由
	_, err = p1.RemoveRoute("publish-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = p1.Unsubscribe([]string{"publish-topic"})
	if err != nil {
		return err
	}

	_, err = p1.Disconnect("")
	if err != nil {
		return err
	}

	pr, err = pubsrv.Publish("publish-topic", []byte(""), false, 1)
	if err != nil {
		return err
	}
	if err := EqualPublish(pr, 0, 1, 0); err != nil {
		return err
	}
	n := client.srvs[1].NotifyResp()
	indexs = append(indexs, n.Index)
	if err := EqualNotify(n, "publish-topic"); err != nil {
		return err
	}

	//拉取消息
	resp, err = p2.Pull("publish-topic", indexs[2], 4)
	if err != nil {
		return err
	}
	if err := EqualPull(resp, 4, "publish-topic", false); err != nil {
		return err
	}

	//p2 离线 删除路由
	_, err = p2.RemoveRoute("publish-topic", client.addrs[1])
	if err != nil {
		return err
	}
	_, err = p2.Unsubscribe([]string{"publish-topic"})
	if err != nil {
		return err
	}

	_, err = p2.Disconnect("")
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

	_, err = p1.AddRoute("publish-topic", client.addrs[0])
	if err != nil {
		return err
	}

	//订阅
	s1, err = p1.Subscribe([]string{"publish-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s1, 0, []byte("")); err != nil {
		return err
	}

	pr, err = pubsrv.Publish("publish-topic", []byte(""), false, 1)
	if err != nil {
		return err
	}
	if err := EqualPublish(pr, 0, 1, 0); err != nil {
		return err
	}

	n = client.srvs[0].NotifyResp()
	if err := EqualNotify(n, "publish-topic"); err != nil {
		return err
	}

	//p2 离线 删除路由
	_, err = p1.RemoveRoute("publish-topic", client.addrs[0])
	if err != nil {
		return err
	}

	_, err = p1.Unsubscribe([]string{"publish-topic"})
	if err != nil {
		return err
	}

	_, err = p1.Disconnect("")
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
