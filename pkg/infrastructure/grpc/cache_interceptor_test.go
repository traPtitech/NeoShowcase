package grpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb/pbconnect"
)

type systemInfoHandler struct {
	pbconnect.UnimplementedAPIServiceHandler
}

func (systemInfoHandler) GetSystemInfo(context.Context, *emptypb.Empty) (*pb.SystemInfo, error) {
	return &pb.SystemInfo{PublicKey: "key"}, nil
}

func TestCacheInterceptor(t *testing.T) {
	server := connect.NewServer(NewCacheInterceptor().Server())
	pbconnect.RegisterAPIServiceHandler(server, systemInfoHandler{})
	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	url := ts.URL + pbconnect.APIServiceGetSystemInfoProcedure
	get := func(t *testing.T, ifNoneMatch string) *http.Response {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, url+"?connect=v1&encoding=json&message={}", nil)
		require.NoError(t, err)
		if ifNoneMatch != "" {
			req.Header.Set("If-None-Match", ifNoneMatch)
		}
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = res.Body.Close() })
		return res
	}

	first := get(t, "")
	require.Equal(t, http.StatusOK, first.StatusCode)
	etag := first.Header.Get("ETag")
	require.NotEmpty(t, etag)

	t.Run("matching If-None-Match returns 304", func(t *testing.T) {
		res := get(t, etag)
		assert.Equal(t, http.StatusNotModified, res.StatusCode)
		assert.Equal(t, etag, res.Header.Get("ETag"))
	})

	t.Run("stale If-None-Match returns the body", func(t *testing.T) {
		res := get(t, `"stale"`)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, etag, res.Header.Get("ETag"))
	})

	t.Run("POST gets no ETag", func(t *testing.T) {
		res, err := http.Post(url, "application/json", strings.NewReader("{}"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = res.Body.Close() })
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.Header.Get("ETag"))
	})
}
