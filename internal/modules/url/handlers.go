package url

import (
	"errors"
	"log/slog"
	"net/http"

	"kira-url/internal/cache"
	"kira-url/internal/constants"
	"kira-url/internal/modules/click"
	"kira-url/internal/request"
	"kira-url/internal/response"
	"kira-url/internal/validator"

	dtourl "kira-url/internal/dto/url"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	Service      *URLService
	cache        *cache.Cache
	clickService *click.ClickService
	log          *slog.Logger
}

func newURLHandler(service *URLService, cache *cache.Cache, clickService *click.ClickService, log *slog.Logger) *URLHandler {
	return &URLHandler{
		cache:        cache,
		Service:      service,
		clickService: clickService,
		log:          log,
	}
}

func (handler *URLHandler) SaveURLShorter(w http.ResponseWriter, r *http.Request) {
	var createUrl dtourl.CreatURL

	if err := request.GetRequestBody(r, &createUrl); err != nil {
		handler.log.Error("Error getting the body", "error", err.Error())
		response.ServerError(w, r)
		return
	}

	if !validator.NotEmpty(createUrl.OriginalURL) {
		handler.log.Warn("Save URL", "ORIGINAL URL EMPTY", createUrl.OriginalURL)
		response.BadRequest(w, r, errors.New("original url can't be empty"))
		return
	}
	if !validator.IsURL(createUrl.OriginalURL) {
		handler.log.Warn("Save URL", "Invalid URL", createUrl.OriginalURL)
		response.BadRequest(w, r, ErrInvalidURL)
		return
	}

	shortURLResponse, found, err := handler.Service.FindByURL(createUrl.OriginalURL)
	if err != nil {
		handler.log.Error("Error searching the URL", "error", err.Error())
		response.ServerError(w, r)
	}

	if found {
		shortURLResponse.CompleteShortURL = constants.BaseDomain + shortURLResponse.ShortURL
		_, err := handler.cache.Get(shortURLResponse.ShortURL)
		if err != nil {
			handler.cache.Set(shortURLResponse.ShortURL, []byte(createUrl.OriginalURL), constants.BASE_TTL)
		}
		response.OK(w, &shortURLResponse, "url found successfully")
		return
	}
	shortURL, err := handler.Service.Save(&createUrl)
	if err != nil {
		if errors.Is(err, ErrInvalidCustomCode) {
			handler.log.Error("Error saving the URL", "error", err.Error())
			response.BadRequest(w, r, err)
			return
		}
		if errors.Is(err, ErrMinRunesCustomCode) {
			handler.log.Error("Error saving the URL", "error", err.Error())
			response.BadRequest(w, r, err)
			return
		}
		if errors.Is(err, ErrMaxRunesCustomCode) {
			handler.log.Error("Error saving the URL", "error", err.Error())
			response.BadRequest(w, r, err)
			return
		}

		handler.log.Error("Error saving the URL", "error", err.Error())
		response.ServerError(w, r)
		return
	}

	response.OK(w, &shortURL, "url created successfully")
}

func (handler *URLHandler) FindURLByShortCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if !validator.NotEmpty(code) {
		response.BadRequest(w, r, errors.New("The code can't be empty"))
		return
	}
	if !validator.MinRunes(code, 6) {
		response.BadRequest(w, r, errors.New("The code must be at least 3 characters"))
		return
	}
	if !validator.MaxRunes(code, 6) {
		response.BadRequest(w, r, errors.New("The code must be at most 10 characters"))
		return
	}
	foundURL, err := handler.cache.Get(code)
	if err != nil {
		handler.log.Debug("Cache", "DON'T FOUND", code)
		url, err := handler.Service.FindByShortURL(code)
		if err != nil {
			if errors.Is(err, ErrURLNotFound) {
				handler.log.Warn("Find short code", "NOT_FOUND", err.Error())
				response.NotFound(w, r)
				return
			}
			handler.log.Error("Error finding the URL", "error", err.Error())
			response.ServerError(w, r)
			return
		}

		if url == nil {
			response.NotFound(w, r)
			return
		}

		handler.clickService.IncrementClicks(code)
		handler.cache.Set(code, []byte(url.OriginalURL), constants.BASE_TTL)

		response.Found(w, &url, "url found successfully", url.OriginalURL)
		return
	}

	handler.log.Debug("Cache", "FOUND", code)

	handler.clickService.IncrementClicks(code)

	urlReponse := dtourl.URLResponse{
		OriginalURL: string(foundURL),
	}

	response.Found(w, &urlReponse, "url found successfully", urlReponse.OriginalURL)
}
