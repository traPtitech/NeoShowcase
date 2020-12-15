package httpserver

import (
	"github.com/volatiletech/null/v8"
	"time"
)

// App defines model for App.
type App struct {
	// アプリID
	Id string `json:"id"`
}

// AppDetail defines model for AppDetail.
type AppDetail struct {
	// アプリID
	Id string `json:"id"`
}

// AppKeys defines model for AppKeys.
type AppKeys struct {
	// MariaDB接続情報
	MariaDB *MariaDbKey `json:"mariadb"`
	// Mongo接続情報
	Mongo *MongoKey `json:"mongo"`
}

// Apps defines model for Apps.
type Apps []App

// BuildLog defines model for BuildLog.
type BuildLog struct {
	// ビルドID
	Id string `json:"id"`
	// ビルドステータス
	Status string `json:"status"`
	// 開始時間
	StartedAt time.Time `json:"startedAt"`
	// 完了時間
	FinishedAt null.Time `json:"finishedAt"`
}

// BuildLogs defines model for BuildLogs.
type BuildLogs []BuildLog

// EnvVar defines model for EnvVar.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// EnvVars defines model for EnvVars.
type EnvVars []EnvVar

// MariaDbKey defines model for MariaDbKey.
type MariaDbKey struct {
	// ホスト名
	Host string `json:"host"`
	// データベース名
	Database string `json:"database"`
	// 接続ユーザー名
	User string `json:"user"`
	// 接続ユーザーパスワード
	Password string `json:"password"`
}

// MongoKey defines model for MongoKey.
type MongoKey struct {
	// ホスト名
	Host string `json:"host"`
	// データベース名
	Database string `json:"database"`
	// 接続ユーザー名
	User string `json:"user"`
	// 接続ユーザーパスワード
	Password string `json:"password"`
}
