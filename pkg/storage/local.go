package storage

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func (ls *LocalStorage) SaveDirToTar(filename string, dstpath string, db *sql.DB, buildid string, sid string) error {
	// filename: ローカルにおけるファイルの名前
	// dstpath: ローカルにおけるファイルのパス
	stat, _ := os.Stat(filename)
	artifact := models.Artifact{
		ID:         sid,
		BuildLogID: buildid,
		Size:       stat.Size(),
		CreatedAt:  time.Now(),
	}

	if err := ls.Move(filename, filepath.Join("/neoshowcase/artifacts", fmt.Sprintf("%s.tar", sid))); err != nil {
		return fmt.Errorf("failed to save artifact tar file: %w", err)
	}

	if err := artifact.Insert(context.Background(), db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert artifact entry: %w", err)
	}
	return nil
}

func (ls *LocalStorage) SaveLogFile(filename string, dstpath string, buildid string) error {
	// filename: ローカルにおけるファイルの名前
	// dstpath: ローカルにおけるファイルのパス
	if err := ls.Move(filename, dstpath); err != nil {
		fmt.Errorf("failed to move build log: %w", err)
	}
	return nil
}

func (ls *LocalStorage) getFilePath(filename string) string {
	return filepath.Join(ls.localDir, filename)
}
