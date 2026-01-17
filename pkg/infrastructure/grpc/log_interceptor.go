package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type LogInterceptor struct {
}

func NewLogInterceptor() *LogInterceptor {
	return &LogInterceptor{}
}

func (l *LogInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		response, err := next(ctx, request)
		elapsed := fmt.Sprintf("%.3fs", time.Since(start).Seconds())
		user := web.GetUser(ctx)
		if err == nil || connect.IsNotModifiedError(err) {
			slog.InfoContext(ctx, "unary request succeeded",
				"user", user.ID,
				"procedure", request.Spec().Procedure,
				"duration_sec", elapsed,
			)
		} else {
			switch connect.CodeOf(err) {
			case connect.CodeUnknown, connect.CodeInternal, connect.CodeUnavailable, connect.CodeDeadlineExceeded, connect.CodeUnimplemented, connect.CodeDataLoss:
				slog.ErrorContext(ctx, "unary request failed with server error",
					"user", user.ID,
					"procedure", request.Spec().Procedure,
					"duration_sec", elapsed,
					"error", err,
					"status", connect.CodeOf(err).String(),
				)

			case connect.CodeCanceled, connect.CodeInvalidArgument, connect.CodeNotFound, connect.CodeAlreadyExists, connect.CodePermissionDenied, connect.CodeFailedPrecondition, connect.CodeAborted, connect.CodeOutOfRange, connect.CodeUnauthenticated, connect.CodeResourceExhausted:
				slog.WarnContext(ctx, "unary request failed with client error",
					"user", user.ID,
					"procedure", request.Spec().Procedure,
					"duration_sec", elapsed,
					"error", err,
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
		open := time.Now()
		user := web.GetUser(ctx)
		streamID := domain.NewID()
		slog.InfoContext(ctx, "stream opened",
			"user", user.ID,
			"stream_id", streamID,
			"procedure", shc.Spec().Procedure,
		)

		err := next(ctx, shc)

		elapsed := fmt.Sprintf("%.3fs", time.Since(open).Seconds())
		if err == nil {
			slog.InfoContext(ctx, "stream closed",
				"user", user.ID,
				"stream_id", streamID,
				"procedure", shc.Spec().Procedure,
				"duration_sec", elapsed,
			)
		} else {
			slog.ErrorContext(ctx, "stream closed",
				"user", user.ID,
				"stream_id", streamID,
				"procedure", shc.Spec().Procedure,
				"duration_sec", elapsed,
				"error", err,
				"status", connect.CodeOf(err).String(),
			)
		}

		return err
	}
}
