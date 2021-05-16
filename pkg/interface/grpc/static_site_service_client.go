package grpc

import (
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"google.golang.org/grpc"
)

type StaticSiteServiceClientConn struct {
	*grpc.ClientConn
}

type StaticSiteServiceClientConfig struct {
	Insecure bool   `mapstructure:"insecure" yaml:"insecure"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
}

func (c *StaticSiteServiceClientConfig) provideClientConfig() ClientConfig {
	return ClientConfig{
		Insecure: c.Insecure,
		Addr:     c.Addr,
	}
}

func NewStaticSiteServiceClientConn(c StaticSiteServiceClientConfig) (*StaticSiteServiceClientConn, error) {
	conn, err := NewClient(c.provideClientConfig())
	if err != nil {
		return nil, err
	}
	return &StaticSiteServiceClientConn{ClientConn: conn}, nil
}

func NewStaticSiteServiceClient(cc *StaticSiteServiceClientConn) pb.StaticSiteServiceClient {
	return pb.NewStaticSiteServiceClient(cc.ClientConn)
}
