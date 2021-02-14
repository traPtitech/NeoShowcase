package appmanager

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/models"
)

// App アプリモデル
type App interface {
	// GetID アプリIDを返します
	GetID() string
	// GetName アプリ名を返します
	GetName() string
	// GetEnvs アプリの全ての環境の配列を返します
	GetEnvs() []Env
	GetEnvByBranchName(branch string) (Env, error)

	CreateEnv(branchName string, buildType BuildType) (Env, error)

	// Start アプリを起動します
	Start(args AppStartArgs) error
	// RequestBuild builderにappのビルドをリクエストする
	RequestBuild(ctx context.Context, envID string) error
}

type AppStartArgs struct {
	// 操作したい環境ID
	EnvironmentID string
	// 起動したいビルドID
	BuildID string
}

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

// Env アプリ環境モデル
type Env interface {
	GetID() string
	GetBranchName() string
	GetBuildType() BuildType

	SetupWebsite(fqdn string, httpPort int) error
}
