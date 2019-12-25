package testcase

import (
	"fmt"

	"github.com/meitu/bifrost/misc/autotest/pushdtest/pushcli"
)

/*
### 验证提连逻辑
*/
func (cli *Client) Kick() error {
	p1, err := pushcli.NewPushdcli(cli.etcdAddrs, cli.service, pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	// ConndAddress
	p1.Connect(cli.addrs[0])
	p2, err := pushcli.NewPushdcli(cli.etcdAddrs, cli.service, pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	p2.Connect(cli.addrs[1])
	cli.srvs[0].DisconnectResp()
	p2.Disconnect(cli.addrs[1])

	p1, err = pushcli.NewPushdcli(cli.etcdAddrs, cli.service, pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	if cli.srvs[1].DisLen() != 0 {
		return fmt.Errorf("kick failed")
	}
	return nil
}

func (cli *Client) KickClean() error {
	p1, err := pushcli.NewPushdcli(cli.etcdAddrs, cli.service)
	if err != nil {
		return err
	}
	// ConndAddress
	p1.Connect(cli.addrs[0])
	p2, err := pushcli.NewPushdcli(cli.etcdAddrs, cli.service)
	if err != nil {
		return err
	}
	p2.Connect(cli.addrs[1])
	cli.srvs[0].DisconnectResp()
	p2.Disconnect(cli.addrs[1])

	p1, err = pushcli.NewPushdcli(cli.etcdAddrs, cli.service, pushcli.SetCleanSession())
	if err != nil {
		return err
	}
	if cli.srvs[1].DisLen() != 0 {
		return fmt.Errorf("kick failed")
	}
	return nil
}
