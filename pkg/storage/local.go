package storage

import (
	"archive/tar"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type LocalStorage struct {
	LocalDir string
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

func (ls *LocalStorage) DeleteAll(dirname string) error {
	path := ls.getFilePath(dirname)
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("not found: %w", err)
	}
	return os.RemoveAll(path)
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

func (ls *LocalStorage) SaveFileAsTar(filename string, dstpath string, db *sql.DB, buildid string, sid string) error {
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
		return fmt.Errorf("failed to move build log: %w", err)
	}
	return nil
}

func (ls *LocalStorage) getFilePath(filename string) string {
	return filepath.Join(ls.LocalDir, filename)
}

func (ls *LocalStorage) ExtractTarToDir(sourcePath, destPath string) error {
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
