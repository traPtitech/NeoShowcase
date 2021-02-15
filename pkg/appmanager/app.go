package appmanager

import "context"

// App アプリモデル
type App interface {
	// GetID アプリIDを返します
	GetID() string
	// GetName アプリ名を返します
	GetName() string
	// GetEnvs アプリの全ての環境の配列を返します
	GetEnvs() []Env
	// GetEnvByBranchName 指定したブランチ名の環境を返します
	GetEnvByBranchName(branch string) (Env, error)
	// CreateEnv アプリに新しく環境を作成します
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
