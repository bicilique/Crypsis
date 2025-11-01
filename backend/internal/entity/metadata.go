package entity

import (
	"time"

	"gorm.io/gorm"
)

type Metadata struct {
	ID        string         `gorm:"type:varchar(36);not null;primaryKey"`
	FileID    string         `gorm:"type:varchar(36);index;not null;constraint:OnDelete:CASCADE"`
	Hash      string         `gorm:"type:varchar(256);not null"`
	EncHash   string         `gorm:"type:varchar(256);index;null"`
	KeyUID    string         `gorm:"type:varchar(256);index;null"`
	EncKey    string         `gorm:"type:text;not null"`
	KeyAlgo   string         `gorm:"type:varchar(64);not null"`
	VersionID string         `gorm:"type:varchar(64);null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Associations
	File Files `gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Metadata) TableName() string {
	return "metadata"
}
