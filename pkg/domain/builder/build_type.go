package builder

import (
	"fmt"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type BuildType int

const (
	BuildTypeImage BuildType = iota
	BuildTypeStatic
)

func (t BuildType) String() string {
	switch t {
	case BuildTypeImage:
		return models.ApplicationsBuildTypeImage
	case BuildTypeStatic:
		return models.ApplicationsBuildTypeStatic
	}
	return ""
}

func BuildTypeFromString(str string) BuildType {
	switch str {
	case models.ApplicationsBuildTypeImage:
		return BuildTypeImage
	case models.ApplicationsBuildTypeStatic:
		return BuildTypeStatic
	default:
		panic(fmt.Errorf("UNKNOWN BUILD TYPE: %s", str))
	}
}
