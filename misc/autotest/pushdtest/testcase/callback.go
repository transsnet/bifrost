package testcase

import (
	"fmt"
	"time"

	"github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"
	"github.com/meitu/bifrost/misc/fake/callback"
)

func (cli *Client) CallbackFunc(service string) error {
	if err := cli.call(service); err != nil {
		return err
	}
	if err := cli.call2(service); err != nil {
		return err
	}
	return nil
}

func (cli *Client) call(service string) error {
	p1, err := pushcli.NewPushdcli(cli.etcdAddrs, cli.service, pushcli.SetService("s1"))
	if err != nil {
		return err
	}

	fnames := "OnConnect,PostSubscribe,OnDisconnect,OnSubscribe,OnUnsubscribe,OnACK,OnPublish,OnOffline"
	c1, err := callback.NewCallbackServer(cli.etcdAddrs, service, "s1", fnames)
	if err != nil {
		return err
	}
	go func() {
		c1.Start("localhost:7742")
	}()
	time.Sleep(time.Second)
	p1.Connect("")
	p1.PostSubscribe()
	p1.Subscribe([]string{"topic"}, []int32{1})
	p1.MQTTPublish("topic", 1, false)
	p1.PutUnack("topic", []byte("1000000000000001"), 1)
	p1.DelUnack(1)
	p1.Unsubscribe([]string{"topic"})
	p1.Disconnect("")
	if c1.RecvLen() != 7 {
		return fmt.Errorf("expect the count of callback %v ,autul the count of callback %v", 7, c1.RecvLen())
	}
	c1.Stop()
	return nil
}

func (cli *Client) call2(service string) error {
	p1, err := pushcli.NewPushdcli(cli.etcdAddrs, cli.service, pushcli.SetService("s2"))
	if err != nil {
		return err
	}

	fnames := "OnConnect,PostSubscribe,OnDisconnect,OnSubscribe,OnUnsubscribe"
	c1, err := callback.NewCallbackServer(cli.etcdAddrs, service, "s2", fnames)
	if err != nil {
		return err
	}
	go func() {
		c1.Start("localhost:7744")
	}()
	time.Sleep(time.Second)
	p1.Connect("")
	p1.PostSubscribe()
	p1.Subscribe([]string{"topic"}, []int32{1})
	p1.MQTTPublish("topic", 1, false)
	p1.PutUnack("topic", []byte("1000000000000001"), 1)
	p1.DelUnack(1)
	p1.Disconnect("")
	if c1.RecvLen() != 4 {
		return fmt.Errorf("expect the count of calling %v ,autul the count of calling %v", 4, c1.RecvLen())
	}
	c1.Stop()
	return nil
}
