package domain

import (
	"context"
)

type (
	StaticServerDocumentRootPath string
	StaticServerPort             int
)

type StaticServerConnectivityConfig struct {
	Service struct {
		Namespace string `mapstructure:"namespace" yaml:"namespace"`
		Kind      string `mapstructure:"kind" yaml:"kind"`
		Name      string `mapstructure:"name" yaml:"name"`
		Port      int    `mapstructure:"port" yaml:"port"`
	} `mapstructure:"service" yaml:"service"`
	URL string `mapstructure:"url" yaml:"url"`
}

type SSEngine interface {
	Start(ctx context.Context) error
	Reconcile(sites []*StaticSite) error
	Shutdown(ctx context.Context) error
}

type StaticSite struct {
	Application *Application
	Website     *Website
	ArtifactID  string
}
