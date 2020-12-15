package dbmanager

import "context"

// Manager アプリデータベースマネージャー
type Manager interface {
	Close(ctx context.Context) error
}

// MariaDBManager MariaDBマネージャー
type MariaDBManager interface {
	// Create データベースが存在しない場合、作成します
	Create(ctx context.Context, args CreateArgs) error
	// Delete データベースが存在する場合、削除します
	Delete(ctx context.Context, args DeleteArgs) error
	Close(ctx context.Context) error
}

// MongoManager Mongoマネージャー
type MongoManager interface {
	// Create データベースが存在しない場合、作成します
	Create(ctx context.Context, args CreateArgs) error
	// Delete データベースが存在する場合、削除します
	Delete(ctx context.Context, args DeleteArgs) error
	Close(ctx context.Context) error
}

// CreateArgs データベース作成引数
type CreateArgs struct {
	// Database データベース名
	Database string
	// Password アクセスパスワード
	Password string
}

// DeleteArgs データベース削除引数
type DeleteArgs struct {
	// Database データベース名
	Database string
}
