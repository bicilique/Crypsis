package entity

import (
	"time"

	"gorm.io/gorm"
)

type Files struct {
	ID         string         `gorm:"type:varchar(36);not null;primaryKey"`
	Name       string         `gorm:"index;not null"`
	AppID      string         `gorm:"type:varchar(36);index"`
	UserID     string         `gorm:"type:varchar(36);index"`
	MimeType   string         `gorm:"type:varchar(255);not null"`
	Size       int64          `gorm:"not null"`
	BucketName string         `gorm:"type:varchar(255)"`
	Location   string         `gorm:"type:text; null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Files) TableName() string {
	return "files"
}
