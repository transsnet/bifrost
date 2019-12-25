package callback

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/meitu/bifrost/commons/efunc"
	lb "github.com/meitu/bifrost/commons/grpc-lb"
	pb "github.com/meitu/bifrost/grpc/callback"
	"google.golang.org/grpc"
)

type Call struct {
	recv    chan int
	cc      clientv3.Config
	service string
	srv     *grpc.Server
	cli     *efunc.Client
}

func (f *Call) OnConnect(ctx context.Context, req *pb.OnConnectRequest) (*pb.OnConnectReply, error) {
	f.recv <- 1
	return &pb.OnConnectReply{Cookie: []byte("connect")}, nil
}
func (f *Call) OnDisconnect(ctx context.Context, req *pb.OnDisconnectRequest) (*pb.OnDisconnectReply, error) {
	f.recv <- 1
	return &pb.OnDisconnectReply{}, nil
}

func (f *Call) OnSubscribe(ctx context.Context, req *pb.OnSubscribeRequest) (*pb.OnSubscribeReply, error) {
	f.recv <- 1
	return &pb.OnSubscribeReply{Cookie: []byte("subscribe")}, nil
}

func (f *Call) PostSubscribe(ctx context.Context, req *pb.PostSubscribeRequest) (*pb.PostSubscribeReply, error) {
	f.recv <- 1
	return &pb.PostSubscribeReply{Cookie: []byte("postsubscribe")}, nil
}

func (f *Call) OnUnsubscribe(ctx context.Context, req *pb.OnUnsubscribeRequest) (*pb.OnUnsubscribeReply, error) {
	f.recv <- 1
	return &pb.OnUnsubscribeReply{Cookie: []byte("unsubscribe")}, nil
}

func (f *Call) OnPublish(ctx context.Context, req *pb.OnPublishRequest) (*pb.OnPublishReply, error) {
	f.recv <- 1
	return &pb.OnPublishReply{Cookie: []byte("publish"), Skip: false}, nil
}

func (f *Call) OnOffline(ctx context.Context, req *pb.OnOfflineRequest) (*pb.OnOfflineReply, error) {
	f.recv <- 1
	return &pb.OnOfflineReply{Cookie: []byte("offline")}, nil
}

func (f *Call) OnACK(ctx context.Context, req *pb.OnACKRequest) (*pb.OnACKReply, error) {
	f.recv <- 1
	return &pb.OnACKReply{Cookie: []byte("postreceiveack")}, nil
}

type CallRegisterMsg struct {
	Service string
	Fnames  string
	Name    string
	Group   string
	Region  string
}

func NewCallbackServer(etcds []string, service, servname, fnames string) (*Call, error) {
	// start server
	cc := clientv3.Config{
		Endpoints: etcds,
	}

	value := &efunc.ClientData{
		ServiceName: servname,
		Funcs:       strings.Split(fnames, ","),
	}
	ecli, err := efunc.NewClient(cc, service, value)
	if err != nil {
		return nil, err
	}

	ss := grpc.NewServer()

	cli := &Call{
		recv:    make(chan int, 1000),
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

func (f *Call) Start(addr string) error {
	register, err := lb.NewNode(f.cc, f.service, addr)
	if err != nil {
		return err
	}
	if err := register.Register(); err != nil {
		return err
	}
	defer register.Deregister()

	if err := f.cli.Register(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return f.srv.Serve(lis)
}

func (f *Call) Stop() {
	f.cli.Deregister()
	f.srv.GracefulStop()
}

func (f *Call) Recv() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	select {
	case <-ctx.Done():
		return 0
	case msg := <-f.recv:
		return msg
	}
	return 0
}

func (f *Call) RecvLen() int {
	return len(f.recv)
}
