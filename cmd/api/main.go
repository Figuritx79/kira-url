package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kira-url/internal/app"
	"kira-url/internal/config"

	"github.com/lmittmann/tint"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func logLevel(cfg *config.Config) slog.Level {

	var logLevel slog.Level
	if cfg.Logger.LogLevel == "info" {
		logLevel = slog.LevelInfo
	}

	if cfg.Logger.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}

	if cfg.Logger.LogLevel == "warn" {
		logLevel = slog.LevelWarn
	}

	if cfg.Logger.LogLevel == "error" {
		logLevel = slog.LevelError
	}
	return logLevel
}

func main() {
	cfg, err := config.New()

	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	level := logLevel(cfg)
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: level}))

	localServer := app.New(cfg, logger)

	httpServer := newHttpServer(cfg, logger, localServer)

	logger.Info("starting server", slog.Group("server", "addr", httpServer.Addr))

	localServer.InitializeProcess()
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(httpServer, done)

	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	//
	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func newHttpServer(config *config.Config, logger *slog.Logger, app *app.App) *http.Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Server.Port),
		Handler:      app.RegisterRoutes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return httpServer
}
