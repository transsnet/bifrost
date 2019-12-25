package callback

import (
	"context"
	"net"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/meitu/bifrost/commons/efunc"
	lb "github.com/meitu/bifrost/commons/grpc-lb"
	pb "github.com/meitu/bifrost/grpc/callback"
	"google.golang.org/grpc"
)

type PushdCookie struct {
	cc      clientv3.Config
	service string
	srv     *grpc.Server
	cli     *efunc.Client
}

func (f *PushdCookie) OnConnect(ctx context.Context, req *pb.OnConnectRequest) (*pb.OnConnectReply, error) {
	return &pb.OnConnectReply{Cookie: []byte("connect")}, nil
}
func (f *PushdCookie) OnDisconnect(ctx context.Context, req *pb.OnDisconnectRequest) (*pb.OnDisconnectReply, error) {
	return &pb.OnDisconnectReply{}, nil
}

func (f *PushdCookie) OnSubscribe(ctx context.Context, req *pb.OnSubscribeRequest) (*pb.OnSubscribeReply, error) {
	return &pb.OnSubscribeReply{Cookie: []byte("subscribe")}, nil
}

func (f *PushdCookie) PostSubscribe(ctx context.Context, req *pb.PostSubscribeRequest) (*pb.PostSubscribeReply, error) {
	return &pb.PostSubscribeReply{Cookie: []byte("postsubscribe")}, nil
}

func (f *PushdCookie) OnUnsubscribe(ctx context.Context, req *pb.OnUnsubscribeRequest) (*pb.OnUnsubscribeReply, error) {
	return &pb.OnUnsubscribeReply{Cookie: []byte("unsubscribe")}, nil
}

func (f *PushdCookie) OnPublish(ctx context.Context, req *pb.OnPublishRequest) (*pb.OnPublishReply, error) {
	return &pb.OnPublishReply{Cookie: []byte("publish"), Skip: true}, nil
}

func (f *PushdCookie) OnOffline(ctx context.Context, req *pb.OnOfflineRequest) (*pb.OnOfflineReply, error) {
	return &pb.OnOfflineReply{Cookie: []byte("offilne")}, nil
}

func (f *PushdCookie) OnACK(ctx context.Context, req *pb.OnACKRequest) (*pb.OnACKReply, error) {
	return &pb.OnACKReply{Cookie: []byte("postreceiveack")}, nil
}

func NewPushdCookieServer(etcds []string, service, group, servname string) (*PushdCookie, error) {
	// start server
	cc := clientv3.Config{
		Endpoints: etcds,
	}

	fnames := "PostSubscribe,OnDisconnect,OnSubscribe,OnUnsubscribe,OnACK,OnPublish,OnOffline"
	value := &efunc.ClientData{
		ServiceName: servname,
		Funcs:       strings.Split(fnames, ","),
	}
	ecli, err := efunc.NewClient(cc, service, value)
	if err != nil {
		return nil, err
	}

	ss := grpc.NewServer()

	cli := &PushdCookie{
		cc:      cc,
		srv:     ss,
		service: servname,
		cli:     ecli,
	}
	//register data
	pb.RegisterOnConnectServer(ss, cli)
	pb.RegisterOnDisconnectServer(ss, cli)
	pb.RegisterOnOfflineServer(ss, cli)
	pb.RegisterOnPublishServer(ss, cli)
	pb.RegisterOnSubscribeServer(ss, cli)
	pb.RegisterOnUnsubscribeServer(ss, cli)
	pb.RegisterPostSubscribeServer(ss, cli)
	pb.RegisterOnACKServer(ss, cli)
	return cli, nil
}

func (f *PushdCookie) Start(addr string) error {
	register, err := lb.NewNode(f.cc, f.service, addr)
	if err != nil {
		return err
	}
	if err := register.Register(); err != nil {
		return err
	}
	defer register.Deregister()

	f.cli.Register()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return f.srv.Serve(lis)
}

func (f *PushdCookie) Stop() {
	f.cli.Deregister()
	f.srv.GracefulStop()
}
