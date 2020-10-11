package apiserver

import (
	"context"
	"github.com/traPtitech/neoshowcase/pkg/apiserver/api"
)

type PingService struct {
	api.UnimplementedPingServer
}

func (p *PingService) Ping(ctx context.Context, request *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Msg: request.GetMsg()}, nil
}
