package fakePush

import (
	"context"
	"log"
	"net"

	"github.com/coreos/etcd/clientv3"
	lb "github.com/meitu/bifrost/commons/grpc-lb"
	pb "github.com/meitu/bifrost/grpc/push"
	"google.golang.org/grpc"
)

type fakeServer struct {
	r   *lb.Node
	srv *grpc.Server
}

func (s *fakeServer) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectResp, error) {
	log.Println("Connect :", req)
	//time.Sleep(time.Duration(delay) * time.Millisecond)
	return &pb.ConnectResp{}, nil
}
func (s *fakeServer) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectResp, error) {
	//time.Sleep(time.Duration(delay) * time.Millisecond)
	log.Println("Disconnect :", req)
	return &pb.DisconnectResp{}, nil
}

func (s *fakeServer) Subscribe(ctx context.Context, req *pb.SubscribeReq) (*pb.SubscribeResp, error) {
	log.Println("Subscribe :", req)
	return &pb.SubscribeResp{}, nil
}

func (s *fakeServer) Unsubscribe(ctx context.Context, req *pb.UnsubscribeReq) (*pb.UnsubscribeResp, error) {
	log.Println("Unsubscribe :", req)
	return &pb.UnsubscribeResp{}, nil
}

func (s *fakeServer) MQTTPublish(ctx context.Context, req *pb.PublishReq) (*pb.PublishResp, error) {
	log.Println("MQTTPublish :", req)
	return &pb.PublishResp{}, nil
}

func (s *fakeServer) Pubrec(ctx context.Context, req *pb.PubrecReq) (*pb.PubrecResp, error) {
	log.Println("Pubrec :", req)
	return &pb.PubrecResp{}, nil
}

func (s *fakeServer) Pubrel(ctx context.Context, req *pb.PubrelReq) (*pb.PubrelResp, error) {
	log.Println("Pubrel :", req)
	return &pb.PubrelResp{}, nil
}

func (s *fakeServer) Pubcomp(ctx context.Context, req *pb.PubcompReq) (*pb.PubcompResp, error) {
	log.Println("Pubcomp :", req)
	return &pb.PubcompResp{}, nil
}

func (s *fakeServer) RangeUnack(ctx context.Context, req *pb.RangeUnackReq) (*pb.RangeUnackResp, error) {
	log.Println("RangeUnack :", req)
	return &pb.RangeUnackResp{}, nil
}

func (s *fakeServer) PutUnack(ctx context.Context, req *pb.PutUnackReq) (*pb.PutUnackResp, error) {
	log.Println("PubtUnack :", req)
	return &pb.PutUnackResp{}, nil
}

func (s *fakeServer) DelUnack(ctx context.Context, req *pb.DelUnackReq) (*pb.DelUnackResp, error) {
	log.Println("DelUnack :", req)
	return &pb.DelUnackResp{}, nil
}

func (s *fakeServer) Pull(ctx context.Context, req *pb.PullReq) (*pb.PullResp, error) {
	log.Println("Pull :", req)
	return &pb.PullResp{}, nil
}

func (s *fakeServer) PostSubscribe(ctx context.Context, req *pb.PostSubscribeReq) (*pb.PostSubscribeResp, error) {
	log.Println("PostSubscribe :", req)
	return &pb.PostSubscribeResp{}, nil
}

func (s *fakeServer) AddRoute(ctx context.Context, req *pb.AddRouteReq) (*pb.AddRouteResp, error) {
	log.Println("AddRoute :", req)
	return &pb.AddRouteResp{}, nil
}

func (s *fakeServer) RemoveRoute(ctx context.Context, req *pb.RemoveRouteReq) (*pb.RemoveRouteResp, error) {
	log.Println("RemoveRouete :", req)
	return &pb.RemoveRouteResp{}, nil
}

func (s *fakeServer) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *fakeServer) Close() {
	s.r.Deregister()
	s.srv.GracefulStop()
}

var gss []*fakeServer

func DefaultPushdServer(service, addr, etcda, group, region string) error {
	cc := clientv3.Config{
		Endpoints: []string{etcda},
	}
	r, err := lb.NewNode(cc, service, addr)
	if err != nil {
		return err
	}
	r.Register()
	ss := grpc.NewServer()

	fs := &fakeServer{
		srv: ss,
		r:   r,
	}

	pb.RegisterPushServiceServer(ss, fs)

	gss = append(gss, fs)

	if err := fs.ListenAndServe(addr); err != nil {
		return err
	}
	return nil
}

func Close() {
	for _, gs := range gss {
		gs.Close()
	}
}
