package url

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type URLModule struct {
	URLHandler *URLHandler
}

func NewURLModule(cache cacheService, clickService clickService, urlService *URLService, TTL int, baseDomain string, log *slog.Logger) *URLModule {
	handler := newURLHandler(urlService, cache, clickService, TTL, baseDomain, log)
	return &URLModule{
		URLHandler: handler,
	}
}

func (module *URLModule) RegisterRoutes(r chi.Router) {
	r.Get("/{code}", module.URLHandler.FindURLByShortCode)
	r.Post("/", module.URLHandler.SaveURLShorter)
}
