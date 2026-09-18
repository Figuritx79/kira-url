package app

import (
	"kira-url/internal/cache"
	"kira-url/internal/database"
	"kira-url/internal/modules/click"
	"kira-url/internal/modules/url"
	"log/slog"
)

type repositories struct {
	URL url.URLRepository
}
type services struct {
	URL   *url.URLService
	Click *click.ClickService
}
type modules struct {
	URL   url.URLModule
	Click click.ClickWorker
}

func buildRepositories(db database.Service) *repositories {
	urlRepository := url.NewURLGormRepository(db.GetDB())

	return &repositories{
		URL: urlRepository,
	}
}
func buildServices() *services {
	clickService := click.NewClickService()
	// urlService := url.New
	return &services{
		Click: clickService,
	}
}
func buildModules(db database.Service, log *slog.Logger, repos repositories, ser services, cache *cache.Cache) *modules {

	// Define click module
	clickWorker := click.NewClickWorker(ser.Click, db.GetDB(), log)
	// Define url module
	urlModule := url.NewURLModule(repos.URL, cache, ser.Click, log)

	return &modules{
		URL:   *urlModule,
		Click: *clickWorker,
	}
}
