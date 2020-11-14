package util

import (
	"archive/tar"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func MoveFile(sourcePath, destPath string) error {
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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func ExtractTarToDir(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	defer inputFile.Close()

	tr := tar.NewReader(inputFile)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("bad tar file: %w", err)
		}

		path := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), header.FileInfo().Mode()|os.ModeDir|100); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			_, err = io.Copy(file, tr)
			_ = file.Close()
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

		default:
			log.Debug("skip:", header)
		}
	}
	return nil
}
