package repository

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/model"
	"fmt"

	"gorm.io/gorm"
)

// appsRepository implements the AppsRepository interface for application data access operations.
type appsRepository struct {
	db *gorm.DB
}

// NewAppsRepository creates a new instance of AppsRepository.
func NewAppsRepository(db *gorm.DB) ApplicationRepository {
	return &appsRepository{db: db}
}

// Create adds a new application record to the database.
func (r *appsRepository) Create(ctx context.Context, app *entity.Apps) error {
	if err := r.db.WithContext(ctx).Create(app).Error; err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}
	return nil
}

// GetByID retrieves an application by its unique ID, including soft-deleted records.
func (r *appsRepository) GetByID(ctx context.Context, id string) (*entity.Apps, error) {
	var app entity.Apps
	if err := r.db.WithContext(ctx).
		Unscoped().
		First(&app, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("app not found")
		}
		return nil, fmt.Errorf("failed to get app by id: %w", err)
	}
	return &app, nil
}

// GetByClientID retrieves an application by its OAuth2 client ID.
func (r *appsRepository) GetByClientID(ctx context.Context, clientID string) (*entity.Apps, error) {
	var app entity.Apps
	if err := r.db.WithContext(ctx).
		First(&app, "client_id = ?", clientID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrAppNotFound
		}
		return nil, fmt.Errorf("failed to get app by client_id: %w", err)
	}
	return &app, nil
}

// GetByClientIDLimited retrieves limited application data (id and is_active) by client ID.
// This is optimized for quick validation checks.
func (r *appsRepository) GetByClientIDLimited(ctx context.Context, clientID string) (*entity.Apps, error) {
	var app entity.Apps
	if err := r.db.WithContext(ctx).
		Select("id", "is_active").
		Where("client_id = ?", clientID).
		First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrAppNotFound
		}
		return nil, fmt.Errorf("failed to get app by client_id: %w", err)
	}
	return &app, nil
}

// GetByName retrieves an application by its name.
func (r *appsRepository) GetByName(ctx context.Context, name string) (*entity.Apps, error) {
	var app entity.Apps
	if err := r.db.WithContext(ctx).First(&app, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("app not found")
		}
		return nil, fmt.Errorf("failed to get app by name: %w", err)
	}
	return &app, nil
}

// Update modifies an existing application record in the database.
func (r *appsRepository) Update(ctx context.Context, app *entity.Apps) error {
	if err := r.db.WithContext(ctx).Save(app).Error; err != nil {
		return fmt.Errorf("failed to update app: %w", err)
	}
	return nil
}

// Delete deactivates an application by setting is_active to false.
// This performs a logical delete rather than removing the record.
func (r *appsRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Apps{}).
		Where("id = ?", id).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to deactivate app: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("app not found")
	}

	return nil
}

// Restore recovers a soft-deleted application and reactivates it.
func (r *appsRepository) Restore(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Apps{}).
		Unscoped().
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": nil,
			"is_active":  true,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to restore app: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("app not found")
	}
	return nil
}

// GetAll retrieves all application records from the database.
func (r *appsRepository) GetAll(ctx context.Context) ([]entity.Apps, error) {
	var apps []entity.Apps
	if err := r.db.WithContext(ctx).Find(&apps).Error; err != nil {
		return nil, fmt.Errorf("failed to get all apps: %w", err)
	}
	return apps, nil
}

// GetListApps retrieves a paginated list of applications with sorting options.
// Includes soft-deleted records for admin visibility.
func (r *appsRepository) GetListApps(ctx context.Context, offset, limit int, orderBy, sort string) (int64, *[]entity.Apps, error) {
	var total int64
	var apps []entity.Apps

	// Default sorting
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	// Validate allowed fields
	allowedOrderFields := map[string]bool{
		"created_at": true,
		"name":       true,
		"client_id":  true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "created_at"
	}

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&entity.Apps{}).
		Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count apps: %w", err)
	}

	// Get paginated list
	if err := r.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", orderBy, sort)).
		Offset(offset).
		Unscoped().
		Limit(limit).
		Find(&apps).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to get list of apps: %w", err)
	}
	return total, &apps, nil
}
