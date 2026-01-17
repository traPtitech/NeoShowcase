package slogutil

import (
	"log/slog"
	"os"
)

// InitLogger initializes the global slog logger with CustomHandler
func InitLogger(level slog.Level) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}
	handler := NewCustomHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
