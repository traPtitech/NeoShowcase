package storage

import (
	"io"
)

type Storage interface {
	Save(filename string, src io.Reader) error
	Open(filename string) (io.ReadCloser, error)
	Delete(filename string) error
	FileExists(filename string) bool
}
