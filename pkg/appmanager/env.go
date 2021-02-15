package appmanager

// Env アプリ環境モデル
type Env interface {
	// GetID 環境IDを返します
	GetID() string
	// GetBranchName 環境に対応するブランチ名を返します
	GetBranchName() string
	// GetBuildType 環境のビルドタイプを返します
	GetBuildType() BuildType
	// SetupWebsite 環境にHTTPウェブサイトを設定します
	SetupWebsite(fqdn string, httpPort int) error
}
