package server

import (
	"errors"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// DeterminePort chooses a port for the web server.
func DeterminePort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

// Starts makes the server listen to the requests and serve the responses.
func Start(srv *http.Server, logger *zap.Logger) {
	logger.Info("Starting the server now")
	err := srv.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Info("Server closed")
	}
}
