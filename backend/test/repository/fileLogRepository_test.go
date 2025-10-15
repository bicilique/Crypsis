package repository_test

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupFileLogTestDB creates an in-memory SQLite database for testing file logs
func setupFileLogTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.FileLogs{})
	require.NoError(t, err)

	return db
}

// createTestFileLog creates a test file log
func createTestFileLog(t *testing.T, db *gorm.DB, fileID, action string) *entity.FileLogs {
	log := &entity.FileLogs{
		FileID:    fileID,
		Action:    action,
		IP:        "127.0.0.1",
		UserAgent: "Mozilla/5.0",
		Timestamp: time.Now(),
	}
	err := db.Create(log).Error
	require.NoError(t, err)
	return log
}

func TestFileLogRepository_Create(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	t.Run("successfully create file log", func(t *testing.T) {
		log := &entity.FileLogs{
			FileID:    "file-123",
			Action:    "upload",
			IP:        "192.168.1.1",
			UserAgent: "Chrome/91.0",
			Timestamp: time.Now(),
		}

		err := repo.Create(ctx, log)
		assert.NoError(t, err)

		// Verify log was created
		var found entity.FileLogs
		err = db.First(&found, "file_id = ?", log.FileID).Error
		assert.NoError(t, err)
		assert.Equal(t, log.Action, found.Action)
		assert.Equal(t, log.IP, found.IP)
	})

	t.Run("successfully create multiple logs for same file", func(t *testing.T) {
		fileID := "file-456"
		actions := []string{"upload", "download", "delete"}

		for _, action := range actions {
			log := &entity.FileLogs{
				FileID:    fileID,
				Action:    action,
				IP:        "10.0.0.1",
				UserAgent: "Firefox/89.0",
				Timestamp: time.Now(),
			}
			err := repo.Create(ctx, log)
			assert.NoError(t, err)
		}

		// Verify all logs were created
		var logs []entity.FileLogs
		err := db.Where("file_id = ?", fileID).Find(&logs).Error
		assert.NoError(t, err)
		assert.Len(t, logs, 3)
	})
}

func TestFileLogRepository_List(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	// Create test logs
	for i := 0; i < 10; i++ {
		log := &entity.FileLogs{
			FileID:    "file-" + string(rune('a'+i)),
			Action:    "upload",
			IP:        "192.168.1.1",
			UserAgent: "TestAgent",
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
		}
		db.Create(log)
	}

	t.Run("successfully get paginated list", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 0, 5, "timestamp", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, *logs, 5)
	})

	t.Run("successfully get with offset", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 5, 5, "timestamp", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, *logs, 5)
	})

	t.Run("successfully sort by different fields", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 0, 10, "file_id", "asc")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, *logs, 10)
	})

	t.Run("handle default ordering", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 0, 5, "", "")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, *logs, 5)
	})

	t.Run("handle invalid order field", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 0, 5, "invalid_field", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, *logs, 5)
	})

	t.Run("handle invalid sort direction", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 0, 5, "timestamp", "invalid")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, *logs, 5)
	})
}

func TestFileLogRepository_GetByFileID(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	fileID := "test-file-789"

	// Create logs for the file
	actions := []string{"upload", "download", "update", "download", "delete"}
	for _, action := range actions {
		createTestFileLog(t, db, fileID, action)
	}

	// Create logs for other files
	createTestFileLog(t, db, "other-file-1", "upload")
	createTestFileLog(t, db, "other-file-2", "download")

	t.Run("successfully get logs by file ID", func(t *testing.T) {
		logs, err := repo.GetByFileID(ctx, fileID)
		assert.NoError(t, err)
		assert.Len(t, *logs, 5)

		// Verify all logs belong to the correct file
		for _, log := range *logs {
			assert.Equal(t, fileID, log.FileID)
		}
	})

	t.Run("get logs for non-existent file", func(t *testing.T) {
		logs, err := repo.GetByFileID(ctx, "non-existent-file")
		assert.NoError(t, err)
		assert.Empty(t, *logs)
	})
}

func TestFileLogRepository_GetByAction(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	// Create logs with different actions
	actions := map[string]int{
		"upload":   3,
		"download": 5,
		"delete":   2,
		"update":   4,
	}

	for action, count := range actions {
		for i := 0; i < count; i++ {
			createTestFileLog(t, db, "file-"+action+"-"+string(rune('a'+i)), action)
		}
	}

	t.Run("successfully get logs by action upload", func(t *testing.T) {
		logs, err := repo.GetByAction(ctx, "upload")
		assert.NoError(t, err)
		assert.Len(t, *logs, 3)

		// Verify all logs have the correct action
		for _, log := range *logs {
			assert.Equal(t, "upload", log.Action)
		}
	})

	t.Run("successfully get logs by action download", func(t *testing.T) {
		logs, err := repo.GetByAction(ctx, "download")
		assert.NoError(t, err)
		assert.Len(t, *logs, 5)

		for _, log := range *logs {
			assert.Equal(t, "download", log.Action)
		}
	})

	t.Run("successfully get logs by action delete", func(t *testing.T) {
		logs, err := repo.GetByAction(ctx, "delete")
		assert.NoError(t, err)
		assert.Len(t, *logs, 2)

		for _, log := range *logs {
			assert.Equal(t, "delete", log.Action)
		}
	})

	t.Run("get logs for non-existent action", func(t *testing.T) {
		logs, err := repo.GetByAction(ctx, "non-existent-action")
		assert.NoError(t, err)
		assert.Empty(t, *logs)
	})
}

func TestFileLogRepository_DeleteOldLogs(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	now := time.Now()

	// Create logs with different timestamps
	testCases := []struct {
		name      string
		timestamp time.Time
	}{
		{"very-old", now.AddDate(0, 0, -100)},  // 100 days old
		{"old", now.AddDate(0, 0, -50)},        // 50 days old
		{"medium", now.AddDate(0, 0, -20)},     // 20 days old
		{"recent", now.AddDate(0, 0, -5)},      // 5 days old
		{"very-recent", now.AddDate(0, 0, -1)}, // 1 day old
		{"today", now},                         // today
	}

	for _, tc := range testCases {
		log := &entity.FileLogs{
			FileID:    "file-" + tc.name,
			Action:    "upload",
			IP:        "127.0.0.1",
			UserAgent: "TestAgent",
			Timestamp: tc.timestamp,
		}
		db.Create(log)
	}

	t.Run("delete logs older than 30 days", func(t *testing.T) {
		err := repo.DeleteOldLogs(ctx, 30)
		assert.NoError(t, err)

		// Verify old logs were deleted
		var logs []entity.FileLogs
		db.Find(&logs)

		// Should have 4 logs remaining (medium, recent, very-recent, today)
		assert.Len(t, logs, 4)

		// Verify remaining logs are all within 30 days
		for _, log := range logs {
			daysDiff := now.Sub(log.Timestamp).Hours() / 24
			assert.LessOrEqual(t, daysDiff, float64(30))
		}
	})

	// Reset database
	db = setupFileLogTestDB(t)
	repo = repository.NewFileLogRepository(db)

	// Recreate logs
	for _, tc := range testCases {
		log := &entity.FileLogs{
			FileID:    "file-" + tc.name,
			Action:    "upload",
			IP:        "127.0.0.1",
			UserAgent: "TestAgent",
			Timestamp: tc.timestamp,
		}
		db.Create(log)
	}

	t.Run("delete logs older than 10 days", func(t *testing.T) {
		err := repo.DeleteOldLogs(ctx, 10)
		assert.NoError(t, err)

		// Verify old logs were deleted
		var logs []entity.FileLogs
		db.Find(&logs)

		// Should have 3 logs remaining (recent, very-recent, today)
		assert.Len(t, logs, 3)

		// Verify remaining logs are all within 10 days
		for _, log := range logs {
			daysDiff := now.Sub(log.Timestamp).Hours() / 24
			assert.LessOrEqual(t, daysDiff, float64(10))
		}
	})
}

func TestFileLogRepository_ConcurrentWrites(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	t.Run("handle concurrent log writes", func(t *testing.T) {
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				log := &entity.FileLogs{
					FileID:    "concurrent-file",
					Action:    "upload",
					IP:        "192.168.1." + string(rune('0'+index)),
					UserAgent: "Goroutine-" + string(rune('0'+index)),
					Timestamp: time.Now(),
				}
				err := repo.Create(ctx, log)
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Verify all logs were created
		var logs []entity.FileLogs
		err := db.Where("file_id = ?", "concurrent-file").Find(&logs).Error
		assert.NoError(t, err)
		assert.Len(t, logs, numGoroutines)
	})
}

func TestFileLogRepository_EmptyDatabase(t *testing.T) {
	db := setupFileLogTestDB(t)
	repo := repository.NewFileLogRepository(db)
	ctx := context.Background()

	t.Run("list on empty database", func(t *testing.T) {
		total, logs, err := repo.List(ctx, 0, 10, "timestamp", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, *logs)
	})

	t.Run("get by file ID on empty database", func(t *testing.T) {
		logs, err := repo.GetByFileID(ctx, "non-existent")
		assert.NoError(t, err)
		assert.Empty(t, *logs)
	})

	t.Run("get by action on empty database", func(t *testing.T) {
		logs, err := repo.GetByAction(ctx, "upload")
		assert.NoError(t, err)
		assert.Empty(t, *logs)
	})

	t.Run("delete old logs on empty database", func(t *testing.T) {
		err := repo.DeleteOldLogs(ctx, 30)
		assert.NoError(t, err)
	})
}
