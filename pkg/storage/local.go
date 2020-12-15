package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
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
		return nil, ErrFileNotFound
	}
	return r, nil
}

// Delete ファイルを削除する
func (ls *LocalStorage) Delete(filename string) error {
	path := ls.getFilePath(filename)
	if _, err := os.Stat(path); err != nil {
		return ErrFileNotFound
	}
	return os.Remove(path)
}

// Move ファイルをdestPathへ移動する
func (ls *LocalStorage) Move(filename, destPath string) error {
	sourcePath := ls.getFilePath(filename)
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	outputFile, err := os.Create(destPath)
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
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}
	return nil
}

func (ls *LocalStorage) getFilePath(filename string) string {
	return ls.localDir + "/" + filename
}
