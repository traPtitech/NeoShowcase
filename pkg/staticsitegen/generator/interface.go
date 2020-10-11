package generator

import (
	"github.com/traPtitech/neoshowcase/pkg/apiserver/api"
)

type Engine interface {
	Generate(sites []api.StaticSite) error
}
