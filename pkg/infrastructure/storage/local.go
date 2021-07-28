package storage

import (
	"errors"
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
		return &LocalStorage{}, errors.New("dir doesn't exist")
	}
	if !fi.IsDir() {
		return &LocalStorage{}, errors.New("dir is not a directory")
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

// Move ファイルをdestPathへ移動する
func (ls *LocalStorage) Move(filename, destPath string) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	outputFile, err := os.Create(ls.getFilePath(destPath))
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %w", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %w", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}
	return nil
}

func (ls *LocalStorage) getFilePath(filename string) string {
	return filepath.Join(ls.localDir, filename)
}
