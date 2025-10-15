package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type FileLogs struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"` // Auto-incrementing ID
	ActorID   string    `gorm:"type:text;not null"`
	ActorType string    `gorm:"type:text;not null;check:actor_type IN ('user', 'client', 'system','admin')"`
	FileID    string    `gorm:"type:uuid;not null;index"` // Keep UUID for file_id
	Action    string    `gorm:"type:text;not null;check:action IN ('upload', 'download', 'update', 'delete', 'recover','encrypt', 'decrypt','re-key')"`
	Timestamp time.Time `gorm:"type:timestamptz;default:now()"` // Auto timestamp
	IP        string    `gorm:"type:inet"`                      // User's IP address
	UserAgent string    `gorm:"type:text"`                      // Client info
	Metadata  JSONB     `gorm:"type:jsonb"`                     // File metadata (size, hash)
}

type JSONB map[string]interface{}

// Scan allows JSONB to read from a PostgreSQL jsonb column
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONB: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value allows JSONB to write into a PostgreSQL jsonb column
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (FileLogs) TableName() string {
	return "file_logs"
}

// UserID    string    `gorm:"type:uuid;not null;index"` // Keep UUID for user_id
