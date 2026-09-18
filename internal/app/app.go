package app

import (
	"log/slog"

	_ "github.com/joho/godotenv/autoload"

	"kira-url/internal/cache"
	"kira-url/internal/config"
	"kira-url/internal/database"
)

type App struct {
	config       *config.Config
	logger       *slog.Logger
	db           database.Service
	Repositories *repositories
	Services     *services
	Modules      *modules
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	// Define db
	db := database.New(cfg)
	// Define fast access layer cache
	cache := cache.NewCache(100 * 1024 * 1024)

	repositories := buildRepositories(db)
	services := buildServices()
	modules := buildModules(db, logger, *repositories, *services, cache)

	NewServer := &App{
		logger:       logger,
		db:           db,
		Repositories: repositories,
		Services:     services,
		Modules:      modules,
	}

	return NewServer
}

func (a *App) InitializeProcess() {
	// In this function we can start diferent process, like cron/schedule,etc
	go a.Modules.Click.Start(a.Modules.URL.URLHandler.Service.BatchUpdate)
}
