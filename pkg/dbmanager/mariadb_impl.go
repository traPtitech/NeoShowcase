package dbmanager

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type mariaDBManagerImpl struct {
	db *sql.DB
}

type MariaDBConfig struct {
	Host          string
	Port          int
	AdminUser     string
	AdminPassword string
}

func NewMariaDBManager(c MariaDBConfig) (MariaDBManager, error) {
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	conf.User = c.AdminUser
	conf.Passwd = c.AdminPassword
	conf.Params = map[string]string{
		"charset": "utf8mb4",
	}
	conf.Collation = "utf8mb4_general_ci"
	conf.ParseTime = true

	// DB接続
	connector, err := mysql.NewConnector(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new connector: %w", err)
	}
	db := sql.OpenDB(connector)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &mariaDBManagerImpl{db: db}, nil
}

func (m *mariaDBManagerImpl) Create(ctx context.Context, args CreateArgs) error {
	panic("implement me") // TODO
}

func (m *mariaDBManagerImpl) Delete(ctx context.Context, args DeleteArgs) error {
	panic("implement me") // TODO
}

func (m *mariaDBManagerImpl) Close(ctx context.Context) error {
	return m.db.Close()
}
