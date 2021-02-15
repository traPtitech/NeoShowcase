package appmanager

import (
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/models"
)

type BuildType int

const (
	BuildTypeImage BuildType = iota
	BuildTypeStatic
)

func (t BuildType) String() string {
	switch t {
	case BuildTypeImage:
		return models.EnvironmentsBuildTypeImage
	case BuildTypeStatic:
		return models.EnvironmentsBuildTypeStatic
	}
	return ""
}

func BuildTypeFromString(str string) BuildType {
	switch str {
	case models.EnvironmentsBuildTypeStatic:
		return BuildTypeStatic
	case models.EnvironmentsBuildTypeImage:
		return BuildTypeImage
	default:
		panic(fmt.Errorf("UNKNOWN BUILD TYPE: %s", str))
	}
}
