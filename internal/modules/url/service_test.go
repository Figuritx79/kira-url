package url

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"kira-url/internal/constants"
	"kira-url/internal/database/models"
	dtourl "kira-url/internal/dto/url"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var errRepositoryFailure = errors.New("repository failure")

var _ URLRepository = (*fakeRepository)(nil)

// fakeRepository is an injectable test double for URLRepository.
type fakeRepository struct {
	urls []models.URL

	findByShortURLResponse *dtourl.URLResponse
	findByShortURLError    error

	findByURLResponse *dtourl.ShortURLResponse
	findByURLError    error

	saveError error
}

func (repository *fakeRepository) FindByShortURL(code string) (*dtourl.URLResponse, error) {
	return repository.findByShortURLResponse, repository.findByShortURLError
}

func (repository *fakeRepository) FindByURL(url string) (*dtourl.ShortURLResponse, error) {
	return repository.findByURLResponse, repository.findByURLError
}

func (repository *fakeRepository) Save(url *models.URL) error {
	if repository.saveError != nil {
		return repository.saveError
	}
	repository.urls = append(repository.urls, *url)
	return nil
}

func (repository *fakeRepository) Update(url models.URL, code string) error {
	return nil
}

func (repository *fakeRepository) Updates(urls []models.URL) error {
	return nil
}

func newTestService(repository URLRepository) *URLService {
	return newURLService(repository, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestFindByShortURL(t *testing.T) {
	tests := []struct {
		name                   string
		findByShortURLResponse *dtourl.URLResponse
		findByShortURLError    error
		want                   *dtourl.URLResponse
		wantErr                error
	}{
		{
			name:                "record not found maps to ErrURLNotFound",
			findByShortURLError: gorm.ErrRecordNotFound,
			wantErr:             ErrURLNotFound,
		},
		{
			name:                   "nil response without error returns nil",
			findByShortURLResponse: nil,
			findByShortURLError:    nil,
			want:                   nil,
			wantErr:                nil,
		},
		{
			name: "returns the stored url response",
			findByShortURLResponse: &dtourl.URLResponse{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OriginalURL: "https://example.com/page",
				VisitCount:  7,
			},
			want: &dtourl.URLResponse{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				OriginalURL: "https://example.com/page",
				VisitCount:  7,
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &fakeRepository{
				findByShortURLResponse: test.findByShortURLResponse,
				findByShortURLError:    test.findByShortURLError,
			}
			service := newTestService(repository)

			got, err := service.FindByShortURL("any-code")

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("FindByShortURL() error = %v, want %v", err, test.wantErr)
			}

			if test.want == nil {
				if got != nil {
					t.Fatalf("FindByShortURL() response = %+v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("FindByShortURL() response = nil, want %+v", test.want)
			}
			if got.ID != test.want.ID || got.OriginalURL != test.want.OriginalURL || got.VisitCount != test.want.VisitCount {
				t.Fatalf("FindByShortURL() response = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestFindByURL(t *testing.T) {
	tests := []struct {
		name              string
		findByURLResponse *dtourl.ShortURLResponse
		findByURLError    error
		want              *dtourl.ShortURLResponse
		wantFound         bool
		wantErr           error
	}{
		{
			name:           "record not found returns no match and no error",
			findByURLError: gorm.ErrRecordNotFound,
			want:           nil,
			wantFound:      false,
			wantErr:        nil,
		},
		{
			name:           "repository failure is propagated",
			findByURLError: errRepositoryFailure,
			want:           nil,
			wantFound:      false,
			wantErr:        errRepositoryFailure,
		},
		{
			name: "returns the stored short url",
			findByURLResponse: &dtourl.ShortURLResponse{
				ShortURL:         "my-slug",
				CompleteShortURL: constants.BaseDomain + "my-slug",
			},
			want: &dtourl.ShortURLResponse{
				ShortURL:         "my-slug",
				CompleteShortURL: constants.BaseDomain + "my-slug",
			},
			wantFound: true,
			wantErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &fakeRepository{
				findByURLResponse: test.findByURLResponse,
				findByURLError:    test.findByURLError,
			}
			service := newTestService(repository)

			got, found, err := service.FindByURL("https://example.com/page")

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("FindByURL() error = %v, want %v", err, test.wantErr)
			}
			if found != test.wantFound {
				t.Fatalf("FindByURL() found = %v, want %v", found, test.wantFound)
			}

			if test.want == nil {
				if got != nil {
					t.Fatalf("FindByURL() response = %+v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("FindByURL() response = nil, want %+v", test.want)
			}
			if got.ShortURL != test.want.ShortURL || got.CompleteShortURL != test.want.CompleteShortURL {
				t.Fatalf("FindByURL() response = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestSaveWithCustomCode(t *testing.T) {
	tests := []struct {
		name         string
		originalURL  string
		customCode   string
		saveError    error
		wantErr      error
		wantResponse *dtourl.URLCompleteResponse
		wantStored   int
		wantSlug     string
		wantIsCustom bool
	}{
		{
			name:        "custom code below the minimum length",
			originalURL: "https://example.com/page",
			customCode:  "ab",
			wantErr:     ErrMinRunesCustomCode,
			wantStored:  0,
		},
		{
			name:        "custom code above the maximum length",
			originalURL: "https://example.com/page",
			customCode:  strings.Repeat("a", 41),
			wantErr:     ErrMaxRunesCustomCode,
			wantStored:  0,
		},
		{
			name:        "custom code with invalid characters",
			originalURL: "https://example.com/page",
			customCode:  "my_slug!",
			wantErr:     ErrInvalidCustomCode,
			wantStored:  0,
		},
		{
			name:        "repository failure is propagated",
			originalURL: "https://example.com/page",
			customCode:  "my-slug",
			saveError:   errRepositoryFailure,
			wantErr:     errRepositoryFailure,
			wantStored:  0,
		},
		{
			name:        "stores the custom code",
			originalURL: "https://example.com/page",
			customCode:  "my-slug",
			wantResponse: &dtourl.URLCompleteResponse{
				ShortURL:    constants.BaseDomain + "my-slug",
				OriginalURL: "https://example.com/page",
			},
			wantStored:   1,
			wantSlug:     "my-slug",
			wantIsCustom: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &fakeRepository{saveError: test.saveError}
			service := newTestService(repository)

			request := &dtourl.CreatURL{
				OriginalURL: test.originalURL,
				CustomCode:  test.customCode,
			}

			got, err := service.Save(request)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Save() error = %v, want %v", err, test.wantErr)
			}

			if test.wantResponse == nil {
				if got != nil {
					t.Fatalf("Save() response = %+v, want nil", got)
				}
			} else {
				if got == nil {
					t.Fatalf("Save() response = nil, want %+v", test.wantResponse)
				}
				if got.ShortURL != test.wantResponse.ShortURL || got.OriginalURL != test.wantResponse.OriginalURL {
					t.Fatalf("Save() response = %+v, want %+v", got, test.wantResponse)
				}
			}

			if len(repository.urls) != test.wantStored {
				t.Fatalf("stored URLs = %d, want %d", len(repository.urls), test.wantStored)
			}

			if test.wantStored == 1 {
				stored := repository.urls[0]
				if stored.ShortURL != test.wantSlug {
					t.Fatalf("stored ShortURL = %q, want %q", stored.ShortURL, test.wantSlug)
				}
				if stored.OriginalURL != test.originalURL {
					t.Fatalf("stored OriginalURL = %q, want %q", stored.OriginalURL, test.originalURL)
				}
				if stored.IsCustom != test.wantIsCustom {
					t.Fatalf("stored IsCustom = %v, want %v", stored.IsCustom, test.wantIsCustom)
				}
			}
		})
	}
}

func TestSaveWithBase62Code(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		saveError   error
		wantErr     error
		wantStored  int
	}{
		{
			name:        "repository failure is propagated",
			originalURL: "https://example.com/page",
			saveError:   errRepositoryFailure,
			wantErr:     errRepositoryFailure,
			wantStored:  0,
		},
		{
			name:        "stores a generated base62 code",
			originalURL: "https://example.com/page",
			wantStored:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &fakeRepository{saveError: test.saveError}
			service := newTestService(repository)

			request := &dtourl.CreatURL{OriginalURL: test.originalURL}

			got, err := service.Save(request)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Save() error = %v, want %v", err, test.wantErr)
			}

			if len(repository.urls) != test.wantStored {
				t.Fatalf("stored URLs = %d, want %d", len(repository.urls), test.wantStored)
			}

			if test.wantStored == 0 {
				if got != nil {
					t.Fatalf("Save() response = %+v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("Save() response = nil, want a generated response")
			}

			stored := repository.urls[0]
			if stored.ShortURL == "" {
				t.Fatalf("stored ShortURL is empty, want a generated base62 code")
			}
			if stored.OriginalURL != test.originalURL {
				t.Fatalf("stored OriginalURL = %q, want %q", stored.OriginalURL, test.originalURL)
			}
			if stored.IsCustom {
				t.Fatalf("stored IsCustom = true, want false")
			}
			if got.ShortURL != constants.BaseDomain+stored.ShortURL {
				t.Fatalf("Save() ShortURL = %q, want %q", got.ShortURL, constants.BaseDomain+stored.ShortURL)
			}
			if got.OriginalURL != test.originalURL {
				t.Fatalf("Save() OriginalURL = %q, want %q", got.OriginalURL, test.originalURL)
			}
		})
	}
}
