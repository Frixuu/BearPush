package util

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForInterrupt blocks the current goroutine
// until it receives either a SIGINT or SIGTERM.
func WaitForInterrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	signal.Stop(quit)
}
