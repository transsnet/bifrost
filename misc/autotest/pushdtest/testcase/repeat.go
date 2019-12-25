package testcase

import (
	"github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"
)

/*
### 重复取消订阅
### 重复订阅
### 重复增加路由信息
### 重复删除ack
### 重复增加ack
### 重复断链
### 重复删除路由信息
*/
func (client *Client) Repeat() error {
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
	_, err = cli.AddRoute("repeat-topic", client.addrs[0])
	if err != nil {
		return err
	}

	_, err = cli.AddRoute("repeat-topic", client.addrs[0])
	if err != nil {
		return err
	}

	//订阅
	s, err := cli.Subscribe([]string{"repeat-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s, 0, []byte("")); err != nil {
		return err
	}

	s, err = cli.Subscribe([]string{"repeat-topic"}, []int32{1})
	if err != nil {
		return err
	}
	if err := EqualSubscribe(s, 0, []byte("")); err != nil {
		return err
	}

	//下发消息
	_, err = cli.MQTTPublish("repeat-topic", 1, false)
	if err != nil {
		return err
	}
	n := client.srvs[0].NotifyResp()
	if err := EqualNotify(n, "repeat-topic"); err != nil {
		return err
	}

	//消息拉取
	resp, err := cli.Pull("repeat-topic", s.Index[0], 2)
	if err != nil {
		return err
	}
	if err := EqualPull(resp, 1, "repeat-topic", false); err != nil {
		return err
	}

	resp, err = cli.Pull("repeat-topic", s.Index[0], 2)
	if err != nil {
		return err
	}
	if err := EqualPull(resp, 1, "repeat-topic", false); err != nil {
		return err
	}

	//发送ack
	_, err = cli.PutUnack("repeat-topic", n.Index, 2)
	if err != nil {
		return err
	}
	_, err = cli.PutUnack("repeat-topic", n.Index, 2)
	if err != nil {
		return err
	}

	//回复ack
	_, err = cli.DelUnack(2)
	if err != nil {
		return err
	}
	/*
		_, err = cli.DelUnack(2 + 1)
		if err != nil {
			return err
		}
	*/

	_, err = cli.RemoveRoute("repeat-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = cli.RemoveRoute("repeat-topic", client.addrs[0])
	if err != nil {
		return err
	}
	_, err = cli.Unsubscribe([]string{"repeat-topic"})
	if err != nil {
		return err
	}
	_, err = cli.Unsubscribe([]string{"repeat-topic"})
	if err != nil {
		return err
	}
	_, err = cli.Disconnect("")
	if err != nil {
		return err
	}
	_, err = cli.Disconnect("")
	if err != nil {
		return err
	}
	return nil
}
