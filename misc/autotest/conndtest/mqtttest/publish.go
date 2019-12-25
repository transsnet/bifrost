package mqtttest

import (
	"context"
	"errors"

	pub "github.com/meitu/bifrost/grpc/publish"
	"google.golang.org/grpc"
)

var (
	ErrPublishResult = errors.New("result failed")
)

type PublishClient struct {
	pub.PublishServiceClient
	payload []byte
}

func NewPublishClient(pubs string, payload []byte) (*PublishClient, error) {
	pubc := &PublishClient{
		payload: payload,
	}
	cli, err := grpc.Dial(pubs, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	pubc.PublishServiceClient = pub.NewPublishServiceClient(cli)
	return pubc, nil
}

func (pubc *PublishClient) Request(tars []*pub.Target) error {
	ctx := context.Background()
	req := &pub.PublishRequest{
		MessageID: []byte("123456"),
		Targets:   tars,
		Payload:   []byte(pubc.payload),
		Weight:    0,
		StatLabel: "statlabel",
		AppKey:    "service-1544516816-1-bb7c77b22be1ff945d3272",
		// ttl:       time.Second,
	}
	resp, err := pubc.Publish(ctx, req)
	if err != nil {
		return err
	}
	for _, result := range resp.Results {
		if result == pub.ErrCode_ErrInternalError {
			return ErrPublishResult
		}
	}
	return nil
}

func (pubc *PublishClient) Close() {
	pubc.Close()
}

func (pubc *PublishClient) Payload() []byte {
	return pubc.payload
}

func (pubc *PublishClient) SetPayload(payload []byte) {
	pubc.payload = payload
}
