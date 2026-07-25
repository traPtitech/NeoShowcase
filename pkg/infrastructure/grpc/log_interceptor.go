package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
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
// would flatten everything to its message — hence the unwrap.
func logError(err error) any {
	if oopsErr, ok := oops.AsOops(err); ok {
		return oopsErr
	}
	return err
}

func (l *LogInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
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

		start := time.Now()
		response, err := next(ctx, request)
		elapsed := fmt.Sprintf("%.3fs", time.Since(start).Seconds())

		if err == nil || connect.IsNotModifiedError(err) {
			slog.InfoContext(ctx, "unary request succeeded",
				"procedure", request.Spec().Procedure,
				"duration_sec", elapsed,
			)
		} else {
			// treat context.Canceled error as connect.CodeCanceled
			if errors.Is(err, context.Canceled) {
				err = connect.NewError(connect.CodeCanceled, err)
			}
			switch connect.CodeOf(err) {
			case connect.CodeUnknown, connect.CodeInternal, connect.CodeUnavailable, connect.CodeDeadlineExceeded, connect.CodeUnimplemented, connect.CodeDataLoss:
				slog.ErrorContext(ctx, "unary request failed with server error",
					"procedure", request.Spec().Procedure,
					"duration_sec", elapsed,
					"error", logError(err),
					"status", connect.CodeOf(err).String(),
				)

			case connect.CodeCanceled, connect.CodeInvalidArgument, connect.CodeNotFound, connect.CodeAlreadyExists, connect.CodePermissionDenied, connect.CodeFailedPrecondition, connect.CodeAborted, connect.CodeOutOfRange, connect.CodeUnauthenticated, connect.CodeResourceExhausted:
				slog.WarnContext(ctx, "unary request failed with client error",
					"procedure", request.Spec().Procedure,
					"duration_sec", elapsed,
					"error", logError(err),
					"status", connect.CodeOf(err).String(),
				)
			}
		}

		return response, err
	}
}

func (l *LogInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (l *LogInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
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

		open := time.Now()
		streamID := domain.NewID()
		slog.InfoContext(ctx, "stream opened",
			"stream_id", streamID,
			"procedure", shc.Spec().Procedure,
		)

		err := next(ctx, shc)

		elapsed := fmt.Sprintf("%.3fs", time.Since(open).Seconds())
		if err == nil {
			slog.InfoContext(ctx, "stream closed",
				"stream_id", streamID,
				"procedure", shc.Spec().Procedure,
				"duration_sec", elapsed,
			)
		} else {
			slog.ErrorContext(ctx, "stream closed",
				"stream_id", streamID,
				"procedure", shc.Spec().Procedure,
				"duration_sec", elapsed,
				"error", logError(err),
				"status", connect.CodeOf(err).String(),
			)
		}

		return err
	}
}
