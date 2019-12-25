package callback

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/meitu/bifrost/commons/efunc"
	lb "github.com/meitu/bifrost/commons/grpc-lb"
	pb "github.com/meitu/bifrost/grpc/callback"
	"google.golang.org/grpc"
)

type ConndCookie struct {
	Discovery bool
	cc        clientv3.Config
	service   string
	srv       *grpc.Server
	cli       *efunc.Client
}

func (f *ConndCookie) OnConnect(ctx context.Context, req *pb.OnConnectRequest) (*pb.OnConnectReply, error) {
	f.Discovery = true
	return &pb.OnConnectReply{Cookie: []byte("connect")}, nil
}
func (f *ConndCookie) OnDisconnect(ctx context.Context, req *pb.OnDisconnectRequest) (*pb.OnDisconnectReply, error) {
	if !bytes.Equal(req.Cookie, []byte("unsubscribe")) {
		log.Fatal("ondisconnect cookie is not equal")
	}
	return &pb.OnDisconnectReply{}, nil
}

func (f *ConndCookie) OnSubscribe(ctx context.Context, req *pb.OnSubscribeRequest) (*pb.OnSubscribeReply, error) {
	if !bytes.Equal(req.Cookie, []byte("connect")) {
		log.Fatal("onsubscribe cookie is not equal")
	}
	return &pb.OnSubscribeReply{Cookie: []byte("subscribe")}, nil
}

func (f *ConndCookie) PostSubscribe(ctx context.Context, req *pb.PostSubscribeRequest) (*pb.PostSubscribeReply, error) {
	if !bytes.Equal(req.Cookie, []byte("subscribe")) {
		log.Fatal("postsubscribe cookie is not equal")
	}
	return &pb.PostSubscribeReply{Cookie: []byte("postsubscribe")}, nil
}

func (f *ConndCookie) OnUnsubscribe(ctx context.Context, req *pb.OnUnsubscribeRequest) (*pb.OnUnsubscribeReply, error) {
	if !bytes.Equal(req.Cookie, []byte("publish")) {
		log.Fatal("onUnsubscribe cookie is not equal")
	}
	return &pb.OnUnsubscribeReply{Cookie: []byte("unsubscribe")}, nil
}

func (f *ConndCookie) OnPublish(ctx context.Context, req *pb.OnPublishRequest) (*pb.OnPublishReply, error) {
	if !bytes.Equal(req.Cookie, []byte("postsubscribe")) {
		log.Fatal("onpublish cookie is not equal")
	}
	return &pb.OnPublishReply{Cookie: []byte("publish"), Skip: true}, nil
}

func (f *ConndCookie) OnOffline(ctx context.Context, req *pb.OnOfflineRequest) (*pb.OnOfflineReply, error) {
	return &pb.OnOfflineReply{Cookie: []byte("offilne")}, nil
}

func (f *ConndCookie) OnACK(ctx context.Context, req *pb.OnACKRequest) (*pb.OnACKReply, error) {
	if !bytes.Equal(req.Cookie, []byte("publish")) {
		log.Fatal("postsubscribe cookie is not equal")
	}
	return &pb.OnACKReply{Cookie: []byte("postreceiveack")}, nil
}

func NewConndCookieServer(etcds []string, service, servname string) (*ConndCookie, error) {
	cc := clientv3.Config{
		Endpoints: etcds,
	}

	fnames := "PostSubscribe,OnDisconnect,OnSubscribe,OnUnsubscribe,OnACK,OnPublish,OnOffline,OnConnect"
	value := &efunc.ClientData{
		ServiceName: servname,
		Funcs:       strings.Split(fnames, ","),
	}
	ecli, err := efunc.NewClient(cc, service, value)
	if err != nil {
		return nil, err
	}

	ss := grpc.NewServer()

	cli := &ConndCookie{
		srv:     ss,
		cli:     ecli,
		service: servname,
	}

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

func (fs *ConndCookie) Start(addr string) error {
	// return fs.gs.ListenAndServe(addr)
	register, err := lb.NewNode(fs.cc, fs.service, addr)
	if err != nil {
		return err
	}
	if err := register.Register(); err != nil {
		return err
	}
	defer register.Deregister()

	fs.cli.Register()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return fs.srv.Serve(lis)
}

func (fs *ConndCookie) Stop() {
	fs.cli.Deregister()
	fs.srv.GracefulStop()
}
