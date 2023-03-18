package builder

import (
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type BuildType int

const (
	BuildTypeRuntime BuildType = iota
	BuildTypeStatic
)

func (t BuildType) String() string {
	switch t {
	case BuildTypeRuntime:
		return models.ApplicationsBuildTypeRuntime
	case BuildTypeStatic:
		return models.ApplicationsBuildTypeStatic
	}
	return ""
}

func BuildTypeFromString(str string) BuildType {
	switch str {
	case models.ApplicationsBuildTypeRuntime:
		return BuildTypeRuntime
	case models.ApplicationsBuildTypeStatic:
		return BuildTypeStatic
	default:
		panic(fmt.Errorf("UNKNOWN BUILD TYPE: %s", str))
	}
}
