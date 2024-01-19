package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

// LocalStorage ローカルストレージ
type LocalStorage struct {
	localDir string
}

// NewLocalStorage LocalStorageを生成する。指定したディレクトリはすでに存在している必要がある。
func NewLocalStorage(dir string) (*LocalStorage, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return &LocalStorage{}, fmt.Errorf("dir %s doesn't exist", dir)
	}
	if !fi.IsDir() {
		return &LocalStorage{}, fmt.Errorf("dir %s is not a directory", dir)
	}

	return &LocalStorage{localDir: dir}, nil
}

// Save ファイルを保存する
func (ls *LocalStorage) Save(filename string, src io.Reader) error {
	file, err := os.Create(ls.getFilePath(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

// Open ファイルを取得する
func (ls *LocalStorage) Open(filename string) (io.ReadCloser, error) {
	r, err := os.Open(ls.getFilePath(filename))
	if err != nil {
		return nil, domain.ErrFileNotFound
	}
	return r, nil
}

// Delete ファイルを削除する
func (ls *LocalStorage) Delete(filename string) error {
	path := ls.getFilePath(filename)
	if _, err := os.Stat(path); err != nil {
		return domain.ErrFileNotFound
	}
	return os.Remove(path)
}

func (ls *LocalStorage) getFilePath(filename string) string {
	return filepath.Join(ls.localDir, filename)
}
