package grpc

import (
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"google.golang.org/grpc"
)

type BuilderServiceClientConn struct {
	*grpc.ClientConn
}

type BuilderServiceClientConfig struct {
	Insecure bool   `mapstructure:"insecure" yaml:"insecure"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
}

func (c *BuilderServiceClientConfig) provideClientConfig() ClientConfig {
	return ClientConfig{
		Insecure: c.Insecure,
		Addr:     c.Addr,
	}
}

func NewBuilderServiceClientConn(c BuilderServiceClientConfig) (*BuilderServiceClientConn, error) {
	conn, err := NewClient(c.provideClientConfig())
	if err != nil {
		return nil, err
	}
	return &BuilderServiceClientConn{ClientConn: conn}, nil
}

func NewBuilderServiceClient(cc *BuilderServiceClientConn) pb.BuilderServiceClient {
	return pb.NewBuilderServiceClient(cc.ClientConn)
}
