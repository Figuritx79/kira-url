package click

import (
	"log/slog"
	"time"

	"kira-url/internal/database/models"
	"kira-url/internal/funcs"

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

func (cw *ClickWorker) Start(task func([]models.URL) error) {
	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()
	cw.logger.Info("====Schedule start=====")
	cw.logger.Info("====Visit Count update=====")
	for range ticker.C {
		taksID, _ := funcs.GenerateUUID()
		cw.logger.Info("Task ID:", "ID", taksID)
		cw.logger.Info("Task executed at:", "TIME", time.Now())
		batch := cw.clickService.FlushClicks()
		if err := task(batch); err != nil {
			cw.logger.Error("Task error:", "ERROR", err)
		}

	}
	cw.logger.Info("====Schedule end=====")
}
