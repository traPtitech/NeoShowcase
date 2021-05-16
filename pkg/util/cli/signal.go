package cli

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitSIGINT() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	signal.Stop(quit)
	for range quit {
	}
	close(quit)
}
