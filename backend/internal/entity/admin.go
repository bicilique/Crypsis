package entity

import (
	"time"

	"gorm.io/gorm"
)

type Admins struct {
	ID        string         `gorm:"type:varchar(36);not null;primaryKey"`
	Username  string         `gorm:"type:varchar(255);not null;unique"`
	ClientID  string         `gorm:"type:varchar(255);not null;unique;index"`
	Secret    string         `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Salt      string         `gorm:"type:varchar(255);not null"`
}

func (Admins) TableName() string {
	return "admins"
}
