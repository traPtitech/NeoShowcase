package domain

import (
	"context"
	"time"
)

type AppMetric struct {
	Time  time.Time
	Value float64
}

type MetricsService interface {
	AvailableNames() []string
	Get(ctx context.Context, name string, app *Application, before time.Time, limit time.Duration) ([]*AppMetric, error)
}
