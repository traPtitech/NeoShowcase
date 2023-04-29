package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

func getEnvOrDefault(key string, fallback string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return fallback
}

func main() {
	header := getEnvOrDefault("HEADER", "X-Showcase-User")
	port := getEnvOrDefault("PORT", "4181")
	user := getEnvOrDefault("USER", "toki")

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set(header, user)
		return c.NoContent(http.StatusOK)
	})

	go func() {
		err := e.Start(":" + port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %+v", err)
		}
	}()

	cli.WaitSIGINT()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := e.Shutdown(ctx)
	if err != nil {
		log.Fatalf("failed to shutdown: %+v", err)
	}
}
