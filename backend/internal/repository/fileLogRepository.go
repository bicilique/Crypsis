package repository

import (
	"context"
	"crypsis-backend/internal/entity"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// fileLogsRepository implements the FileLogsRepository interface for file log data access operations.
type fileLogsRepository struct {
	db *gorm.DB
}

// NewFileLogRepository creates a new instance of FileLogsRepository.
func NewFileLogRepository(db *gorm.DB) FileLogsRepository {
	return &fileLogsRepository{db: db}
}

// Create adds a new file log entry to the database.
func (r *fileLogsRepository) Create(ctx context.Context, log *entity.FileLogs) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		slog.Error("Failed to insert file log", slog.Any("error", err))
		return fmt.Errorf("failed to insert file log: %w", err)
	}
	return nil
}

// List retrieves a paginated list of file logs with sorting options.
func (r *fileLogsRepository) List(ctx context.Context, offset, limit int, orderBy, sort string) (int64, *[]entity.FileLogs, error) {
	var total int64
	var logs []entity.FileLogs

	// Default handling
	if orderBy == "" {
		orderBy = "timestamp"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}
	allowedOrderFields := map[string]bool{
		"id":         true,
		"file_id":    true,
		"action":     true,
		"timestamp":  true,
		"ip":         true,
		"user_agent": true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "timestamp" // or another sensible default
	}

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&entity.FileLogs{}).
		Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count file logs: %w", err)
	}

	// Get paginated list
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order(orderBy + " " + sort).
		Find(&logs).Error; err != nil {
		slog.Error("Failed to get file logs", slog.Any("error", err))
		return 0, nil, fmt.Errorf("failed to get file logs: %w", err)
	}
	return total, &logs, nil
}

// GetByFileID retrieves all log entries for a specific file ID.
func (r *fileLogsRepository) GetByFileID(ctx context.Context, fileID string) (*[]entity.FileLogs, error) {
	var logs *[]entity.FileLogs
	if err := r.db.WithContext(ctx).Where("file_id = ?", fileID).Find(&logs).Error; err != nil {
		slog.Error("Failed to get logs for file_id", slog.Any("file_id", fileID), slog.Any("error", err))
		return nil, fmt.Errorf("failed to get logs for file_id %s: %w", fileID, err)
	}
	return logs, nil
}

// GetByAction retrieves all log entries for a specific action type.
func (r *fileLogsRepository) GetByAction(ctx context.Context, action string) (*[]entity.FileLogs, error) {
	var logs *[]entity.FileLogs
	if err := r.db.WithContext(ctx).Where("action = ?", action).Find(&logs).Error; err != nil {
		slog.Error("Failed to get logs for action", slog.Any("action", action), slog.Any("error", err))
		return nil, fmt.Errorf("failed to get logs for action %s: %w", action, err)
	}
	return logs, nil
}

// DeleteOldLogs removes log entries older than the specified number of days.
// This is useful for log retention management and database cleanup.
func (r *fileLogsRepository) DeleteOldLogs(ctx context.Context, days int) error {
	expiryDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.WithContext(ctx).Where("timestamp < ?", expiryDate).Delete(&entity.FileLogs{}).Error; err != nil {
		slog.Error("Failed to delete old logs", slog.Any("error", err))
		return fmt.Errorf("failed to delete old logs: %w", err)
	}
	return nil
}
