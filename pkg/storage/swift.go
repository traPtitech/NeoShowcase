package storage

import (
	"fmt"
	"io"

	"github.com/ncw/swift"
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

func (ss *SwiftStorage) FileExists(filename string) bool {
	_, _, err := ss.conn.Object(ss.container, filename)
	return err == nil
}
