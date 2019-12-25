package scene

import (
	"context"
	"errors"
	"log"

	pub "github.com/meitu/bifrost/grpc/publish"
	"google.golang.org/grpc"
)

type Pubcli struct {
	addr   string
	appkey string
	cli    pub.PublishServiceClient
}

func NewPubcli(addr, appkey string) *Pubcli {
	pubc := &Pubcli{
		addr:   addr,
		appkey: appkey,
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("connect pushd fail %s", err)
		return nil
	}
	pubc.cli = pub.NewPublishServiceClient(conn)
	return pubc
}

func (pubc *Pubcli) Request(payload []byte, tars []*pub.Target) error {
	ctx := context.Background()
	req := &pub.PublishRequest{
		MessageID: []byte("publish"),
		Targets:   tars,
		Payload:   payload,
		Weight:    0,
		AppKey:    pubc.appkey,
	}
	resp, err := pubc.cli.Publish(ctx, req)
	if err != nil {
		return err
	}
	for _, r := range resp.Results {
		if r == pub.ErrCode_ErrInternalError {
			return errors.New("publish failed")
		}
	}
	return nil
}
