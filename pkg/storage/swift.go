package storage

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ncw/swift"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SwiftStorage struct {
	container string
	conn      *swift.Connection
}

func NewSwiftStorage(container, userName, apiKey, tenant, tenantID, authURL string) (SwiftStorage, error) {
	conn := &swift.Connection{
		AuthUrl:  authURL,
		UserName: userName,
		ApiKey:   apiKey,
		Tenant:   tenant,
		TenantId: tenantID,
	}

	if err := conn.Authenticate(); err != nil {
		return SwiftStorage{}, err
	}

	if _, _, err := conn.Container(container); err != nil {
		return SwiftStorage{}, err
	}

	s := SwiftStorage{
		container: container,
		conn:      conn,
	}

	return s, nil
}

func (ss *SwiftStorage) Save(filename string, src io.Reader) error {
	_, err := ss.conn.ObjectPut(ss.container, filename, src, true, "", "", swift.Headers{})
	return err
}

func (ss *SwiftStorage) Open(filename string) (io.ReadCloser, error) {
	file, _, err := ss.conn.ObjectOpen(ss.container, filename, true, nil)
	if err != nil {
		if err == swift.ObjectNotFound {
			return nil, fmt.Errorf("not found: %w", err)
		}
		return nil, err
	}
	return file, nil
}

func (ss *SwiftStorage) Delete(filename string) error {
	err := ss.conn.ObjectDelete(ss.container, filename)
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
	_, err = ss.conn.ObjectPut(ss.container, destPath, inputFile, true, "", "", swift.Headers{})
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

func (ss *SwiftStorage) SaveDirToTar(filename string, dstpath string, db *sql.DB, buildid string, sid string) error {
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
	_, _, err := ss.conn.Object(ss.container, filename)
	return err == nil
}
