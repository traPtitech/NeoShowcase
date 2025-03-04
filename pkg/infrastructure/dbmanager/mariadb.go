package dbmanager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/go-sql-driver/mysql"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type mariaDBManagerImpl struct {
	db *sql.DB
	c  MariaDBConfig
}

type MariaDBConfig struct {
	Host          string `mapstructure:"host" yaml:"host"`
	Port          int    `mapstructure:"port" yaml:"port"`
	AdminUser     string `mapstructure:"adminUser" yaml:"adminUser"`
	AdminPassword string `mapstructure:"adminPassword" yaml:"adminPassword"`
}

func NewMariaDBManager(c MariaDBConfig) (domain.MariaDBManager, error) {
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	conf.User = c.AdminUser
	conf.Passwd = c.AdminPassword
	// TODO: As of go-sql-driver/mysql v1.9.0, the "charset" parameter is not functional.
	// Support for this parameter is expected in a future version:
	// https://github.com/go-sql-driver/mysql/pull/1679
	conf.Params = map[string]string{
		"charset": "utf8mb4",
	}
	conf.Collation = "utf8mb4_general_ci"
	conf.ParseTime = true
	conf.InterpolateParams = true

	// DB接続
	connector, err := mysql.NewConnector(conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new connector")
	}
	db := sql.OpenDB(connector)
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}
	db.SetMaxOpenConns(1024)
	db.SetMaxIdleConns(1024)
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(time.Minute)

	return &mariaDBManagerImpl{db: db, c: c}, nil
}

func (m *mariaDBManagerImpl) GetHost() (host string, port int) {
	return m.c.Host, m.c.Port
}

func (m *mariaDBManagerImpl) Create(ctx context.Context, args domain.CreateArgs) error {
	if strings.ContainsRune(args.Database, '`') {
		return fmt.Errorf("backtick(`) in database name are not permitted")
	}
	if _, err := m.db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", args.Database)); err != nil {
		return err
	}
	if _, err := m.db.ExecContext(ctx, "CREATE USER IF NOT EXISTS ? IDENTIFIED BY ?", args.Database, args.Password); err != nil {
		return err
	}
	if _, err := m.db.ExecContext(ctx, fmt.Sprintf("GRANT ALL ON `%s`.* TO ?", args.Database), args.Database); err != nil {
		return err
	}
	return nil
}

func (m *mariaDBManagerImpl) Delete(ctx context.Context, args domain.DeleteArgs) error {
	if strings.ContainsRune(args.Database, '`') {
		return fmt.Errorf("backtick(`) in database name are not permitted")
	}
	if _, err := m.db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", args.Database)); err != nil {
		return err
	}
	if _, err := m.db.ExecContext(ctx, "DROP USER IF EXISTS ?", args.Database); err != nil {
		return err
	}
	return nil
}

func (m *mariaDBManagerImpl) IsExist(ctx context.Context, name string) (bool, error) {
	rows, err := m.db.QueryContext(ctx, "SHOW DATABASES WHERE `Database` = ?", name)
	if err != nil {
		return false, err
	}

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func (m *mariaDBManagerImpl) Close(_ context.Context) error {
	return m.db.Close()
}
