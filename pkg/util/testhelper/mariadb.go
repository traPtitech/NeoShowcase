package testhelper

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
)

type mysqlContainer struct {
	testcontainers.Container
	host string
	port string
}

var container *mysqlContainer

func init() {
	req := testcontainers.ContainerRequest{
		Image: "mysql:8",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
		},
		ExposedPorts: []string{"3306/tcp"},
	}
	_container, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, err := _container.Host(context.Background())
	if err != nil {
		panic(err)
	}
	port, err := _container.MappedPort(context.Background(), "3306")
	if err != nil {
		panic(err)
	}

	container = &mysqlContainer{
		Container: _container,
		host:      host,
		port:      port.Port(),
	}

	ensureHealthy(container.host, container.port)

}

func OpenDB(t *testing.T) *sql.DB {
	// create database
	dbName := randomDBName()
	_, _, err := container.Exec(t.Context(), []string{"mysql", "-u", "root", "-ppassword", "-e", "CREATE DATABASE " + dbName})
	if err != nil {
		t.Fatal(err)
	}

	// migration
	migrationScript := filepath.Join(getProjectRoot(), "migrations", "entrypoint.sh")
	cmd := exec.Command(migrationScript, filepath.Join(getProjectRoot(), "migrations", "schema.sql"))
	envs := map[string]string{
		"DB_HOST": container.host,
		"DB_PORT": string(container.port),
		"DB_NAME": dbName,
	}
	cmd.Env = append(cmd.Env, os.Environ()...)
	for key, value := range envs {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	if out, err := cmd.Output(); err != nil {
		t.Logf("%s", out)
		t.Logf("%s", err.(*exec.ExitError).Stderr)
		t.Fatal(err)
	}

	cfg := mysql.Config{
		User:      "root",
		Passwd:    "password",
		Addr:      net.JoinHostPort(container.host, container.port),
		Net:       "tcp",
		DBName:    dbName,
		ParseTime: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	return db
}

func randomDBName() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var buf [10]byte
	for i := range buf {
		buf[i] = charset[rand.IntN(len(charset))]
	}
	return "test_" + string(buf[:])
}

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir { // Reached root directory
			break
		}
		dir = parentDir
	}

	panic("go.mod not found in any parent directory")
}

func ensureHealthy(host string, port string) {
	mysql.SetLogger(&mysql.NopLogger{})
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Addr:   net.JoinHostPort(host, port),
		Net:    "tcp",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	for {
		if err := db.Ping(); err != nil {
			time.Sleep(500 * time.Millisecond)
		} else {
			return
		}
	}
}
