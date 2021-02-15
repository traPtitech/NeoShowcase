package appmanager

import (
	"context"
	"errors"
)

// ErrNotFound 汎用エラー 見つかりません
var ErrNotFound = errors.New("not found")

// Manager アプリマネージャー
type Manager interface {
	// GetRepo 指定したリポジトリURLのRepoを取得します
	GetRepoByURL(url string) (Repo, error)
	// GetApp 指定したIDのアプリを取得します
	GetApp(appID string) (App, error)
	// GetAppByRepository 指定したリポジトリURLのアプリを取得します
	GetAppByRepository(repo string) (App, error)
	// GetAppByEnvironment 指定した環境IDのアプリを取得します
	GetAppByEnvironment(envID string) (App, error)
	CreateApp(args CreateAppArgs) (App, error)
	Shutdown(ctx context.Context) error
}

type CreateAppArgs struct {
	Owner         string
	Name          string
	RepositoryURL string
	BranchName    string
	BuildType     BuildType
}
