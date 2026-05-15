package click

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type ClickWorker struct {
	clickService *ClickService
	db           *gorm.DB
	logger       *slog.Logger
}

func NewClickWorker(service *ClickService, db *gorm.DB, logger *slog.Logger) *ClickWorker {
	return &ClickWorker{
		clickService: service,
		db:           db,
		logger:       logger,
	}
}
func (cw *ClickWorker) Start() {
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()
	cw.logger.Info("====Schedule start=====")
	cw.logger.Info("====Visit Count update=====")
	for range ticker.C {
		// batch := cw.clickService.FlushClicks()
		cw.logger.Info("Task executed at:", "TIME", time.Now())
	}
}
