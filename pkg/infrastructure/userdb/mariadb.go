package userdb

import (
	"context"
)

// MariaDBManager MariaDBマネージャー
type MariaDBManager interface {
	// Create データベースが存在しない場合、作成します
	Create(ctx context.Context, args CreateArgs) error
	// Delete データベースが存在する場合、削除します
	Delete(ctx context.Context, args DeleteArgs) error
	Close(ctx context.Context) error
}
