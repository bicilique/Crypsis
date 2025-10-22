package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type FileLogs struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"` // Auto-incrementing ID
	ActorID   string    `gorm:"type:text;not null"`
	ActorType string    `gorm:"type:text;not null;check:actor_type IN ('user', 'client', 'system','admin')"`
	FileID    string    `gorm:"not null;index"` // Removed type:uuid to support SQLite
	Action    string    `gorm:"type:text;not null;check:action IN ('upload', 'download', 'update', 'delete', 'recover','encrypt', 'decrypt','re-key')"`
	Timestamp time.Time `gorm:"autoCreateTime"` // Changed to autoCreateTime for SQLite compatibility
	IP        string    `gorm:"type:text"`      // Changed from inet to text for SQLite
	UserAgent string    `gorm:"type:text"`      // Client info
	Metadata  JSONB     `gorm:"type:text"`      // Changed from jsonb to text for SQLite, will serialize to JSON
}

// BeforeCreate hook to set timestamp if not already set
func (f *FileLogs) BeforeCreate(tx *gorm.DB) error {
	if f.Timestamp.IsZero() {
		f.Timestamp = time.Now()
	}
	return nil
}

// GormDataType returns the SQL data type for FileID based on the database driver
func (FileLogs) GormDataType(dialect gorm.Dialector) string {
	switch dialect.Name() {
	case "postgres":
		return "uuid"
	default:
		return "text"
	}
}

// GormDBDataType returns the SQL data type for Timestamp based on the database driver
func (f FileLogs) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		switch field.Name {
		case "FileID":
			return "uuid"
		case "Timestamp":
			return "timestamptz"
		case "IP":
			return "inet"
		case "Metadata":
			return "jsonb"
		}
	case "sqlite":
		switch field.Name {
		case "FileID":
			return "text"
		case "Timestamp":
			return "datetime"
		case "IP":
			return "text"
		case "Metadata":
			return "text"
		}
	}
	return ""
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
