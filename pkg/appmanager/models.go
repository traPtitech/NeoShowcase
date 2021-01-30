package appmanager

// App アプリモデル
type App interface {
	// GetID アプリIDを返します
	GetID() string
	// GetName アプリ名を返します
	GetName() string

	// Start アプリを起動します
	Start(args AppStartArgs) error
}

type AppStartArgs struct {
	// 起動したいビルドID
	BuildID string
}
