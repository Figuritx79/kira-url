package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"kira-url/internal/database"
	"kira-url/internal/env"
)

var autoMigrate = env.GetEnvBool("AUTO_MIGRATE", true)

type Server struct {
	port   int
	logger *slog.Logger
	db     database.Service
}

func NewServer(logger *slog.Logger) *http.Server {
	port := env.GetEnvInt("PORT", 8080)
	NewServer := &Server{
		port:   port,
		logger: logger,
		db:     database.New(autoMigrate),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		ErrorLog:     slog.NewLogLogger(NewServer.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("starting server", slog.Group("server", "addr", server.Addr))
	return server
}
