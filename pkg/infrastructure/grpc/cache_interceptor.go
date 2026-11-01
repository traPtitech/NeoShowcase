package grpc

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"

	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type CacheInterceptor struct{}

func NewCacheInterceptor() *CacheInterceptor {
	return &CacheInterceptor{}
}

func (c *CacheInterceptor) hash(resp any) (string, bool) {
	respStringer, ok := resp.(fmt.Stringer)
	if ok {
		respStr := respStringer.String()
		return hash.XXH3Hex([]byte(respStr)), true
	}
	return "", false
}

// Server returns the interceptor that sets ETag on unary GET responses and
// answers 304 Not Modified when it matches If-None-Match.
func (c *CacheInterceptor) Server() connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			if spec.StreamType != connect.StreamTypeUnary {
				return next(ctx, spec, stream)
			}
			httpInfo, ok := connecthttp.ServerInfoForContext(ctx)
			if !ok || httpInfo.HTTPMethod() != http.MethodGet {
				return next(ctx, spec, stream)
			}
			info, _ := connect.CallInfoForServerContext(ctx)
			return next(ctx, spec, &cacheServerStream{ServerStream: stream, cache: c, info: info})
		}
	}
}

// cacheServerStream intercepts the single response of a unary call. The
// generated handler returns the error of Send, so returning the not-modified
// error from Send ends the call with 304 instead of sending the body.
type cacheServerStream struct {
	connect.ServerStream
	cache *CacheInterceptor
	info  *connect.CallInfo
}

func (s *cacheServerStream) Send(msg any) error {
	respHash, ok := s.cache.hash(msg)
	if !ok {
		return s.ServerStream.Send(msg)
	}
	respETag := `"` + respHash + `"`
	s.info.ResponseHeader().Set("ETag", respETag)
	if s.info.RequestHeader().Get("If-None-Match") == respETag {
		return connecthttp.NewNotModifiedError()
	}
	return s.ServerStream.Send(msg)
}
