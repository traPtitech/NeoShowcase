package domain

import "context"

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

// MariaDBManager MariaDBマネージャー
type MariaDBManager interface {
	// Create データベースが存在しない場合、作成します
	Create(ctx context.Context, args CreateArgs) error
	// Delete データベースが存在する場合、削除します
	Delete(ctx context.Context, args DeleteArgs) error
	// IsExist データベースが存在するか確認します
	IsExist(ctx context.Context, dbname string) (bool, error)
	Close(ctx context.Context) error
}

// MongoDBManager Mongoマネージャー
type MongoDBManager interface {
	// Create データベースが存在しない場合、作成します
	Create(ctx context.Context, args CreateArgs) error
	// Delete データベースが存在する場合、削除します
	Delete(ctx context.Context, args DeleteArgs) error
	// IsExist データベースが存在するか確認します
	IsExist(ctx context.Context, dbname string) (bool, error)
	Close(ctx context.Context) error
}
