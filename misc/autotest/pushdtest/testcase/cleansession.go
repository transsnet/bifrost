package testcase

import "github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"

/*
## 绘画恢复场景使用: 两个客户端登录
* c1 c2 登录
* 增加路由信息
* 增加订阅信息
* 下发消息（6 条）
* c1 客户端 拉去消息 拉去最大数量为4条.
* c2 客户端 拉去消息 拉去最大数量为2条.
* c1 ack 4条 消息 del ack 2条
* 删除c1 路由
* c1 客户端离线
* 下发一条消息
* c2 拉去 6 条消息
* c2 删除路由
* c2 离线
* c1 登录
* c1 增加路由信息
* c1 订阅信息
* c1 拉取消息 2 次 一次最多2条
* ack 两条消息.
* c1 c2 删除路由信息
* c2 离线
* c1 离线
*/

func (client *Client) CleanSession() error {
	p1, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetClientID("c1"))
	if err != nil {
		return err
	}
	p2, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetClientID("c2"))
	if err != nil {
		return err
	}

	//链接
	c, err := p1.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, true, 0); err != nil {
		return err
	}

	c, err = p2.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, true, 0); err != nil {
		return err
	}

	//路由信息,由于订阅相同信息 因此c1 增加topic 不同路由效果是一样的
	for _, addr := range client.addrs {
		_, err = p1.AddRoute("clean-topic", addr)
		if err != nil {
			return err
		}
	}

	//订阅
	s1, err := p1.Subscribe([]string{"clean-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s1, 0, []byte("")); err != nil {
		return err
	}

	s2, err := p2.Subscribe([]string{"clean-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s2, 0, []byte("")); err != nil {
		return err
	}

	//下发消息
	var indexs [][]byte
	for i := 0; i < 6; i++ {
		_, err = p1.MQTTPublish("clean-topic", 1, false)
		if err != nil {
			return err
		}
		var index []byte
		for _, srv := range client.srvs {
			n := srv.NotifyResp()
			index = n.Index
			if err := EqualNotify(n, "clean-topic"); err != nil {
				return err
			}
		}
		indexs = append(indexs, index)
	}

	//消息拉取
	resp, err := p1.Pull("clean-topic", s1.Index[0], 4)
	if err != nil {
		return err
	}
	if err := EqualPull(resp, 4, "clean-topic", false); err != nil {
		return err
	}

	resp, err = p2.Pull("clean-topic", s1.Index[0], 2)
	if err != nil {
		return err
	}
	if err := EqualPull(resp, 2, "clean-topic", false); err != nil {
		return err
	}

	//发送ack
	var i int64
	for i = 0; i < 4; i++ {
		_, err = p1.PutUnack("clean-topic", indexs[i], i+1)
		if err != nil {
			return err
		}
	}

	//回复ack
	for i = 0; i < 2; i++ {
		_, err = p1.DelUnack(i + 1)
		if err != nil {
			return err
		}
	}

	//p1 离线 删除路由
	_, err = p1.RemoveRoute("clean-topic", client.addrs[0])
	if err != nil {
		return err
	}

	_, err = p1.Unsubscribe([]string{"thord-topic"})
	if err != nil {
		return err
	}
	_, err = p1.Disconnect("")
	if err != nil {
		return err
	}

	_, err = p2.MQTTPublish("clean-topic", 1, false)
	if err != nil {
		return err
	}
	n := client.srvs[1].NotifyResp()
	indexs = append(indexs, n.Index)
	if err := EqualNotify(n, "clean-topic"); err != nil {
		return err
	}

	//拉取消息
	resp, err = p2.Pull("clean-topic", s2.Index[0], 6)
	if err != nil {
		return err
	}
	if err := EqualPull(resp, 6, "clean-topic", false); err != nil {
		return err
	}

	//p2 离线 删除路由
	_, err = p2.RemoveRoute("clean-topic", client.addrs[1])
	if err != nil {
		return err
	}
	_, err = p2.Unsubscribe([]string{"thord-topic"})
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
	if err := EqualConnect(c, 1, true, 4); err != nil {
		return err
	}

	if err := EqualRecords(c.Records[0], indexs[4], indexs[6], "clean-topic"); err != nil {
		return err
	}

	//拉取出未读ack消息
	r, err := p1.RangeUnack([]byte{0}, 1)
	if err != nil {
		return err
	}
	if err := EqualRange(r, 1, false, []byte{4}, "clean-topic"); err != nil {
		return err
	}

	r, err = p1.RangeUnack(r.Offset, 2)
	if err != nil {
		return err
	}
	if err := EqualRange(r, 1, true, []byte{5}, "clean-topic"); err != nil {
		return err
	}

	_, err = p1.AddRoute("clean-topic", client.addrs[0])
	if err != nil {
		return err
	}

	//订阅
	s1, err = p1.Subscribe([]string{"clean-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s1, 0, []byte("")); err != nil {
		return err
	}

	_, err = p1.MQTTPublish("clean-topic", 1, false)
	if err != nil {
		return err
	}
	n = client.srvs[0].NotifyResp()
	indexs = append(indexs, n.Index)
	if err := EqualNotify(n, "clean-topic"); err != nil {
		return err
	}

	//p1 离线 删除路由
	_, err = p1.RemoveRoute("clean-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = p1.Unsubscribe([]string{"thord-topic"})
	if err != nil {
		return err
	}

	_, err = p1.Disconnect("")
	if err != nil {
		return err
	}
	return nil
}
