package dbmanager

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/traPtitech/neoshowcase/pkg/domain"
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

func NewMariaDBManager(c MariaDBConfig) (domain.MariaDBManager, error) {
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

func (m *mariaDBManagerImpl) Create(ctx context.Context, args domain.CreateArgs) error {
	db := m.db
	_, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", args.Database))
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s'", args.Database, args.Password))
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, fmt.Sprintf("GRANT ALL ON %s.* TO %s", args.Database, args.Database))
	if err != nil {
		return err
	}
	return nil
}

func (m *mariaDBManagerImpl) Delete(ctx context.Context, args domain.DeleteArgs) error {
	db := m.db
	_, err := db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", args.Database))
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, fmt.Sprintf("DROP USER %s", args.Database))
	if err != nil {
		return err
	}
	return nil
}

func (m *mariaDBManagerImpl) Poll(ctx context.Context) error {
	db := m.db
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	err = db.QueryRowContext(ctx, "SELECT 1 LIMIT 1").Scan(new(int))
	if err != nil {
		return err
	}
	return nil
}

func (m *mariaDBManagerImpl) Close(_ context.Context) error {
	return m.db.Close()
}
