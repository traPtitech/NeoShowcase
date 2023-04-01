package domain

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb/pbconnect"
)

type ComponentService interface {
	pbconnect.ComponentServiceHandler
	BroadcastBuilder(req *pb.BuilderRequest)
	BroadcastSSGen(req *pb.SSGenRequest)
}

type ComponentServiceClient interface {
	ConnectBuilder(ctx context.Context, onRequest func(req *pb.BuilderRequest), response <-chan *pb.BuilderResponse) error
	ConnectSSGen(ctx context.Context, onRequest func(req *pb.SSGenRequest)) error
}
