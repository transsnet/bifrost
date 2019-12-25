package testcase

import "github.com/meitu/bifrost/misc/autotest/pushdtest/conndsrv"

type Client struct {
	// *pushcli.Pushdcli
	etcdAddrs []string
	service   string
	group     string
	srvs      []*conndsrv.FakeServer
	addrs     []string
}

//srvs 默认只使用srvs中的第一个进行消息的注册分发
//im test 系统要求至少两个
func NewClient(etcdAddrs, addrs []string, service, group string) (*Client, error) {
	var srvs []*conndsrv.FakeServer
	for _, addr := range addrs {
		srv, err := conndsrv.NewConndSrv(addr)
		if err != nil {
			return nil, err
		}
		srvs = append(srvs, srv)
	}
	return &Client{
		etcdAddrs: etcdAddrs,
		service:   service,
		group:     group,
		srvs:      srvs,
		addrs:     addrs,
	}, nil
}
