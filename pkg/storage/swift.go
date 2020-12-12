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

	"github.com/ncw/swift"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SwiftStorage struct {
	Container string
	Conn      *swift.Connection
}

func (ss *SwiftStorage) Save(filename string, src io.Reader) error {
	_, err := ss.Conn.ObjectPut(ss.Container, filename, src, true, "", "", swift.Headers{})
	return err
}

func (ss *SwiftStorage) Open(filename string) (io.ReadCloser, error) {
	file, _, err := ss.Conn.ObjectOpen(ss.Container, filename, true, nil)
	if err != nil {
		if err == swift.ObjectNotFound {
			return nil, fmt.Errorf("not found: %w", err)
		}
		return nil, err
	}
	return file, nil
}

func (ss *SwiftStorage) Delete(filename string) error {
	err := ss.Conn.ObjectDelete(ss.Container, filename)
	if err != nil {
		if err == swift.ObjectNotFound {
			return fmt.Errorf("not found: %w", err)
		}
		return err
	}
	return nil
}

func (ss *SwiftStorage) Move(sourcePath, destPath string) error {
	// Move LocalDir to Swift Storage
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	_, err = ss.Conn.ObjectPut(ss.Container, destPath, inputFile, true, "", "", swift.Headers{})
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

func (ss *SwiftStorage) SaveFileAsTar(filename string, dstpath string, db *sql.DB, buildid string, sid string) error {
	// filename: ローカルにおけるファイルの名前
	// dstpath: SwiftStorageにおけるファイルのパス(階層構造ではないのでファイル名)
	stat, _ := os.Stat(filename)
	artifact := models.Artifact{
		ID:         sid,
		BuildLogID: buildid,
		Size:       stat.Size(),
		CreatedAt:  time.Now(),
	}

	if err := ss.Move(filename, dstpath); err != nil {
		return fmt.Errorf("failed to save artifact tar file: %w", err)
	}

	if err := artifact.Insert(context.Background(), db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert artifact entry: %w", err)
	}
	return nil
}

func (ss *SwiftStorage) SaveLogFile(filename string, dstpath string, buildid string) error {
	// filename: ローカルにおけるファイルの名前
	// dstpath: SwiftStorageにおけるファイルのパス(階層構造ではないのでファイル名)
	if err := ss.Move(filename, dstpath); err != nil {
		return fmt.Errorf("failed to move build log: %w", err)
	}
	return nil
}

func (ss *SwiftStorage) FileExists(filename string) bool {
	_, _, err := ss.Conn.Object(ss.Container, filename)
	return err == nil
}

func (ss *SwiftStorage) ExtractTarToDir(sourcePath, destPath string) error {
	// filename: SwiftStorageにおけるファイルのパス(階層構造ではないのでファイル名)
	// dstpath: ローカルにおけるファイルの名前
	inputFile, err := ss.Open(sourcePath)
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
