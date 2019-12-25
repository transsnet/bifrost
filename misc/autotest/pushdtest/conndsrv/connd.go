package conndsrv

import (
	"context"
	"log"

	"github.com/facebookgo/grace/gracenet"
	pb "github.com/meitu/bifrost/grpc/conn"
	"google.golang.org/grpc"
)

type FakeServer struct {
	DisChan    chan *pb.DisconnectReq
	NotifyChan chan *pb.NotifyReq
	srv        *grpc.Server
}

func (s *FakeServer) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectResp, error) {
	s.DisChan <- req
	return &pb.DisconnectResp{}, nil
}

func (s *FakeServer) Notify(ctx context.Context, req *pb.NotifyReq) (*pb.NotifyResp, error) {
	s.NotifyChan <- req
	return &pb.NotifyResp{}, nil
}

func (s *FakeServer) NotifyResp() *pb.NotifyReq {
	return <-s.NotifyChan
}

func (s *FakeServer) NotifyRespLen() int {
	return len(s.NotifyChan)
}

func (s *FakeServer) DisconnectResp() *pb.DisconnectReq {
	return <-s.DisChan
}

func (s *FakeServer) DisLen() int {
	return len(s.DisChan)
}

func (s *FakeServer) Close() {
	s.srv.Stop()
}

func NewConndSrv(address string) (*FakeServer, error) {
	gnet := &gracenet.Net{}
	lis, err := gnet.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	srv := grpc.NewServer()

	server := &FakeServer{
		srv:        srv,
		NotifyChan: make(chan *pb.NotifyReq, 10),
		DisChan:    make(chan *pb.DisconnectReq, 10),
		// reg:        reg,
	}

	pb.RegisterConnServiceServer(srv, server)
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	return server, nil
}
