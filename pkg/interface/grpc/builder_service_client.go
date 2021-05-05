package grpc

import (
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"google.golang.org/grpc"
)

type BuilderServiceClientConn struct {
	*grpc.ClientConn
}

type BuilderServiceClientConfig struct {
	ClientConfig
}

func NewBuilderServiceClientConn(c BuilderServiceClientConfig) (*BuilderServiceClientConn, error) {
	conn, err := NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &BuilderServiceClientConn{ClientConn: conn}, nil
}

func NewBuilderServiceClient(cc *BuilderServiceClientConn) pb.BuilderServiceClient {
	return pb.NewBuilderServiceClient(cc.ClientConn)
}
