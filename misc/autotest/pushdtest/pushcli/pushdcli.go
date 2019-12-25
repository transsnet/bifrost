package pushcli

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	lb "github.com/meitu/bifrost/commons/grpc-lb"
	pb "github.com/meitu/bifrost/grpc/push"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type Pushdcli struct {
	// clientid     string
	// cleansession bool
	// service      string
	push pb.PushServiceClient
	*config
}

//Pushdcli 创建一个pushd cli
func NewPushdcli(etcdAddrs []string, service string, opts ...opt) (*Pushdcli, error) {
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

	push := pb.NewPushServiceClient(conn)
	c := &config{}
	for _, o := range opts {
		o(c)
	}

	if c.clientid == "" {
		c.clientid = "clientid"
	}
	if c.service == "" {
		c.service = "service"
	}
	if len(c.payload) == 0 {
		c.payload = []byte("payload")
	}

	return &Pushdcli{
		config: c,
		push:   push,
	}, nil
}

func (p *Pushdcli) Connect(addr string) (*pb.ConnectResp, error) {
	req := &pb.ConnectReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
		GrpcAddress:  addr,
	}
	return p.push.Connect(context.Background(), req)
}

func (p *Pushdcli) RangeUnack(offset []byte, limit int64) (*pb.RangeUnackResp, error) {
	req := &pb.RangeUnackReq{
		ClientID: p.clientid,
		Service:  p.service,
		Offset:   offset,
		Limit:    limit,
	}
	return p.push.RangeUnack(context.Background(), req)
}

func (p *Pushdcli) PutUnack(topic string, index []byte, id int64) (*pb.PutUnackResp, error) {
	req := &pb.PutUnackReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
		Messages: []*pb.UnackDesc{
			&pb.UnackDesc{
				Topic:     topic,
				Index:     index,
				MessageID: id,
			},
		},
	}
	return p.push.PutUnack(context.Background(), req)
}

func (p *Pushdcli) Disconnect(addr string) (*pb.DisconnectResp, error) {
	req := &pb.DisconnectReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
		GrpcAddress:  addr,
	}
	return p.push.Disconnect(context.Background(), req)
}

func (p *Pushdcli) Subscribe(topics []string, qoss []int32) (*pb.SubscribeResp, error) {
	req := &pb.SubscribeReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
		Topics:       topics,
		Qoss:         qoss,
	}
	return p.push.Subscribe(context.Background(), req)
}

func (p *Pushdcli) PostSubscribe() (*pb.PostSubscribeResp, error) {
	req := &pb.PostSubscribeReq{
		ClientID: p.clientid,
		Service:  p.service,
		Topics:   []string{"topic"},
		Qoss:     []int32{1},
	}
	return p.push.PostSubscribe(context.Background(), req)
}

func (p *Pushdcli) Unsubscribe(topics []string) (*pb.UnsubscribeResp, error) {
	req := &pb.UnsubscribeReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
		Topics:       topics,
	}
	return p.push.Unsubscribe(context.Background(), req)
}

func (p *Pushdcli) MQTTPublish(topic string, qos int32, retain bool) (*pb.PublishResp, error) {
	req := &pb.PublishReq{
		ClientID: p.clientid,
		Service:  p.service,
		Message: &pb.Message{
			Topic:   topic,
			Qos:     qos,
			Retain:  retain,
			Payload: p.payload,
		},
	}
	return p.push.MQTTPublish(context.Background(), req)
}

func (p *Pushdcli) Pubrec() (*pb.PubrecResp, error) {
	req := &pb.PubrecReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
	}
	return p.push.Pubrec(context.Background(), req)
}

func (p *Pushdcli) Pubrel() (*pb.PubrelResp, error) {
	req := &pb.PubrelReq{
		ClientID: p.clientid,
		Service:  p.service,
	}
	return p.push.Pubrel(context.Background(), req)
}

func (p *Pushdcli) Pubcomp() (*pb.PubcompResp, error) {
	req := &pb.PubcompReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
	}
	return p.push.Pubcomp(context.Background(), req)
}

func (p *Pushdcli) DelUnack(id int64) (*pb.DelUnackResp, error) {
	req := &pb.DelUnackReq{
		ClientID:     p.clientid,
		Service:      p.service,
		CleanSession: p.cleansession,
		MessageID:    id,
	}
	return p.push.DelUnack(context.Background(), req)
}

func (p *Pushdcli) Pull(topic string, offset []byte, limit int64) (*pb.PullResp, error) {
	req := &pb.PullReq{
		Service: p.service,
		Topic:   topic,
		Offset:  offset,
		Limit:   limit,
	}
	return p.push.Pull(context.Background(), req)
}

func (p *Pushdcli) AddRoute(topic string, addr string) (*pb.AddRouteResp, error) {
	req := &pb.AddRouteReq{
		Service:     p.service,
		Topic:       topic,
		GrpcAddress: addr,
		Version:     0,
	}
	return p.push.AddRoute(context.Background(), req)
}

func (p *Pushdcli) RemoveRoute(topic string, addr string) (*pb.RemoveRouteResp, error) {
	req := &pb.RemoveRouteReq{
		Service:     p.service,
		Topic:       topic,
		GrpcAddress: addr,
		Version:     1,
	}
	return p.push.RemoveRoute(context.Background(), req)
}
