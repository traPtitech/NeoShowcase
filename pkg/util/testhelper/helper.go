package testhelper

import "go.uber.org/dig"

// Container is a wrapper of dig.Container
type Container struct {
	*dig.Container
}

func (c *Container) Provide(constructor any, opts ...dig.ProvideOption) {
	if err := c.Container.Provide(constructor, opts...); err != nil {
		panic(err)
	}
}

type ContainerOption func(*Container)

func NewContainer(opts ...ContainerOption) *Container {
	c := &Container{Container: dig.New()}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func Resolve[T any](c *Container) T {
	var v T
	if err := c.Container.Invoke(func(res T) {
		v = res
	}); err != nil {
		panic(err)
	}
	return v
}
