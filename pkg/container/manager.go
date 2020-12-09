package container

import (
	"context"
)

type Manager interface {
	Create(ctx context.Context, args CreateArgs) (*CreateResult, error)
	// TODO Start(ctx context.Context, ) k8sではCreateと同じ
	Stop(ctx context.Context, args StopArgs) (*StopResult, error)          // k8sではDestroyと同じ
	Restart(ctx context.Context, args RestartArgs) (*RestartResult, error) // k8sではCreateと同じ
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
