package server

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// Starts makes the server listen to the requests and serve the responses.
func Start(srv *http.Server, logger *zap.Logger) {
	logger.Info("Starting the server now")
	err := srv.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Info("Server closed")
	}
}
