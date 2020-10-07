package generator

import (
	"github.com/traPtitech/neoshowcase/pkg/apiserver/grpc/api"
)

type Engine interface {
	Generate(sites []api.StaticSite) error
}
