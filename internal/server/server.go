package server

import (
	"log/slog"

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
	ClickWorker  *click.ClickWorker
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
		ClickWorker:  clickWorker,
		clickService: clickService,
	}

	return NewServer
}
