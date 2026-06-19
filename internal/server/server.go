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
	Port         int
	logger       *slog.Logger
	db           database.Service
	urlModule    *url.URLModule
	cache        *cache.Cache
	clickWorker  *click.ClickWorker
	clickService *click.ClickService
}

func NewServer(logger *slog.Logger) *Server {
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
		Port:         port,
		logger:       logger,
		db:           db,
		urlModule:    urlModule,
		cache:        cache,
		clickWorker:  clickWorker,
		clickService: clickService,
	}

	return NewServer
}
func (s *Server) InitializeProcess() {
	// In this function we can start diferent process, like cron/schedule,etc
	go s.clickWorker.Start(s.urlModule.URLHandler.Service.BatchUpdate)
}

func NewHttpServer(server *Server, logger *slog.Logger) *http.Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", server.Port),
		Handler:      server.RegisterRoutes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return httpServer
}
