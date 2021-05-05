package grpc

import (
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"google.golang.org/grpc"
)

type StaticSiteServiceClientConn struct {
	*grpc.ClientConn
}

type StaticSiteServiceClientConfig struct {
	ClientConfig
}

func NewStaticSiteServiceClientConn(c StaticSiteServiceClientConfig) (*StaticSiteServiceClientConn, error) {
	conn, err := NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &StaticSiteServiceClientConn{ClientConn: conn}, nil
}

func NewStaticSiteServiceClient(cc *StaticSiteServiceClientConn) pb.StaticSiteServiceClient {
	return pb.NewStaticSiteServiceClient(cc.ClientConn)
}
