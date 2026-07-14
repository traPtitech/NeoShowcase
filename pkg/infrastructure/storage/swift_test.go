package storage

import (
	"io"
	"strings"
	"testing"

	"github.com/ncw/swift/v2"
	"github.com/ncw/swift/v2/swifttest"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func TestSwiftStorage(t *testing.T) {
	server, err := swifttest.NewSwiftServer("localhost")
	require.NoError(t, err)
	t.Cleanup(server.Close)

	const container = "test-container"
	conn := &swift.Connection{
		AuthUrl:  server.AuthURL,
		UserName: swifttest.TEST_ACCOUNT,
		ApiKey:   swifttest.TEST_ACCOUNT,
	}
	require.NoError(t, conn.Authenticate(t.Context()))
	require.NoError(t, conn.ContainerCreate(t.Context(), container, nil))

	storage, err := NewSwiftStorage(
		container,
		swifttest.TEST_ACCOUNT,
		swifttest.TEST_ACCOUNT,
		"",
		"",
		server.AuthURL,
	)
	require.NoError(t, err)

	const (
		filename = "test-file"
		content  = "test-content"
	)
	require.NoError(t, storage.Save(filename, strings.NewReader(content)))

	reader, err := storage.Open(filename)
	require.NoError(t, err)
	actual, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	require.Equal(t, content, string(actual))

	require.NoError(t, storage.Delete(filename))
	_, err = storage.Open(filename)
	require.ErrorIs(t, err, domain.ErrFileNotFound)
	require.ErrorIs(t, storage.Delete(filename), domain.ErrFileNotFound)
}
