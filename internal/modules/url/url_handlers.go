package url

import (
	"errors"
	"log/slog"
	"net/http"

	"kira-url/internal/cache"
	"kira-url/internal/constants"
	"kira-url/internal/httperrors"
	"kira-url/internal/request"
	"kira-url/internal/response"
	"kira-url/internal/transport/httptransport"
	"kira-url/internal/validator"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *urlService
	cache   *cache.Cache
	log     *slog.Logger
}

func newURLHandler(service *urlService, cache *cache.Cache, log *slog.Logger) *URLHandler {
	return &URLHandler{
		cache:   cache,
		service: service,
		log:     log,
	}
}

func (handler *URLHandler) SaveURLShorter(w http.ResponseWriter, r *http.Request) {
	var createUrl CreatURL

	if err := request.GetRequestBody[CreatURL](r, &createUrl); err != nil {
		handler.log.Error("Error getting the body", "error", err.Error())
		httperrors.ServerError(w, r, err)
		return
	}

	if !validator.NotEmpty(createUrl.OriginalURL) {
		handler.log.Warn("Save URL", "ORIGINAL URL EMPTY", createUrl.OriginalURL)
		httperrors.BadRequest(w, r, errors.New("original url can't be empty"))
		return
	}
	if !validator.IsURL(createUrl.OriginalURL) {
		handler.log.Warn("Save URL", "Invalid URL", createUrl.OriginalURL)
		httperrors.BadRequest(w, r, ErrInvalidURL)
		return
	}

	shortURLResponse, found, err := handler.service.FindByURL(createUrl.OriginalURL)
	if err != nil {
		handler.log.Error("Error searching the URL", "error", err.Error())
		httperrors.ServerError(w, r, err)
	}
	if found {
		shortURLResponse.CompleteShortURL = constants.BaseDomain + shortURLResponse.ShortURL
		_, err := handler.cache.Get(shortURLResponse.ShortURL)
		if err != nil {
			handler.cache.Set(shortURLResponse.ShortURL, []byte(createUrl.OriginalURL), constants.BASE_TTL)
		}
		json := httptransport.JSONResponse[ShortURLResponse]{
			Data:    shortURLResponse,
			Type:    httptransport.Success,
			Message: "Url found successfully",
		}

		if err := response.JSON[httptransport.JSONResponse[ShortURLResponse]](w, http.StatusOK, json); err != nil {
			handler.log.Error("Error sending the response", "error", err.Error())
			httperrors.ServerError(w, r, err)
			return
		}
		return
	}
	shortURL, err := handler.service.Save(&createUrl)
	if err != nil {
		handler.log.Error("Error saving the URL", "error", err.Error())
		httperrors.ServerError(w, r, err)
		return
	}

	json := httptransport.JSONResponse[URLCompleteResponse]{
		Data:    shortURL,
		Type:    httptransport.Success,
		Message: "Url created successfully",
	}

	if err := response.JSON[httptransport.JSONResponse[URLCompleteResponse]](w, http.StatusOK, json); err != nil {
		handler.log.Error("Error sending the response", "error", err.Error())
		httperrors.ServerError(w, r, err)
		return
	}
}

func (handler *URLHandler) FindURLByShortCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if !validator.NotEmpty(code) {
		httperrors.BadRequest(w, r, errors.New("The code can't be empty"))
		return
	}
	if !validator.MinRunes(code, 6) {
		httperrors.BadRequest(w, r, errors.New("The code must be at least 3 characters"))
		return
	}
	if !validator.MaxRunes(code, 6) {
		httperrors.BadRequest(w, r, errors.New("The code must be at most 10 characters"))
		return
	}
	foundURL, err := handler.cache.Get(code)
	if err != nil {
		handler.log.Debug("Cache", "DON'T FOUND", code)
		url, err := handler.service.FindByShortURL(code)
		if err != nil {
			if errors.Is(err, ErrURLNotFound) {
				handler.log.Warn("Find short code", "NOT_FOUND", err.Error())
				httperrors.NotFound(w, r)
				return
			}
			handler.log.Error("Error finding the URL", "error", err.Error())
			httperrors.ServerError(w, r, err)
			return
		}

		if url == nil {
			httperrors.NotFound(w, r)
			return
		}

		json := httptransport.JSONResponse[URLResponse]{
			Data:    url,
			Type:    httptransport.Success,
			Message: "Url found successfully",
		}
		handler.cache.Set(code, []byte(url.OriginalURL), constants.BASE_TTL)
		if err := response.JSON[httptransport.JSONResponse[URLResponse]](w, http.StatusSeeOther, json); err != nil {
			handler.log.Error("Error sending the response", "error", err.Error())
			httperrors.ServerError(w, r, err)
			return
		}
		return
	}

	handler.log.Debug("Cache", "FOUND", code)
	urlReponse := URLResponse{
		OriginalURL: string(foundURL),
	}
	json := httptransport.JSONResponse[URLResponse]{
		Data:    &urlReponse,
		Type:    httptransport.Success,
		Message: "Url found successfully",
	}

	if err := response.JSON[httptransport.JSONResponse[URLResponse]](w, http.StatusSeeOther, json); err != nil {
		handler.log.Error("Error sending the response", "error", err.Error())
		httperrors.ServerError(w, r, err)
		return
	}
}
