package staticsitegen

import (
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/generator"
)

type Config struct {
	ServerType  string
	NginxConfig generator.Nginx
}
