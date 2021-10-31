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
		return models.BranchesBuildTypeImage
	case BuildTypeStatic:
		return models.BranchesBuildTypeStatic
	}
	return ""
}

func BuildTypeFromString(str string) BuildType {
	switch str {
	case models.BranchesBuildTypeStatic:
		return BuildTypeStatic
	case models.BranchesBuildTypeImage:
		return BuildTypeImage
	default:
		panic(fmt.Errorf("UNKNOWN BUILD TYPE: %s", str))
	}
}
