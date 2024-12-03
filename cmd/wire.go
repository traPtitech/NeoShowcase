//go:build wireinject
// +build wireinject

package main

import (
	"github.com/friendsofgo/errors"
	"github.com/google/wire"

	authdev "github.com/traPtitech/neoshowcase/cmd/auth-dev"
	"github.com/traPtitech/neoshowcase/cmd/builder"
	buildpackhelper "github.com/traPtitech/neoshowcase/cmd/buildpack-helper"
	"github.com/traPtitech/neoshowcase/cmd/controller"
	"github.com/traPtitech/neoshowcase/cmd/gateway"
	giteaintegration "github.com/traPtitech/neoshowcase/cmd/gitea-integration"
	"github.com/traPtitech/neoshowcase/cmd/ssgen"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/k8simpl"
	ubuilder "github.com/traPtitech/neoshowcase/pkg/usecase/builder"
	buildermock "github.com/traPtitech/neoshowcase/pkg/usecase/builder/mock"
)

func NewAuthDev(c Config) (component, error) {
	wire.Build(
		providers,
		wire.Bind(new(component), new(*authdev.Server)),
	)
	return nil, nil
}

func NewBuilder(c Config) (component, error) {
	if c.Components.Builder.Mock {
		return newMockBuilder(c)
	} else {
		return newBuilder(c)
	}
}

func newBuilder(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(BuilderConfig), "Buildpack"),
		wire.Bind(new(ubuilder.Service), new(*ubuilder.ServiceImpl)),
		wire.Bind(new(component), new(*builder.Server)),
		wire.Struct(new(builder.Server), "*"),
	)
	return nil, nil
}

func newMockBuilder(c Config) (component, error) {
	wire.Build(
		providers,
		buildermock.NewBuilderServiceMock,
		wire.Bind(new(ubuilder.Service), new(*buildermock.BuilderServiceMock)),
		wire.Bind(new(component), new(*builder.Server)),
		wire.Struct(new(builder.Server), "*"),
	)
	return nil, nil
}

func NewBuildpackHelper(c Config) (component, error) {
	wire.Build(
		providers,
		wire.Bind(new(component), new(*buildpackhelper.Server)),
		wire.Struct(new(buildpackhelper.Server), "*"),
	)
	return nil, nil
}

func NewController(c Config) (component, error) {
	switch c.Components.Controller.Mode {
	case "docker":
		return NewControllerDocker(c)
	case "k8s", "kubernetes":
		return NewControllerK8s(c)
	}
	return nil, errors.New("unknown mode: " + c.Components.Controller.Mode)
}

func NewControllerDocker(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(ControllerConfig), "Docker", "SSH", "Webhook"),
		wire.Bind(new(domain.Backend), new(*dockerimpl.Backend)),
		wire.Bind(new(component), new(*controller.Server)),
		wire.Struct(new(controller.Server), "*"),
	)
	return nil, nil
}

func NewControllerK8s(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(ControllerConfig), "K8s", "SSH", "Webhook"),
		wire.Bind(new(domain.Backend), new(*k8simpl.Backend)),
		wire.Bind(new(component), new(*controller.Server)),
		wire.Struct(new(controller.Server), "*"),
	)
	return nil, nil
}

func NewGateway(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(GatewayConfig), "AvatarBaseURL", "AuthHeader", "Controller", "MariaDB", "MongoDB"),
		wire.Bind(new(component), new(*gateway.Server)),
		wire.Struct(new(gateway.Server), "*"),
	)
	return nil, nil
}

func NewGiteaIntegration(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(GiteaIntegrationConfig), "Controller"),
		wire.Bind(new(component), new(*giteaintegration.Server)),
		wire.Struct(new(giteaintegration.Server), "*"),
	)
	return nil, nil
}

func NewSSGen(c Config) (component, error) {
	wire.Build(
		providers,
		wire.FieldsOf(new(SSGenConfig), "HealthPort", "Controller"),
		wire.Bind(new(component), new(*ssgen.Server)),
		wire.Struct(new(ssgen.Server), "*"),
	)
	return nil, nil
}
