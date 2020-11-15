package container

import (
	"context"
)

type Manager interface {
	Create(ctx context.Context, args CreateArgs) (*CreateResult, error)
	// TODO Start(ctx context.Context, )
	// TODO Stop(ctx context.Context, )
	// TODO Restart(ctx context.Context, )
	Destroy(ctx context.Context, args DestroyArgs) (*DestroyResult, error)
	List(ctx context.Context) (*ListResult, error)
}

type CreateArgs struct {
	ApplicationID string
	ImageName     string
	ImageTag      string
	Labels        map[string]string

	Domain   string
	HTTPPort int

	NoStart bool
}

type CreateResult struct {
}

type DestroyArgs struct {
	ApplicationID string
}

type DestroyResult struct {
}

type ListResult struct {
	Containers []Container
}

type Container struct {
	ApplicationID string
	State         string
}
