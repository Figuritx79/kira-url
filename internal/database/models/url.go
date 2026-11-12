package models

import (
	"strings"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type URL struct {
	ID          uuid.UUID      `gorm:"primaryKey; type:uuid;"`
	ShortURL    string         `gorm:"unique; not null; index; type:varchar(40)"`
	OriginalURL string         `gorm:"not null; type:text" `
	VisitCount  int64          `gorm:"default:0"`
	IsCustom    bool           `gorm:"type:boolean; default:false; not null;"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (url *URL) BeforeCreate(tx *gorm.DB) error {
	if url.ID == uuid.Nil() {
		url.ID = uuid.NewV7()
	}
	if url.IsCustom {
		url.ShortURL = strings.ToLower(url.ShortURL)
	}
	return nil
}
