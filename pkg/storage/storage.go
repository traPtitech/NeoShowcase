package storage

import (
	"archive/tar"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	// ErrFileNotFound ファイルが存在しない
	ErrFileNotFound = errors.New("not found")
)

// Storage ストレージインターフェース
type Storage interface {
	Save(filename string, src io.Reader) error
	Open(filename string) (io.ReadCloser, error)
	Delete(filename string) error
	Move(filename, destPath string) error // LocalFile to Storage
}

// SaveArtifact Artifactをtar形式で保存する
func SaveArtifact(strg Storage, filename string, dstpath string, db *sql.DB, buildid string, sid string) error {
	stat, _ := os.Stat(filename)
	artifact := models.Artifact{
		ID:         sid,
		BuildLogID: buildid,
		Size:       stat.Size(),
		CreatedAt:  time.Now(),
	}

	if err := strg.Move(filename, dstpath); err != nil {
		return fmt.Errorf("failed to save artifact tar file: %w", err)
	}

	if err := artifact.Insert(context.Background(), db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert artifact entry: %w", err)
	}
	return nil
}

// SaveLogFile ログファイルを保存する
func SaveLogFile(strg Storage, filename string, dstpath string, buildid string) error {
	if err := strg.Move(filename, dstpath); err != nil {
		return fmt.Errorf("failed to move build log: %w", err)
	}
	return nil
}

// ExtractTarToDir tarファイルをディレクトリに展開する
func ExtractTarToDir(strg Storage, sourcePath, destPath string) error {
	inputFile, err := strg.Open(sourcePath)
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
