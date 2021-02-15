package appmanager

// Repo リポジトリモデル
type Repo interface {
	// GetID リポジトリIDを返します
	GetID() string
	// GetGitURL リポジトリのGitURLを返します
	GetGitURL() string
	// GetWebhookSecret リポジトリのWebhookSecret文字列を返します
	GetWebhookSecret() string
	// SetWebhookSecret リポジトリのWebhookSecret文字列を変更します
	SetWebhookSecret(secret string) error
}
