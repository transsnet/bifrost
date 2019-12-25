package fakeConnd

import (
	"context"
	"log"

	"github.com/facebookgo/grace/gracenet"
	"google.golang.org/grpc"

	pb "github.com/meitu/bifrost/grpc/conn"
)

type fakeServer struct {
}

func (s *fakeServer) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectResp, error) {
	log.Printf("Disconnect %#v\n", req)
	return &pb.DisconnectResp{}, nil
}

func (s *fakeServer) Notify(ctx context.Context, req *pb.NotifyReq) (*pb.NotifyResp, error) {
	log.Printf("Notify %#v\n", req)
	return &pb.NotifyResp{}, nil
}

var server fakeServer

func DefaultConndServer(etcd []string, service, address string) {
	gnet := &gracenet.Net{}
	lis, err := gnet.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	gs := grpc.NewServer()
	pb.RegisterConnServiceServer(gs, &server)
	if err := gs.Serve(lis); err != nil {
		// log.Fatal(err)
	}
}
