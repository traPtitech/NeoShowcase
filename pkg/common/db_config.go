package common

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

// DBConfig データベース接続設定
type DBConfig struct {
	// Host ホスト名 (default: 127.0.0.1)
	Host string `mapstructure:"host" yaml:"host"`
	// Port ポート番号 (default: 3306)
	Port int `mapstructure:"port" yaml:"port"`
	// Username ユーザー名 (default: root)
	Username string `mapstructure:"username" yaml:"username"`
	// Password パスワード (default: password)
	Password string `mapstructure:"password" yaml:"password"`
	// Database データベース名 (default: neoshowcase)
	Database string `mapstructure:"database" yaml:"database"`
	// Connection コネクション設定
	Connection struct {
		// MaxOpen 最大オープン接続数. 0は無制限 (default: 0)
		MaxOpen int `mapstructure:"maxOpen" yaml:"maxOpen"`
		// MaxIdle 最大アイドル接続数 (default: 2)
		MaxIdle int `mapstructure:"maxIdle" yaml:"maxIdle"`
		// LifeTime 待機接続維持時間. 0は無制限 (default: 0)
		LifeTime int `mapstructure:"lifetime" yaml:"lifetime"`
	} `mapstructure:"connection" yaml:"connection"`
}

// Connect この設定でDBに接続
func (c *DBConfig) Connect() (*sql.DB, error) {
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	conf.DBName = c.Database
	conf.User = c.Username
	conf.Passwd = c.Password
	conf.Params = map[string]string{
		"charset": "utf8mb4",
	}
	conf.Collation = "utf8mb4_general_ci"
	conf.ParseTime = true

	connector, err := mysql.NewConnector(conf)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(c.Connection.MaxOpen)
	db.SetMaxIdleConns(c.Connection.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(c.Connection.LifeTime) * time.Second)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
