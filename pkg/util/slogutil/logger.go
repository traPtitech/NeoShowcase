package slogutil

import (
	"log/slog"
	"os"
)

var levelVar = new(slog.LevelVar) // default level is INFO

// InitLogger initializes the global slog logger with CustomHandler
func InitLogger(service, version, revision string) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     levelVar,
	}
	handler := NewCustomHandler(os.Stdout, opts)
	logger := slog.New(handler).With(
		slog.Group("service",
			"name", service,
			"version", version,
			"revision", revision,
		),
	)
	slog.SetDefault(logger)
}

func SetLogLevel(level slog.Level) {
	levelVar.Set(level)
}
