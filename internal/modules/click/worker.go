package click

import (
	"context"
	"log/slog"
	"time"
	"uuid"

	"kira-url/internal/database/models"
)

type urlService interface {
	BatchUpdate(ctx context.Context, updateURLs []models.URL) error
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

func (cw *ClickWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()
	cw.logger.Info("====Schedule start=====")
	cw.logger.Info("====Visit Count update=====")

	for {
		select {
		case <-ctx.Done():
			// ctx is already cancelled here, so use a fresh bounded context
			// to persist the pending clicks during shutdown.
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			cw.flush(flushCtx)
			cancel()
			cw.logger.Info("====Schedule end=====")
			return
		case <-ticker.C:
			taskID := uuid.NewV7()
			cw.logger.Info("Task ID:", "ID", taskID)
			cw.logger.Info("Task executed at:", "TIME", time.Now())
			cw.flush(ctx)
		}
	}
}

func (cw *ClickWorker) flush(ctx context.Context) {
	batch := cw.clickService.FlushClicks()
	if len(batch) == 0 {
		return
	}
	if err := cw.urlService.BatchUpdate(ctx, batch); err != nil {
		cw.logger.Error("Task error:", "ERROR", err)
	}
}
