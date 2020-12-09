package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	localDir string
}

func NewLocalStorage(dir string) (LocalStorage, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return LocalStorage{}, fmt.Errorf("dir doesn't exist: %w", err)
	}
	if !fi.IsDir() {
		return LocalStorage{}, fmt.Errorf("dir is not a directory: %w", err)
	}

	return LocalStorage{localDir: dir}, nil
}

func (ls *LocalStorage) Save(filename string, src io.Reader) error {
	file, err := os.Create(ls.getFilePath(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

func (ls *LocalStorage) Open(filename string) (io.ReadCloser, error) {
	r, err := os.Open(ls.getFilePath(filename))
	if err != nil {
		return nil, fmt.Errorf("not found: %w", err)
	}
	return r, nil
}

func (ls *LocalStorage) Delete(filename string) error {
	path := ls.getFilePath(filename)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("not found: %w", err)
	}
	return os.Remove(path)
}

func (ls *LocalStorage) Move(sourcePath, destPath string) error {
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

func (ls *LocalStorage) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (ls *LocalStorage) getFilePath(filename string) string {
	return filepath.Join(ls.localDir, filename)
}
