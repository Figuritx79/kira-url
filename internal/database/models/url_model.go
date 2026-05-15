package models

import (
	"strings"
	"time"

	"kira-url/internal/funcs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID          uuid.UUID      `gorm:"primaryKey; type:uuid;"`
	ShortURL    string         `gorm:"unique; not null; index; type:varchar(10)"`
	OriginalURL string         `gorm:"not null; type:text" `
	VisitCount  int64          `gorm:"default:0"`
	IsCustom    bool           `gorm:"type:boolean; default:false; not null;"`
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
	if url.IsCustom {
		url.ShortURL = strings.ToLower(url.ShortURL)
	}
	return nil
}
