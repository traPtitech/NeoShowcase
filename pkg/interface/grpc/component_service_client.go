package grpc

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
)

type ComponentServiceClientConn struct {
	*grpc.ClientConn
}

type ComponentServiceClientConfig struct {
	Insecure bool   `mapstructure:"insecure" yaml:"insecure"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
}

func (c *ComponentServiceClientConfig) provideClientConfig() ClientConfig {
	return ClientConfig{
		Insecure: c.Insecure,
		Addr:     c.Addr,
	}
}

func NewComponentServiceClientConn(c ComponentServiceClientConfig) (*ComponentServiceClientConn, error) {
	conn, err := NewClient(c.provideClientConfig())
	if err != nil {
		return nil, err
	}
	return &ComponentServiceClientConn{ClientConn: conn}, nil
}

type ComponentServiceClient struct {
	client pb.ComponentServiceClient
}

func NewComponentServiceClient(cc *ComponentServiceClientConn) domain.ComponentServiceClient {
	return &ComponentServiceClient{
		client: pb.NewComponentServiceClient(cc.ClientConn),
	}
}

func (c *ComponentServiceClient) ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := c.client.ConnectBuilder(ctx)
	if err != nil {
		return err
	}

	go func() {
		defer cancel()
		defer func() {
			err := conn.CloseSend()
			if err != nil {
				log.WithError(err).Error("failed to close send stream")
			}
		}()

		for {
			select {
			case res, ok := <-response:
				if !ok {
					return
				}
				err := conn.Send(res)
				if err != nil {
					log.WithError(err).Error("failed to send builder response")
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		onRequest(req)
	}
}

func (c *ComponentServiceClient) ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := c.client.ConnectSSGen(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		cancel()
		err := conn.CloseSend()
		if err != nil {
			log.WithError(err).Error("failed to close send stream")
		}
	}()

	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		onRequest(req)
	}
}
