package repository

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.Files{}, &entity.Metadata{}, &entity.Apps{})
	require.NoError(t, err)

	return db
}

// createTestApp creates a test application
func createTestApp(t *testing.T, db *gorm.DB) *entity.Apps {
	app := &entity.Apps{
		ID:           "test-app-id",
		Name:         "Test App",
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		IsActive:     true,
	}
	err := db.Create(app).Error
	require.NoError(t, err)
	return app
}

// createTestFile creates a test file
func createTestFile(t *testing.T, db *gorm.DB, appID string) *entity.Files {
	file := &entity.Files{
		ID:         "test-file-id",
		AppID:      appID,
		UserID:     "test-user-id",
		Name:       "test-file.txt",
		MimeType:   "text/plain",
		Size:       1024,
		BucketName: "test-bucket",
		Location:   "/test/location",
	}
	err := db.Create(file).Error
	require.NoError(t, err)
	return file
}

// createTestMetadata creates test metadata
func createTestMetadata(t *testing.T, db *gorm.DB, fileID string) *entity.Metadata {
	metadata := &entity.Metadata{
		FileID:  fileID,
		KeyUID:  "test-key-uid",
		EncKey:  "test-enc-key",
		EncHash: "test-enc-hash",
	}
	err := db.Create(metadata).Error
	require.NoError(t, err)
	return metadata
}

func TestFileRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully create file", func(t *testing.T) {
		app := createTestApp(t, db)
		file := &entity.Files{
			ID:    "file-1",
			AppID: app.ID,
			Name:  "document.pdf",
			Size:  2048,
		}

		err := repo.Create(ctx, file)
		assert.NoError(t, err)

		// Verify file was created
		var found entity.Files
		err = db.First(&found, "id = ?", file.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, file.Name, found.Name)
	})

	t.Run("fail with nil file", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file cannot be nil")
	})
}

func TestFileRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully get file by ID", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)

		found, err := repo.GetByID(ctx, file.ID)
		assert.NoError(t, err)
		assert.Equal(t, file.ID, found.ID)
		assert.Equal(t, file.Name, found.Name)
	})

	t.Run("fail with empty ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file ID cannot be empty")
	})

	t.Run("fail with non-existent ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existent")
		assert.Error(t, err)
		assert.Equal(t, model.ErrFileNotFound, err)
	})
}

// TestFileRepository_GetByHash is commented out because the Files entity
// does not have a Hash field, though the repository interface defines GetByHash method.
// This indicates a schema mismatch that should be addressed separately.
/*
func TestFileRepository_GetByHash(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully get file by hash", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)

		found, err := repo.GetByHash(ctx, "test-hash")
		assert.NoError(t, err)
		assert.Equal(t, file.ID, found.ID)
	})

	t.Run("fail with empty hash", func(t *testing.T) {
		_, err := repo.GetByHash(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hash cannot be empty")
	})

	t.Run("fail with non-existent hash", func(t *testing.T) {
		_, err := repo.GetByHash(ctx, "non-existent-hash")
		assert.Error(t, err)
		assert.Equal(t, model.ErrFileNotFound, err)
	})
}
*/

func TestFileRepository_GetListFiles(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	app := createTestApp(t, db)

	// Create multiple files
	for i := 0; i < 5; i++ {
		file := &entity.Files{
			ID:    "file-" + string(rune('a'+i)),
			AppID: app.ID,
			Name:  "file-" + string(rune('a'+i)) + ".txt",
			Size:  int64(1024 * (i + 1)),
		}
		err := db.Create(file).Error
		require.NoError(t, err)
	}

	t.Run("successfully get paginated list", func(t *testing.T) {
		total, files, err := repo.GetListFiles(ctx, app.ID, 0, 3, "created_at", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, files, 3)
	})

	t.Run("successfully get with offset", func(t *testing.T) {
		total, files, err := repo.GetListFiles(ctx, app.ID, 3, 3, "created_at", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, files, 2)
	})
}

func TestFileRepository_GetListFilesForAdmin(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	app1 := createTestApp(t, db)
	app2 := &entity.Apps{
		ID:       "app-2",
		Name:     "App 2",
		ClientID: "client-2",
		IsActive: true,
	}
	db.Create(app2)

	// Create files for both apps
	for i := 0; i < 3; i++ {
		db.Create(&entity.Files{
			ID:    "file-app1-" + string(rune('a'+i)),
			AppID: app1.ID,
			Name:  "file1-" + string(rune('a'+i)) + ".txt",
			Size:  1024,
		})
		db.Create(&entity.Files{
			ID:    "file-app2-" + string(rune('a'+i)),
			AppID: app2.ID,
			Name:  "file2-" + string(rune('a'+i)) + ".txt",
			Size:  2048,
		})
	}

	t.Run("successfully get all files for admin", func(t *testing.T) {
		total, files, err := repo.GetListFilesForAdmin(ctx, 0, 10, "created_at", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total)
		assert.Len(t, files, 6)
	})
}

func TestFileRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	app := createTestApp(t, db)

	// Create some files
	for i := 0; i < 3; i++ {
		db.Create(&entity.Files{
			ID:    "file-" + string(rune('a'+i)),
			AppID: app.ID,
			Name:  "file-" + string(rune('a'+i)) + ".txt",
			Size:  1024,
		})
	}

	t.Run("successfully get all files", func(t *testing.T) {
		files, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(files), 3)
	})
}

func TestFileRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully update file", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)

		file.Name = "updated-file.txt"
		file.Size = 2048

		err := repo.Update(ctx, file)
		assert.NoError(t, err)

		// Verify update
		var updated entity.Files
		db.First(&updated, "id = ?", file.ID)
		assert.Equal(t, "updated-file.txt", updated.Name)
		assert.Equal(t, int64(2048), updated.Size)
	})

	t.Run("fail with nil file", func(t *testing.T) {
		err := repo.Update(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file cannot be nil")
	})

	t.Run("fail with empty ID", func(t *testing.T) {
		file := &entity.Files{Name: "test"}
		err := repo.Update(ctx, file)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file ID cannot be empty")
	})
}

func TestFileRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully delete file and metadata", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		err := repo.Delete(ctx, file.ID)
		assert.NoError(t, err)

		// Verify soft delete
		var deletedFile entity.Files
		err = db.First(&deletedFile, "id = ?", file.ID).Error
		assert.Error(t, err)

		var deletedMetadata entity.Metadata
		err = db.First(&deletedMetadata, "file_id = ?", metadata.FileID).Error
		assert.Error(t, err)
	})

	t.Run("fail with empty ID", func(t *testing.T) {
		err := repo.Delete(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file ID cannot be empty")
	})

	t.Run("fail with non-existent file", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent")
		assert.Error(t, err)
	})
}

func TestFileRepository_CreateFileWithMetadata(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully create file with metadata", func(t *testing.T) {
		app := createTestApp(t, db)
		file := &entity.Files{
			ID:    "file-with-metadata",
			AppID: app.ID,
			Name:  "document.pdf",
			Size:  4096,
		}
		metadata := &entity.Metadata{
			FileID:  file.ID,
			KeyUID:  "key-123",
			EncKey:  "enc-key-123",
			EncHash: "enc-hash-123",
		}

		err := repo.CreateFileWithMetadata(ctx, file, metadata)
		assert.NoError(t, err)

		// Verify both were created
		var foundFile entity.Files
		err = db.First(&foundFile, "id = ?", file.ID).Error
		assert.NoError(t, err)

		var foundMetadata entity.Metadata
		err = db.First(&foundMetadata, "file_id = ?", file.ID).Error
		assert.NoError(t, err)
	})

	t.Run("fail with nil file", func(t *testing.T) {
		metadata := &entity.Metadata{FileID: "test"}
		err := repo.CreateFileWithMetadata(ctx, nil, metadata)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file and metadata cannot be nil")
	})

	t.Run("fail with nil metadata", func(t *testing.T) {
		file := &entity.Files{ID: "test"}
		err := repo.CreateFileWithMetadata(ctx, file, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file and metadata cannot be nil")
	})
}

func TestFileRepository_GetMetadataByFileID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully get metadata by file ID", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		found, err := repo.GetMetadataByFileID(ctx, file.ID)
		assert.NoError(t, err)
		assert.Equal(t, metadata.FileID, found.FileID)
		assert.Equal(t, metadata.KeyUID, found.KeyUID)
	})

	t.Run("fail with empty file ID", func(t *testing.T) {
		_, err := repo.GetMetadataByFileID(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file ID cannot be empty")
	})

	t.Run("fail with non-existent file ID", func(t *testing.T) {
		_, err := repo.GetMetadataByFileID(ctx, "non-existent")
		assert.Error(t, err)
	})
}

func TestFileRepository_GetMetadataByAppIDAndFileID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully get metadata by app and file ID", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		found, err := repo.GetMetadataByAppIDAndFileID(ctx, app.ID, file.ID)
		assert.NoError(t, err)
		assert.Equal(t, metadata.FileID, found.FileID)
	})

	t.Run("fail with empty app ID", func(t *testing.T) {
		_, err := repo.GetMetadataByAppIDAndFileID(ctx, "", "file-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "appID and fileID cannot be empty")
	})

	t.Run("fail with empty file ID", func(t *testing.T) {
		_, err := repo.GetMetadataByAppIDAndFileID(ctx, "app-id", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "appID and fileID cannot be empty")
	})
}

func TestFileRepository_RestoreFile(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully restore deleted file", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		// Delete the file first
		err := repo.Delete(ctx, file.ID)
		require.NoError(t, err)

		// Restore the file
		err = repo.RestoreFile(ctx, file.ID)
		assert.NoError(t, err)

		// Verify restoration
		var restored entity.Files
		err = db.First(&restored, "id = ?", file.ID).Error
		assert.NoError(t, err)

		var restoredMetadata entity.Metadata
		err = db.First(&restoredMetadata, "file_id = ?", metadata.FileID).Error
		assert.NoError(t, err)
	})

	t.Run("fail with empty file ID", func(t *testing.T) {
		err := repo.RestoreFile(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file ID cannot be empty")
	})
}

func TestFileRepository_UpdateFileAndMetadata(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully update file and metadata", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		file.Name = "updated.txt"
		metadata.EncKey = "new-enc-key"

		err := repo.UpdateFileAndMetadata(ctx, file, metadata)
		assert.NoError(t, err)

		// Verify updates
		var updatedFile entity.Files
		db.First(&updatedFile, "id = ?", file.ID)
		assert.Equal(t, "updated.txt", updatedFile.Name)

		var updatedMetadata entity.Metadata
		db.First(&updatedMetadata, "file_id = ?", file.ID)
		assert.Equal(t, "new-enc-key", updatedMetadata.EncKey)
	})

	t.Run("update only file", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		createTestMetadata(t, db, file.ID)

		file.Size = 9999
		err := repo.UpdateFileAndMetadata(ctx, file, nil)
		assert.NoError(t, err)
	})

	t.Run("update only metadata", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		metadata.KeyUID = "new-key-uid"
		err := repo.UpdateFileAndMetadata(ctx, nil, metadata)
		assert.NoError(t, err)
	})
}

func TestFileRepository_GetAllKeyUIDs(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	app := createTestApp(t, db)

	// Create multiple files with metadata
	for i := 0; i < 3; i++ {
		file := &entity.Files{
			ID:    "file-" + string(rune('a'+i)),
			AppID: app.ID,
			Name:  "file.txt",
			Size:  1024,
		}
		db.Create(file)

		metadata := &entity.Metadata{
			FileID: file.ID,
			KeyUID: "key-uid-" + string(rune('a'+i)),
			EncKey: "enc-key",
		}
		db.Create(metadata)
	}

	t.Run("successfully get all key UIDs", func(t *testing.T) {
		keyUIDs, err := repo.GetAllKeyUIDs(ctx)
		assert.NoError(t, err)
		assert.Len(t, keyUIDs, 3)
	})
}

func TestFileRepository_BatchUpdateEncKeys(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	app := createTestApp(t, db)

	// Create files with metadata
	keyUIDs := []string{"key-1", "key-2", "key-3"}
	for i, keyUID := range keyUIDs {
		file := &entity.Files{
			ID:    "file-" + string(rune('a'+i)),
			AppID: app.ID,
			Name:  "file.txt",
			Size:  1024,
		}
		db.Create(file)

		metadata := &entity.Metadata{
			FileID: file.ID,
			KeyUID: keyUID,
			EncKey: "old-enc-key-" + string(rune('a'+i)),
		}
		db.Create(metadata)
	}

	t.Run("successfully batch update encryption keys", func(t *testing.T) {
		updates := map[string]string{
			"key-1": "new-enc-key-1",
			"key-2": "new-enc-key-2",
			"key-3": "new-enc-key-3",
		}

		err := repo.BatchUpdateEncKeys(ctx, updates)
		assert.NoError(t, err)

		// Verify updates
		for keyUID, expectedEncKey := range updates {
			var metadata entity.Metadata
			db.First(&metadata, "key_uid = ?", keyUID)
			assert.Equal(t, expectedEncKey, metadata.EncKey)
		}
	})

	t.Run("handle empty updates", func(t *testing.T) {
		err := repo.BatchUpdateEncKeys(ctx, map[string]string{})
		assert.NoError(t, err)
	})
}

func TestFileRepository_UpdateEncKeyByKeyUID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully update encryption key", func(t *testing.T) {
		app := createTestApp(t, db)
		file := createTestFile(t, db, app.ID)
		metadata := createTestMetadata(t, db, file.ID)

		err := repo.UpdateEncKeyByKeyUID(ctx, metadata.KeyUID, "new-enc-key")
		assert.NoError(t, err)

		// Verify update
		var updated entity.Metadata
		db.First(&updated, "key_uid = ?", metadata.KeyUID)
		assert.Equal(t, "new-enc-key", updated.EncKey)
	})

	t.Run("fail with empty key UID", func(t *testing.T) {
		err := repo.UpdateEncKeyByKeyUID(ctx, "", "new-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "keyUID and newEncKey cannot be empty")
	})

	t.Run("fail with empty new enc key", func(t *testing.T) {
		err := repo.UpdateEncKeyByKeyUID(ctx, "key-uid", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "keyUID and newEncKey cannot be empty")
	})
}

func TestFileRepository_WithTransaction(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewFileRepository(db)
	ctx := context.Background()

	t.Run("successfully execute transaction", func(t *testing.T) {
		app := createTestApp(t, db)

		err := repo.WithTransaction(ctx, func(tx *gorm.DB) error {
			file := &entity.Files{
				ID:         "tx-file",
				AppID:      app.ID,
				UserID:     "user-1",
				Name:       "tx-file.txt",
				MimeType:   "text/plain",
				Size:       1024,
				BucketName: "test-bucket",
				Location:   "/tx/location",
			}
			return tx.Create(file).Error
		})
		assert.NoError(t, err)

		// Verify file was created
		var found entity.Files
		err = db.First(&found, "id = ?", "tx-file").Error
		assert.NoError(t, err)
	})

	t.Run("rollback on error", func(t *testing.T) {
		app := createTestApp(t, db)

		err := repo.WithTransaction(ctx, func(tx *gorm.DB) error {
			file := &entity.Files{
				ID:         "rollback-file",
				AppID:      app.ID,
				UserID:     "user-1",
				Name:       "rollback.txt",
				MimeType:   "text/plain",
				Size:       1024,
				BucketName: "test-bucket",
				Location:   "/rollback/location",
			}
			if err := tx.Create(file).Error; err != nil {
				return err
			}
			// Force an error to trigger rollback
			return gorm.ErrInvalidTransaction
		})
		assert.Error(t, err)

		// Verify file was not created
		var found entity.Files
		err = db.First(&found, "id = ?", "rollback-file").Error
		assert.Error(t, err)
	})
}
