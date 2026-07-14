package storage

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func TestS3Storage(t *testing.T) {
	const (
		bucket   = "test-bucket"
		filename = "test-file"
		content  = "test-content"
	)

	var object []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+bucket+"/"+filename {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodPut:
			var err error
			object, err = io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("failed to read request body: %v", err)
				http.Error(w, "failed to read request body", http.StatusInternalServerError)
				return
			}
			w.Header().Set("ETag", `"test-etag"`)
		case http.MethodGet:
			if object == nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Length", strconv.Itoa(len(object)))
			_, err := w.Write(object)
			if err != nil {
				t.Errorf("failed to write response body: %v", err)
			}
		case http.MethodDelete:
			object = nil
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)

	storage, err := NewS3Storage(bucket, "access-key", "access-secret", "us-east-1", server.URL)
	require.NoError(t, err)
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
}
