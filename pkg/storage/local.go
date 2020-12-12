package storage

import (
	"fmt"
	"io"
	"os"
)

type LocalStorage struct{}

func (ls *LocalStorage) Save(filename string, src io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

func (ls *LocalStorage) Open(filename string) (io.ReadCloser, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("not found: %w", err)
	}
	return r, nil
}

func (ls *LocalStorage) Delete(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return fmt.Errorf("not found: %w", err)
	}
	return os.Remove(filename)
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
