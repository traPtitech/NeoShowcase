package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/samber/oops"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/slogutil"
	"go.opentelemetry.io/otel/trace"
)

type LogInterceptor struct {
}

func NewLogInterceptor() *LogInterceptor {
	return &LogInterceptor{}
}

// logError picks the value to log for err.
//
// oops.OopsError implements slog.LogValuer, so logging it expands into structured
// fields (the wrap chain, the code, With() attributes and the stack trace). The
// *connect.Error the boundary returns implements neither, so logging it directly
// would flatten everything to its message — hence the unwrap. For the same reason
// an error without oops in its chain is logged by its cause: the message of a
// *connect.Error is the client-facing one and carries none of the detail.
func logError(err error) any {
	if oopsErr, ok := oops.AsOops(err); ok {
		return oopsErr
	}
	if connectErr, ok := errors.AsType[*connect.Error](err); ok {
		if cause := connectErr.Unwrap(); cause != nil {
			return cause.Error()
		}
	}
	return err
}

// errorLevel maps err to the level the boundary logs it at.
//
// A client error means the server behaved correctly and rejected the request, so
// only failures on our side reach Error. Unclassified codes are ours by default.
func errorLevel(err error) slog.Level {
	if errors.Is(err, context.Canceled) {
		return slog.LevelWarn
	}
	switch connect.CodeOf(err) {
	case connect.CodeCanceled, connect.CodeInvalidArgument, connect.CodeNotFound,
		connect.CodeAlreadyExists, connect.CodePermissionDenied, connect.CodeFailedPrecondition,
		connect.CodeAborted, connect.CodeOutOfRange, connect.CodeUnauthenticated,
		connect.CodeResourceExhausted:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

// Server returns the interceptor that logs each request at the boundary.
func (l *LogInterceptor) Server() connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			// Add trace_id to context if available
			if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
				traceID := span.SpanContext().TraceID().String()
				ctx = context.WithValue(ctx, slogutil.TraceIDKey, traceID)
			}

			// Add user_id to context if exists
			user, ok := web.TryGetUser(ctx)
			if ok && user.ID != "" {
				ctx = context.WithValue(ctx, slogutil.UserIDKey, user.ID)
			}

			if spec.StreamType == connect.StreamTypeUnary {
				return logUnary(ctx, spec, stream, next)
			}
			return logStream(ctx, spec, stream, next)
		}
	}
}

func logUnary(ctx context.Context, spec connect.Spec, stream connect.ServerStream, next connect.ServerFunc) error {
	start := time.Now()
	err := next(ctx, spec, stream)
	elapsed := fmt.Sprintf("%.3fs", time.Since(start).Seconds())

	if err == nil || connecthttp.IsNotModifiedError(err) {
		slog.InfoContext(ctx, "unary request succeeded",
			"procedure", spec.Procedure,
			"duration_sec", elapsed,
		)
	} else {
		// treat context.Canceled error as connect.CodeCanceled
		if errors.Is(err, context.Canceled) {
			err = connect.NewError(connect.CodeCanceled, err.Error()).WithCause(err)
		}
		slog.Log(ctx, errorLevel(err), "unary request failed",
			"procedure", spec.Procedure,
			"duration_sec", elapsed,
			"error", logError(err),
			"status", connect.CodeOf(err).String(),
		)
	}

	return err
}

func logStream(ctx context.Context, spec connect.Spec, stream connect.ServerStream, next connect.ServerFunc) error {
	open := time.Now()
	streamID := domain.NewID()
	slog.InfoContext(ctx, "stream opened",
		"stream_id", streamID,
		"procedure", spec.Procedure,
	)

	err := next(ctx, spec, stream)

	elapsed := fmt.Sprintf("%.3fs", time.Since(open).Seconds())
	if err == nil {
		slog.InfoContext(ctx, "stream closed",
			"stream_id", streamID,
			"procedure", spec.Procedure,
			"duration_sec", elapsed,
		)
	} else {
		slog.Log(ctx, errorLevel(err), "stream closed",
			"stream_id", streamID,
			"procedure", spec.Procedure,
			"duration_sec", elapsed,
			"error", logError(err),
			"status", connect.CodeOf(err).String(),
		)
	}

	return err
}
