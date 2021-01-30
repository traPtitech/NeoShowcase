package appmanager

import (
	"context"
	"errors"
)

// ErrNotFound 汎用エラー 見つかりません
var ErrNotFound = errors.New("not found")

// Manager アプリマネージャー
type Manager interface {
	// GetApp 指定したIDのアプリを取得します
	GetApp(appID string) (App, error)
	// GetAppByRepository 指定したリポジトリURLのアプリを取得します
	GetAppByRepository(repo string) (App, error)
	Shutdown(ctx context.Context) error
}
