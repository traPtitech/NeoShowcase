package appmanager

import "context"

// App アプリモデル
type App interface {
	// GetID アプリIDを返します
	GetID() string
	// GetName アプリ名を返します
	GetName() string

	// Start アプリを起動します
	Start(args AppStartArgs) error
	RequestBuild(ctx context.Context, envID string) error
}

type AppStartArgs struct {
	// 操作したい環境ID
	EnvironmentID string
	// 起動したいビルドID
	BuildID string
}
