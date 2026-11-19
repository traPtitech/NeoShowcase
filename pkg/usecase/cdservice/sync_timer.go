package cdservice

import (
	"log/slog"
	"time"
)

// syncTimer records how long each step of a deployment sync takes.
// The zero value is ready to use.
type syncTimer struct {
	steps []any
}

// measure runs fn and appends its duration under name, whether fn fails or not.
func (t *syncTimer) measure(name string, fn func() error) error {
	start := time.Now()
	err := fn()
	t.steps = append(t.steps, slog.Duration(name, time.Since(start)))
	return err
}

// attr returns the recorded durations as a single log attribute, in the order they were measured.
func (t *syncTimer) attr() slog.Attr {
	return slog.Group("steps", t.steps...)
}
