package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

func New(c Config) (*sql.DB, error) {
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	conf.DBName = c.Database
	conf.User = c.Username
	conf.Passwd = c.Password
	conf.ParseTime = true
	conf.Apply(mysql.Charset("utf8mb4", "utf8mb4_general_ci"))

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
