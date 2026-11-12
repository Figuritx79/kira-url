package app

import (
	"log/slog"

	"kira-url/internal/cache"
	"kira-url/internal/config"
	"kira-url/internal/database"
	"kira-url/internal/modules/click"
	"kira-url/internal/modules/url"
	repository "kira-url/internal/repository/url"
)

type repositories struct {
	url url.URLRepository
}
type services struct {
	click *click.ClickService
	url   *url.URLService
}
type modules struct {
	url   url.URLModule
	click click.ClickWorker
}

func buildRepositories(db database.Service) *repositories {
	urlRepository := repository.NewURLRepository(db.GetDB())

	return &repositories{
		url: urlRepository,
	}
}

func buildServices(log *slog.Logger, repos *repositories, cfg *config.Config) *services {
	clickService := click.NewClickService()
	urlService := url.NewURLService(repos.url, cfg.Server.Domain, log)
	return &services{
		click: clickService,
		url:   urlService,
	}
}

func buildModules(db database.Service, log *slog.Logger, repos *repositories, ser *services, cache *cache.Cache, cfg *config.Config) *modules {
	// Define click module
	// Define url module
	urlModule := url.NewURLModule(repos.url, cache, ser.click, ser.url, cfg.Cache.TTL, cfg.Server.Domain, log)

	clickWorker := click.NewClickWorker(ser.click, ser.url, log)

	return &modules{
		url:   *urlModule,
		click: *clickWorker,
	}
}
