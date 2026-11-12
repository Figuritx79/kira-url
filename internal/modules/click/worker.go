package click

import (
	"log/slog"
	"time"
	"uuid"

	"kira-url/internal/database/models"
)

type urlService interface {
	BatchUpdate(updateURLs []models.URL) error
}
type ClickWorker struct {
	clickService *ClickService
	logger       *slog.Logger
	urlService   urlService
}

func NewClickWorker(clickservice *ClickService, urlService urlService, logger *slog.Logger) *ClickWorker {
	return &ClickWorker{
		clickService: clickservice,
		logger:       logger,
		urlService:   urlService,
	}
}

func (cw *ClickWorker) Start() {
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()
	cw.logger.Info("====Schedule start=====")
	cw.logger.Info("====Visit Count update=====")
	for range ticker.C {
		taksID := uuid.NewV7()
		cw.logger.Info("Task ID:", "ID", taksID)
		cw.logger.Info("Task executed at:", "TIME", time.Now())
		batch := cw.clickService.FlushClicks()
		if err := cw.urlService.BatchUpdate(batch); err != nil {
			cw.logger.Error("Task error:", "ERROR", err)
		}

	}
	cw.logger.Info("====Schedule end=====")
}
