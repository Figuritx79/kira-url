package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"kira-url/internal/cache"
	"kira-url/internal/database"
	"kira-url/internal/env"
	"kira-url/internal/modules/click"
	"kira-url/internal/modules/url"
)

var autoMigrate = env.GetEnvBool("AUTO_MIGRATE", true)

type Server struct {
	port         int
	logger       *slog.Logger
	db           database.Service
	urlModule    *url.URLModule
	cache        *cache.Cache
	clickWorker  *click.ClickWorker
	clickService *click.ClickService
}

func NewServer(logger *slog.Logger) *http.Server {
	// Define db
	db := database.New(autoMigrate)
	// Define fast access layer cache
	cache := cache.NewCache(100 * 1024 * 1024)
	// Define click module
	clickService := click.NewClickService()
	clickWorker := click.NewClickWorker(clickService, db.GetDB(), logger)
	// Define url module
	urlRepository := url.NewURLGormRepository(db.GetDB())
	urlModule := url.NewURLModule(urlRepository, cache, clickService, logger)
	port := env.GetEnvInt("PORT", 8080)
	NewServer := &Server{
		port:      port,
		logger:    logger,
		db:        db,
		urlModule: urlModule,
		cache:     cache,
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
	go clickWorker.Start()
	return server
}
