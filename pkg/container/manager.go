package container

import (
	"context"
)

type Manager interface {
	Create(ctx context.Context, args CreateArgs) (*CreateResult, error)
	Start(ctx context.Context, args StartArgs) (*StartResult, error)
	Stop(ctx context.Context, args StopArgs) (*StopResult, error)
	Restart(ctx context.Context, args RestartArgs) (*RestartResult, error)
	Destroy(ctx context.Context, args DestroyArgs) (*DestroyResult, error)
	List(ctx context.Context) (*ListResult, error)
	Dispose(ctx context.Context) error
}

type CreateArgs struct {
	ApplicationID string
	ImageName     string
	ImageTag      string
	Labels        map[string]string
	HTTPProxy     *HTTPProxy
	NoStart       bool
}

type HTTPProxy struct {
	Domain string
	Port   int
}

type CreateResult struct {
}

type StartArgs struct {
	ApplicationID string
}

type StartResult struct {
}

type DestroyArgs struct {
	ApplicationID string
}

type DestroyResult struct {
}

type StopArgs struct {
	ApplicationID string
}

type StopResult struct {
}

type RestartArgs struct {
	ApplicationID string
}

type RestartResult struct {
}

type ListResult struct {
	Containers []Container
}

type Container struct {
	ApplicationID string
	State         string
}
