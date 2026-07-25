package slogutil

import (
	"context"
	"io"
	"log/slog"
)

// CustomHandler adds trace_id and user_id from the context to every record.
//
// Rendering errors is deliberately not its job. Errors created with
// github.com/samber/oops implement slog.LogValuer, so the base handler expands
// their wrap chain, code, attached attributes and stack trace by itself — for any
// output format.
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
	traceID, hasTrace := ctx.Value(TraceIDKey).(string)
	userID, hasUser := ctx.Value(UserIDKey).(string)
	hasTrace = hasTrace && traceID != ""
	hasUser = hasUser && userID != ""

	if hasTrace || hasUser {
		// Clone before adding: the record's attributes may share backing storage
		// with the caller's.
		record = record.Clone()
		if hasTrace {
			record.AddAttrs(slog.String("trace_id", traceID))
		}
		if hasUser {
			record.AddAttrs(slog.String("user_id", userID))
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
