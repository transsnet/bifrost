package testcase

import "github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"

/*
## 推送测试场景：客户端以cleansession为true 的情况登录，订阅一个topic，下发消息。
* 增加路由信息
* 客户端登录
* 订阅topic
* 下发消息（5 条消息）
* 拉去消息 拉去最大数量为2条 拉取两次 : 每次会返回两条消息.
* 回复ack 回去最前面的两条数据ack
* 客户端离线
* 客户端重新登录
* 下发一条消息
* 拉取消息 : 仅仅可以得到一条消息
* 删除路由信息
*/

func (client *Client) Thord() error {
	cli, err := pushcli.NewPushdcli(client.etcdAddrs, client.service, pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	//链接
	c, err := cli.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, false, 0); err != nil {
		return err
	}

	//路由信息
	_, err = cli.AddRoute("thord-topic", client.addrs[0])
	if err != nil {
		return err
	}

	//订阅
	s, err := cli.Subscribe([]string{"thord-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s, 0, []byte("")); err != nil {
		return err
	}

	//下发消息
	var indexs [][]byte
	for i := 0; i < 5; i++ {
		_, err = cli.MQTTPublish("thord-topic", 1, false)
		if err != nil {
			return err
		}
		n := client.srvs[0].NotifyResp()
		indexs = append(indexs, n.Index)
		if err := EqualNotify(n, "thord-topic"); err != nil {
			return err
		}
	}

	//消息拉取
	count := 2
	resp, err := cli.Pull("thord-topic", s.Index[0], int64(count))
	if err != nil {
		return err
	}
	if err := EqualPull(resp, count, "thord-topic", false); err != nil {
		return err
	}

	//消息拉取
	resp, err = cli.Pull("thord-topic", indexs[2], int64(count))
	if err != nil {
		return err
	}
	if err := EqualPull(resp, count, "thord-topic", false); err != nil {
		return err
	}

	//发送ack
	var i int64
	for i = 1; i < 5; i++ {
		_, err = cli.PutUnack("thord-topic", indexs[i], i+1)
		if err != nil {
			return err
		}
	}

	//回复ack
	for i = 1; i < 5; i++ {
		_, err = cli.DelUnack(i + 1)
		if err != nil {
			return err
		}
	}

	_, err = cli.RemoveRoute("thord-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = cli.Unsubscribe([]string{"thord-topic"})
	if err != nil {
		return err
	}
	_, err = cli.Disconnect("")
	if err != nil {
		return err
	}

	c, err = cli.Connect("")
	if err != nil {
		return err
	}
	if err := EqualConnect(c, 0, false, 0); err != nil {
		return err
	}

	//路由信息
	_, err = cli.AddRoute("thord-topic", client.addrs[0])
	if err != nil {
		return err
	}

	s, err = cli.Subscribe([]string{"thord-topic"}, []int32{1})
	if err := EqualSubscribe(s, 0, []byte("")); err != nil {
		return err
	}

	_, err = cli.MQTTPublish("thord-topic", 1, false)
	if err != nil {
		return err
	}

	n := client.srvs[0].NotifyResp()
	if err := EqualNotify(n, "thord-topic"); err != nil {
		return err
	}

	_, err = cli.RemoveRoute("thord-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = cli.Unsubscribe([]string{"thord-topic"})
	if err != nil {
		return err
	}

	_, err = cli.Disconnect("")
	if err != nil {
		return err
	}
	return nil
}
