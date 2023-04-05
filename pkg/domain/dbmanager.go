package domain

import "context"

func DBName(applicationID string) string {
	return "nsapp_" + applicationID
}

// CreateArgs データベース作成引数
type CreateArgs struct {
	// Database データベース/ユーザー名
	Database string
	// Password アクセスパスワード
	Password string
}

// DeleteArgs データベース削除引数
type DeleteArgs struct {
	// Database データベース/ユーザー名
	Database string
}

// MariaDBManager MariaDBマネージャー
type MariaDBManager interface {
	GetHost() (host string, port int)
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
	GetHost() (host string, port int)
	// Create データベースが存在しない場合、作成します
	Create(ctx context.Context, args CreateArgs) error
	// Delete データベースが存在する場合、削除します
	Delete(ctx context.Context, args DeleteArgs) error
	// IsExist データベースが存在するか確認します
	IsExist(ctx context.Context, dbname string) (bool, error)
	Close(ctx context.Context) error
}
