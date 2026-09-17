package url

import (
	"log/slog"

	"kira-url/internal/cache"
	"kira-url/internal/modules/click"

	"github.com/go-chi/chi/v5"
)

type URLModule struct {
	URLHandler *URLHandler
}

func NewURLModule(urlRepository URLRepository, cache *cache.Cache, clickService *click.ClickService, log *slog.Logger) *URLModule {
	service := newURLService(urlRepository, log)
	handler := newURLHandler(service, cache, clickService, log)
	return &URLModule{
		URLHandler: handler,
	}
}

func (module *URLModule) RegisterRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/urls", func(r chi.Router) {
		r.Get("/{code}", module.URLHandler.FindURLByShortCode)
		r.Post("/", module.URLHandler.SaveURLShorter)
	})

	return router
}
