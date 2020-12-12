package storage

import (
	"database/sql"
	"io"
)

type Storage interface {
	Save(filename string, src io.Reader) error
	Open(filename string) (io.ReadCloser, error)
	Delete(filename string) error
	DeleteAll(dirname string) error
	Move(sourcePath, destPath string) error // LocalDir to Storage
	FileExists(filename string) bool
	SaveFileAsTar(filename string, dstpath string, db *sql.DB, buildid string, sid string) error
	SaveLogFile(filename string, dstpath string, buildid string) error
	ExtractTarToDir(sourcePath, destPath string) error // Storage to LocalDir
}
