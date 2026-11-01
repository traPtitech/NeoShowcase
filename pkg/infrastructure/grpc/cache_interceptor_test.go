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

type fakeAPIServiceHandler struct {
	pbconnect.UnimplementedAPIServiceHandler
}

func (fakeAPIServiceHandler) GetSystemInfo(context.Context, *emptypb.Empty) (*pb.SystemInfo, error) {
	return &pb.SystemInfo{PublicKey: "key"}, nil
}

func TestCacheInterceptor(t *testing.T) {
	server := connect.NewServer(NewCacheInterceptor().Server())
	pbconnect.RegisterAPIServiceHandler(server, fakeAPIServiceHandler{})
	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)

	getURL := pbconnect.APIServiceGetSystemInfoProcedure + "?connect=v1&encoding=json&message={}"
	get := func(ifNoneMatch string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, getURL, nil)
		if ifNoneMatch != "" {
			req.Header.Set("If-None-Match", ifNoneMatch)
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}

	first := get("")
	require.Equal(t, http.StatusOK, first.Code)
	etag := first.Header().Get("ETag")
	require.NotEmpty(t, etag)

	t.Run("matching If-None-Match returns 304", func(t *testing.T) {
		rec := get(etag)
		assert.Equal(t, http.StatusNotModified, rec.Code)
		assert.Equal(t, etag, rec.Header().Get("ETag"))
	})

	t.Run("stale If-None-Match returns the body", func(t *testing.T) {
		rec := get(`"stale"`)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, etag, rec.Header().Get("ETag"))
	})

	t.Run("POST gets no ETag", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, pbconnect.APIServiceGetSystemInfoProcedure, strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Header().Get("ETag"))
	})
}
