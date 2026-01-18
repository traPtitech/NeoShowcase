package slogutil

import (
	"context"
	"io"
	"log/slog"
)

// CustomHandler adds trace_id and user_id from context
type CustomHandler struct {
	baseHandler slog.Handler
}

// NewCustomHandler creates a new CustomHandler
func NewCustomHandler(w io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		baseHandler: slog.NewTextHandler(w, opts),
	}
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	// Extract trace_id from context
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if traceIDStr, ok := traceID.(string); ok && traceIDStr != "" {
			record.AddAttrs(slog.String("trace_id", traceIDStr))
		}
	}

	// Extract user_id from context
	if userID := ctx.Value(UserIDKey); userID != nil {
		if userIDStr, ok := userID.(string); ok && userIDStr != "" {
			record.AddAttrs(slog.String("user_id", userIDStr))
		}
	}

	return h.baseHandler.Handle(ctx, record)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		baseHandler: h.baseHandler.WithAttrs(attrs),
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		baseHandler: h.baseHandler.WithGroup(name),
	}
}
