package repository

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/model"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

// fileRepository implements the FileRepository interface for file data access operations.
type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new instance of FileRepository.
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

// Create adds a new file record to the database.
func (r *fileRepository) Create(ctx context.Context, file *entity.Files) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		slog.Error("Failed to create file with ID "+file.ID, slog.Any("error", err))
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

// GetByID retrieves a file by its unique ID.
func (r *fileRepository) GetByID(ctx context.Context, id string) (*entity.Files, error) {
	if id == "" {
		return nil, errors.New("file ID cannot be empty")
	}
	var file entity.Files
	if err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to retrieve file by ID: %w", err)
	}
	return &file, nil
}

// GetByHash retrieves a file by its hash value.
func (r *fileRepository) GetByHash(ctx context.Context, hash string) (*entity.Files, error) {
	if hash == "" {
		return nil, errors.New("hash cannot be empty")
	}
	var file entity.Files
	if err := r.db.WithContext(ctx).First(&file, "hash = ?", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to retrieve file by hash: %w", err)
	}
	return &file, nil
}

// GetListFiles retrieves a paginated list of files for a specific application.
func (r *fileRepository) GetListFiles(ctx context.Context, appID string, offset, limit int, orderBy, sort string) (int64, []entity.Files, error) {
	return r.getListFiles(ctx, map[string]interface{}{"app_id": appID}, offset, limit, orderBy, sort)
}

// GetListFilesForAdmin retrieves a paginated list of all files for admin users.
func (r *fileRepository) GetListFilesForAdmin(ctx context.Context, offset, limit int, orderBy, sort string) (int64, []entity.Files, error) {
	return r.getListFiles(ctx, nil, offset, limit, orderBy, sort)
}

// GetAll retrieves all files from the database.
func (r *fileRepository) GetAll(ctx context.Context) ([]entity.Files, error) {
	var files []entity.Files
	if err := r.db.WithContext(ctx).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve all files: %w", err)
	}
	return files, nil
}

// Update modifies an existing file record in the database.
func (r *fileRepository) Update(ctx context.Context, file *entity.Files) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}
	if file.ID == "" {
		return errors.New("file ID cannot be empty")
	}
	if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

// Delete performs a soft delete on a file and its associated metadata.
func (r *fileRepository) Delete(ctx context.Context, fileID string) error {
	if fileID == "" {
		return errors.New("file ID cannot be empty")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Soft delete metadata first
		result := tx.Where("file_id = ?", fileID).Delete(&entity.Metadata{})
		if result.Error != nil {
			slog.Error("Failed to delete metadata", slog.String("fileID", fileID), slog.Any("error", result.Error))
			return fmt.Errorf("failed to delete metadata: %w", result.Error)
		}

		// Soft delete file
		result = tx.Where("id = ?", fileID).Delete(&entity.Files{})
		if result.Error != nil {
			slog.Error("Failed to delete file", slog.String("fileID", fileID), slog.Any("error", result.Error))
			return fmt.Errorf("failed to delete file: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return model.ErrFileNotFound
		}
		return nil
	})

}

// WithTransaction executes a function within a database transaction.
func (r *fileRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// CreateFileWithMetadata creates a new file and its metadata in a single transaction.
func (r *fileRepository) CreateFileWithMetadata(ctx context.Context, file *entity.Files, metadata *entity.Metadata) error {
	if file == nil || metadata == nil {
		return errors.New("file and metadata cannot be nil")
	}
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(file).Error; err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		if err := tx.Create(metadata).Error; err != nil {
			return fmt.Errorf("failed to create metadata: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetMetadataByFileID retrieves metadata for a file by its file ID.
func (r *fileRepository) GetMetadataByFileID(ctx context.Context, fileID string) (*entity.Metadata, error) {
	if fileID == "" {
		return nil, errors.New("file ID cannot be empty")
	}
	var metadata entity.Metadata
	if err := r.db.WithContext(ctx).
		Preload("File"). // Preload the related File
		Where("file_id = ?", fileID).
		First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to retrieve metadata: %w", err)
	}

	return &metadata, nil
}

// GetMetadataByAppIDAndFileID retrieves metadata by both application ID and file ID.
func (r *fileRepository) GetMetadataByAppIDAndFileID(ctx context.Context, appID, fileID string) (*entity.Metadata, error) {
	if appID == "" || fileID == "" {
		return nil, errors.New("appID and fileID cannot be empty")
	}
	var metadata entity.Metadata

	if err := r.db.WithContext(ctx).
		Joins("JOIN files ON files.id = metadata.file_id").
		Where("metadata.file_id = ? AND files.app_id = ?", fileID, appID).
		Preload("File").
		First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to retrieve file information: %w", err)
	}

	return &metadata, nil
}

// GetDeletedMetadataByAppIDAndFileID retrieves soft-deleted metadata by application ID and file ID.
func (r *fileRepository) GetDeletedMetadataByAppIDAndFileID(ctx context.Context, appID, fileID string) (*entity.Metadata, error) {
	if appID == "" || fileID == "" {
		return nil, errors.New("appID and fileID cannot be empty")
	}
	var metadata entity.Metadata
	if err := r.db.WithContext(ctx).
		Unscoped().
		Joins("JOIN files ON files.id = metadata.file_id").
		Where("metadata.file_id = ? AND files.app_id = ?", fileID, appID).
		Preload("File").
		First(&metadata).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to retrieve file information: %w", err)
	}

	return &metadata, nil
}

// GetMetadataByEncHash retrieves metadata by its encrypted hash value.
func (r *fileRepository) GetMetadataByEncHash(ctx context.Context, encHash string) (*entity.Metadata, error) {
	if encHash == "" {
		return nil, errors.New("encrypted hash cannot be empty")
	}
	var metadata entity.Metadata
	if err := r.db.WithContext(ctx).First(&metadata, "enc_hash = ?", encHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("metadata not found")
		}
		return nil, fmt.Errorf("failed to retrieve metadata by enc_hash: %w", err)
	}
	return &metadata, nil
}

// GetAllMetadata retrieves all metadata records from the database.
func (r *fileRepository) GetAllMetadata(ctx context.Context) ([]entity.Metadata, error) {
	var metadata []entity.Metadata
	if err := r.db.WithContext(ctx).Find(&metadata).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve all metadata: %w", err)
	}
	return metadata, nil
}

// RestoreFile restores a soft-deleted file and its metadata.
func (r *fileRepository) RestoreFile(ctx context.Context, fileID string) error {
	if fileID == "" {
		return errors.New("file ID cannot be empty")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Restore metadata first
		result := tx.Model(&entity.Metadata{}).Unscoped().Where("file_id = ?", fileID).Update("deleted_at", nil)
		if result.Error != nil {
			slog.Error("Failed to restore metadata", slog.String("fileID", fileID), slog.Any("error", result.Error))
			return fmt.Errorf("failed to restore metadata: %w", result.Error)
		}

		// Restore file
		result = tx.Model(&entity.Files{}).Unscoped().Where("id = ?", fileID).Update("deleted_at", nil)
		if result.Error != nil {
			slog.Error("Failed to restore file", slog.String("fileID", fileID), slog.Any("error", result.Error))
			return fmt.Errorf("failed to restore file: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("file not found or not deleted")
		}

		return nil
	})
}

// UpdateFileAndMetadata updates both file and metadata records in a single transaction.
func (r *fileRepository) UpdateFileAndMetadata(ctx context.Context, file *entity.Files, metadata *entity.Metadata) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update file
		if file != nil {
			if file.ID == "" {
				return errors.New("file ID cannot be empty")
			}
			result := tx.Model(&entity.Files{}).Where("id = ?", file.ID).Updates(file)
			if result.Error != nil {
				slog.Error("Failed to update file", slog.String("fileID", file.ID), slog.Any("error", result.Error))
				return fmt.Errorf("failed to update file: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return errors.New("file not found or no changes made")
			}
		}

		// Update metadata
		if metadata != nil {
			if metadata.FileID == "" {
				return errors.New("file ID in metadata cannot be empty")
			}
			result := tx.Model(&entity.Metadata{}).Where("file_id = ?", metadata.FileID).Updates(metadata)
			if result.Error != nil {
				slog.Error("Failed to update metadata", slog.String("fileID", metadata.FileID), slog.Any("error", result.Error))
				return fmt.Errorf("failed to update metadata: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return errors.New("metadata not found or no changes made")
			}
		}

		return nil
	})
}

// GetAllKeyUIDs retrieves all key UIDs from metadata records.
func (r *fileRepository) GetAllKeyUIDs(ctx context.Context) ([]string, error) {
	var keyUIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.Metadata{}).
		Where("key_uid IS NOT NULL").
		Pluck("key_uid", &keyUIDs).Error; err != nil {
		return nil, errors.New("failed to retrieve key_uids: " + err.Error())
	}
	return keyUIDs, nil
}

// BatchUpdateEncKeys performs a batch update of encryption keys for multiple key UIDs.
// Uses raw SQL for better performance compared to individual updates.
func (r *fileRepository) BatchUpdateEncKeys(ctx context.Context, updates map[string]string) error {
	if len(updates) == 0 {
		return nil // nothing to do
	}

	var caseSQL strings.Builder
	var args []interface{}
	var uids []interface{}

	caseSQL.WriteString("UPDATE metadata SET enc_key = CASE key_uid\n")

	for uid, encKey := range updates {
		caseSQL.WriteString("WHEN ? THEN ?\n")
		args = append(args, uid, encKey)
		uids = append(uids, uid)
	}

	caseSQL.WriteString("ELSE enc_key END\n")
	caseSQL.WriteString("WHERE key_uid IN (")

	placeholders := make([]string, len(uids))
	for i := range uids {
		placeholders[i] = "?"
	}
	caseSQL.WriteString(strings.Join(placeholders, ","))
	caseSQL.WriteString(");")

	args = append(args, uids...)

	// Execute the raw SQL
	if err := r.db.WithContext(ctx).Exec(caseSQL.String(), args...).Error; err != nil {
		return fmt.Errorf("batch update failed: %w", err)
	}

	return nil
}

// UpdateEncKeyByKeyUID updates the encryption key for a specific key UID.
func (r *fileRepository) UpdateEncKeyByKeyUID(ctx context.Context, keyUID string, newEncKey string) error {
	if keyUID == "" || newEncKey == "" {
		return fmt.Errorf("keyUID and newEncKey cannot be empty")
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Metadata{}).
		Where("key_uid = ?", keyUID).
		Update("enc_key", newEncKey).Error

	if err != nil {
		return fmt.Errorf("failed to update enc_key for key_uid %s: %w", keyUID, err)
	}

	return nil
}

// getListFiles is a helper function that retrieves paginated file lists with optional filtering.
// It supports both application-specific and admin queries.
func (r *fileRepository) getListFiles(
	ctx context.Context,
	filter map[string]interface{},
	offset, limit int,
	orderBy, sort string,
) (int64, []entity.Files, error) {
	var total int64
	files := make([]entity.Files, 0)

	// Default handling
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}
	allowedOrderFields := map[string]bool{
		"created_at": true,
		"name":       true,
		"size":       true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "created_at"
	}

	// Apply filters
	query := r.db.WithContext(ctx).Model(&entity.Files{})
	if filter != nil {
		query = query.Where(filter)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count files: %w", err)
	}

	if filter != nil {
		if err := r.db.WithContext(ctx).
			Model(&entity.Files{}).
			Where("app_id = ?", filter["app_id"]).
			Offset(offset).
			Limit(limit).
			Order(fmt.Sprintf("%s %s", orderBy, sort)).
			Find(&files).Error; err != nil {
			return 0, nil, fmt.Errorf("failed to get list of files: %w", err)
		}
	} else {
		if err := r.db.WithContext(ctx).
			Model(&entity.Files{}).
			Offset(offset).
			Limit(limit).
			Unscoped().
			Order(fmt.Sprintf("%s %s", orderBy, sort)).
			Find(&files).Error; err != nil {
			return 0, nil, fmt.Errorf("failed to get list of files: %w", err)
		}
	}
	return total, files, nil
}
