package storage

import (
	"io"

	"github.com/ncw/swift"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

// SwiftStorage OpenStack Swiftストレージ
type SwiftStorage struct {
	container string
	conn      *swift.Connection
}

// NewSwiftStorage 引数の情報でOpenStack Swiftストレージを生成する
func NewSwiftStorage(container, userName, apiKey, tenant, tenantID, authURL string) (*SwiftStorage, error) {
	conn := &swift.Connection{
		AuthUrl:  authURL,
		UserName: userName,
		ApiKey:   apiKey,
		Tenant:   tenant,
		TenantId: tenantID,
	}

	if err := conn.Authenticate(); err != nil {
		return &SwiftStorage{}, err
	}

	if _, _, err := conn.Container(container); err != nil {
		return &SwiftStorage{}, err
	}

	s := SwiftStorage{
		container: container,
		conn:      conn,
	}
	return &s, nil
}

// Save ファイルを保存する
func (ss *SwiftStorage) Save(filename string, src io.Reader) error {
	_, err := ss.conn.ObjectPut(ss.container, filename, src, true, "", "", swift.Headers{})
	return err
}

// Open ファイルを取得する
func (ss *SwiftStorage) Open(filename string) (io.ReadCloser, error) {
	file, _, err := ss.conn.ObjectOpen(ss.container, filename, true, nil)
	if err != nil {
		if err == swift.ObjectNotFound {
			return nil, domain.ErrFileNotFound
		}
		return nil, err
	}
	return file, nil
}

// Delete ファイルを削除する
func (ss *SwiftStorage) Delete(filename string) error {
	err := ss.conn.ObjectDelete(ss.container, filename)
	if err != nil {
		if err == swift.ObjectNotFound {
			return domain.ErrFileNotFound
		}
		return err
	}
	return nil
}
