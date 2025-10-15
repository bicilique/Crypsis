package entity

import (
	"time"

	"gorm.io/gorm"
)

type Apps struct {
	ID           string         `gorm:"type:varchar(36);not null;primaryKey"`
	Name         string         `gorm:"index;not null"`
	ClientID     string         `gorm:"type:varchar(255);not null;index"`
	ClientSecret string         `gorm:"type:varchar(255);not null"`
	IsActive     bool           `gorm:"not null"`
	Uri          string         `gorm:"type:text; null"`
	RedirectUri  string         `gorm:"type:text; null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Apps) TableName() string {
	return "apps"
}
