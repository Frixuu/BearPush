package server

import (
	"os"
)

// DeterminePort chooses a port for the web server.
func DeterminePort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
