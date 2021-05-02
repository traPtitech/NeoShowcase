package userdb

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
