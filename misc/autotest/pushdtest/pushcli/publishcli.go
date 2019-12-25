package pushcli

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	lb "github.com/meitu/bifrost/commons/grpc-lb"
	pb "github.com/meitu/bifrost/grpc/publish"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type Publishcli struct {
	instance pb.PublishServiceClient
}

func NewPublishcli(etcdAddrs []string, service string) (*Publishcli, error) {
	cc := clientv3.Config{
		Endpoints: etcdAddrs,
	}
	r := lb.NewResolver(cc, service)
	resolver.Register(r)

	var ropts []grpc.DialOption
	ropts = append(ropts, grpc.WithBalancerName("round_robin"))
	ropts = append(ropts, grpc.WithInsecure())
	ropts = append(ropts, grpc.WithBlock())
	conn, err := grpc.Dial(r.URL(), ropts...)
	if err != nil {
		return nil, err
	}

	return &Publishcli{
		instance: pb.NewPublishServiceClient(conn),
	}, nil
}

func (pub *Publishcli) Publish(topic string, payload []byte, retain bool, qos int32) (*pb.PublishReply, error) {
	req := &pb.PublishRequest{
		Payload: payload,
		AppKey:  "service-1544516816-1-bb7c77b22be1ff945d3272",
		Targets: []*pb.Target{
			&pb.Target{
				Topic:    topic,
				Qos:      qos,
				IsRetain: retain,
			},
		},
	}
	return pub.instance.Publish(context.TODO(), req)
}
