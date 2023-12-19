package grpc

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type CacheInterceptor struct{}

var _ connect.Interceptor = &CacheInterceptor{}

func NewCacheInterceptor() *CacheInterceptor {
	return &CacheInterceptor{}
}

func (c *CacheInterceptor) hash(resp connect.AnyResponse) (string, bool) {
	respStringer, ok := resp.Any().(fmt.Stringer)
	if ok {
		respStr := respStringer.String()
		return hash.XXH3Hex([]byte(respStr)), true
	}
	return "", false
}

func (c *CacheInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		resp, err := unaryFunc(ctx, request)
		if err != nil {
			return resp, err
		}
		if request.HTTPMethod() == http.MethodGet {
			// Calculate hash and set etag
			respHash, ok := c.hash(resp)
			respETag := `"` + respHash + `"`
			if ok {
				resp.Header().Set("ETag", respETag)
			}
			// Send not modified response if matched
			reqETag := request.Header().Get("If-None-Match")
			if ok && respETag == reqETag {
				return nil, connect.NewNotModifiedError(resp.Header())
			}
		}
		return resp, err
	}
}

func (c *CacheInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (c *CacheInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
