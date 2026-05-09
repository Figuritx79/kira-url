package models

import (
	"time"

	"kira-url/internal/funcs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID          uuid.UUID      `gorm:"primaryKey; type:uuid;"`
	ShortURL    string         `gorm:"unique; not null; index; type:varchar(50)"`
	OriginalURL string         `gorm:"not null; type:text" `
	VisitCount  int64          `gorm:"default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (url *URL) BeforeCreate(tx *gorm.DB) error {
	if url.ID == uuid.Nil {
		newUUID, err := funcs.GenerateUUID()
		if err != nil {
			return err
		}
		url.ID = newUUID
	}
	return nil
}
