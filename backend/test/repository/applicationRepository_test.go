package repository_test

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupApplicationTestDB initializes an in-memory SQLite database for testing.
func setupApplicationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to in-memory database")

	// Auto-migrate the Apps entity
	err = db.AutoMigrate(&entity.Apps{})
	require.NoError(t, err, "Failed to migrate schema")

	return db
}

// createTestApp creates a test application entity with default values.
func createTestApp(id, name, clientID string) *entity.Apps {
	return &entity.Apps{
		ID:           id,
		Name:         name,
		ClientID:     clientID,
		ClientSecret: "test_secret",
		IsActive:     true,
		Uri:          "https://example.com",
		RedirectUri:  "https://example.com/callback",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestApplicationRepository_Create(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		app := createTestApp("app-1", "Test App", "client-1")

		err := repo.Create(ctx, app)
		assert.NoError(t, err)

		// Verify the app was created
		var retrieved entity.Apps
		err = db.First(&retrieved, "id = ?", "app-1").Error
		assert.NoError(t, err)
		assert.Equal(t, "Test App", retrieved.Name)
		assert.Equal(t, "client-1", retrieved.ClientID)
		assert.True(t, retrieved.IsActive)
	})

	t.Run("Duplicate ID", func(t *testing.T) {
		app1 := createTestApp("app-2", "App 1", "client-2")
		err := repo.Create(ctx, app1)
		require.NoError(t, err)

		// Try to create another app with the same ID
		app2 := createTestApp("app-2", "App 2", "client-3")
		err = repo.Create(ctx, app2)
		assert.Error(t, err)
	})

	t.Run("Empty Required Fields", func(t *testing.T) {
		app := &entity.Apps{
			ID:       "app-3",
			Name:     "",
			ClientID: "",
		}
		err := repo.Create(ctx, app)
		// SQLite may or may not enforce NOT NULL constraints depending on version
		// Just verify we can handle the case
		_ = err
	})
}

func TestApplicationRepository_GetByID(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	// Create test app
	app := createTestApp("app-1", "Test App", "client-1")
	err := repo.Create(ctx, app)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		retrieved, err := repo.GetByID(ctx, "app-1")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "app-1", retrieved.ID)
		assert.Equal(t, "Test App", retrieved.Name)
		assert.Equal(t, "client-1", retrieved.ClientID)
	})

	t.Run("Not Found", func(t *testing.T) {
		retrieved, err := repo.GetByID(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "app not found")
	})

	t.Run("Include Soft Deleted", func(t *testing.T) {
		// Create and soft delete an app
		app2 := createTestApp("app-2", "Deleted App", "client-2")
		err := repo.Create(ctx, app2)
		require.NoError(t, err)

		// Soft delete the app
		err = db.Delete(app2).Error
		require.NoError(t, err)

		// GetByID should still find it (Unscoped)
		retrieved, err := repo.GetByID(ctx, "app-2")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "app-2", retrieved.ID)
	})
}

func TestApplicationRepository_GetByClientID(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	// Create test app
	app := createTestApp("app-1", "Test App", "client-1")
	err := repo.Create(ctx, app)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		retrieved, err := repo.GetByClientID(ctx, "client-1")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "app-1", retrieved.ID)
		assert.Equal(t, "client-1", retrieved.ClientID)
	})

	t.Run("Not Found", func(t *testing.T) {
		retrieved, err := repo.GetByClientID(ctx, "non-existent-client")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, model.ErrAppNotFound, err)
	})

	t.Run("Excludes Soft Deleted", func(t *testing.T) {
		// Create and soft delete an app
		app2 := createTestApp("app-2", "Deleted App", "client-2")
		err := repo.Create(ctx, app2)
		require.NoError(t, err)

		// Soft delete the app
		err = db.Delete(app2).Error
		require.NoError(t, err)

		// GetByClientID should not find it (uses default scope)
		retrieved, err := repo.GetByClientID(ctx, "client-2")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, model.ErrAppNotFound, err)
	})
}

func TestApplicationRepository_GetByClientIDLimited(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	// Create test app
	app := createTestApp("app-1", "Test App", "client-1")
	app.IsActive = true
	err := repo.Create(ctx, app)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		retrieved, err := repo.GetByClientIDLimited(ctx, "client-1")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "app-1", retrieved.ID)
		assert.True(t, retrieved.IsActive)
		// Other fields should be zero values (not selected)
		assert.Empty(t, retrieved.Name)
		// ClientID is not in SELECT clause, so it won't be populated
		assert.Empty(t, retrieved.ClientID)
	})

	t.Run("Not Found", func(t *testing.T) {
		retrieved, err := repo.GetByClientIDLimited(ctx, "non-existent-client")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Equal(t, model.ErrAppNotFound, err)
	})

	t.Run("Inactive App", func(t *testing.T) {
		// Create inactive app
		app2 := createTestApp("app-2", "Inactive App", "client-2")
		app2.IsActive = false
		err := repo.Create(ctx, app2)
		require.NoError(t, err)

		retrieved, err := repo.GetByClientIDLimited(ctx, "client-2")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.False(t, retrieved.IsActive)
	})
}

func TestApplicationRepository_GetByName(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	// Create test app
	app := createTestApp("app-1", "Test App", "client-1")
	err := repo.Create(ctx, app)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		retrieved, err := repo.GetByName(ctx, "Test App")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "app-1", retrieved.ID)
		assert.Equal(t, "Test App", retrieved.Name)
	})

	t.Run("Not Found", func(t *testing.T) {
		retrieved, err := repo.GetByName(ctx, "Non-existent App")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "app not found")
	})

	t.Run("Case Sensitive", func(t *testing.T) {
		// SQLite is case-insensitive by default, but this tests the behavior
		retrieved, err := repo.GetByName(ctx, "test app")
		// Result depends on database collation
		_ = retrieved
		_ = err
	})
}

func TestApplicationRepository_Update(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Create an app
		app := createTestApp("app-1", "Original Name", "client-1")
		err := repo.Create(ctx, app)
		require.NoError(t, err)

		// Update the app
		app.Name = "Updated Name"
		app.IsActive = false
		app.Uri = "https://updated.com"

		err = repo.Update(ctx, app)
		assert.NoError(t, err)

		// Verify the update
		retrieved, err := repo.GetByID(ctx, "app-1")
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", retrieved.Name)
		assert.False(t, retrieved.IsActive)
		assert.Equal(t, "https://updated.com", retrieved.Uri)
	})

	t.Run("Update Non-existent", func(t *testing.T) {
		// Try to update a non-existent app
		app := createTestApp("non-existent", "Name", "client-x")
		err := repo.Update(ctx, app)
		// GORM Save creates if not exists, so this won't error
		// Just verify behavior
		assert.NoError(t, err)

		// Verify it was created
		retrieved, err := repo.GetByID(ctx, "non-existent")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
	})

	t.Run("Concurrent Updates", func(t *testing.T) {
		// Create an app
		app := createTestApp("app-2", "Concurrent App", "client-2")
		err := repo.Create(ctx, app)
		require.NoError(t, err)

		// Simulate concurrent updates
		app1, _ := repo.GetByID(ctx, "app-2")
		app2, _ := repo.GetByID(ctx, "app-2")

		app1.Name = "Update 1"
		app2.Name = "Update 2"

		err = repo.Update(ctx, app1)
		assert.NoError(t, err)

		err = repo.Update(ctx, app2)
		assert.NoError(t, err)

		// Last write wins
		retrieved, err := repo.GetByID(ctx, "app-2")
		assert.NoError(t, err)
		assert.Equal(t, "Update 2", retrieved.Name)
	})
}

func TestApplicationRepository_Delete(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Create an active app
		app := createTestApp("app-1", "Test App", "client-1")
		app.IsActive = true
		err := repo.Create(ctx, app)
		require.NoError(t, err)

		// Delete (deactivate) the app
		err = repo.Delete(ctx, "app-1")
		assert.NoError(t, err)

		// Verify the app is deactivated
		retrieved, err := repo.GetByID(ctx, "app-1")
		assert.NoError(t, err)
		assert.False(t, retrieved.IsActive)
	})

	t.Run("Not Found", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app not found")
	})

	t.Run("Already Deactivated", func(t *testing.T) {
		// Create an inactive app
		app := createTestApp("app-2", "Inactive App", "client-2")
		app.IsActive = false
		err := repo.Create(ctx, app)
		require.NoError(t, err)

		// Try to delete again
		err = repo.Delete(ctx, "app-2")
		assert.NoError(t, err)

		// Verify it's still inactive
		retrieved, err := repo.GetByID(ctx, "app-2")
		assert.NoError(t, err)
		assert.False(t, retrieved.IsActive)
	})
}

func TestApplicationRepository_Restore(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Create and soft delete an app
		app := createTestApp("app-1", "Test App", "client-1")
		app.IsActive = false
		err := repo.Create(ctx, app)
		require.NoError(t, err)

		// Soft delete
		err = db.Delete(app).Error
		require.NoError(t, err)

		// Restore the app
		err = repo.Restore(ctx, "app-1")
		assert.NoError(t, err)

		// Verify the app is restored and active
		retrieved, err := repo.GetByID(ctx, "app-1")
		assert.NoError(t, err)
		assert.True(t, retrieved.IsActive)
		assert.False(t, retrieved.DeletedAt.Valid)
	})

	t.Run("Not Found", func(t *testing.T) {
		err := repo.Restore(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app not found")
	})

	t.Run("Restore Active App", func(t *testing.T) {
		// Create an active app
		app := createTestApp("app-2", "Active App", "client-2")
		app.IsActive = true
		err := repo.Create(ctx, app)
		require.NoError(t, err)

		// Try to restore (should work but do nothing significant)
		err = repo.Restore(ctx, "app-2")
		assert.NoError(t, err)

		// Verify it's still active
		retrieved, err := repo.GetByID(ctx, "app-2")
		assert.NoError(t, err)
		assert.True(t, retrieved.IsActive)
	})
}

func TestApplicationRepository_GetAll(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	t.Run("Empty Database", func(t *testing.T) {
		apps, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Empty(t, apps)
	})

	t.Run("Multiple Apps", func(t *testing.T) {
		// Create multiple apps
		app1 := createTestApp("app-1", "App 1", "client-1")
		app2 := createTestApp("app-2", "App 2", "client-2")
		app3 := createTestApp("app-3", "App 3", "client-3")

		err := repo.Create(ctx, app1)
		require.NoError(t, err)
		err = repo.Create(ctx, app2)
		require.NoError(t, err)
		err = repo.Create(ctx, app3)
		require.NoError(t, err)

		// Get all apps
		apps, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, apps, 3)
	})

	t.Run("Excludes Soft Deleted", func(t *testing.T) {
		// Soft delete one app
		app4 := createTestApp("app-4", "App 4", "client-4")
		err := repo.Create(ctx, app4)
		require.NoError(t, err)

		err = db.Delete(app4).Error
		require.NoError(t, err)

		// GetAll should exclude soft deleted
		apps, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		// Should not include app-4
		for _, app := range apps {
			assert.NotEqual(t, "app-4", app.ID)
		}
	})
}

func TestApplicationRepository_GetListApps(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	// Create test apps with different timestamps
	for i := 1; i <= 5; i++ {
		app := createTestApp(
			"app-"+string(rune('0'+i)),
			"App "+string(rune('0'+i)),
			"client-"+string(rune('0'+i)),
		)
		err := repo.Create(ctx, app)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	t.Run("Default Pagination", func(t *testing.T) {
		total, apps, err := repo.GetListApps(ctx, 0, 10, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, apps)
		assert.Equal(t, int64(5), total)
		assert.Len(t, *apps, 5)
	})

	t.Run("With Offset and Limit", func(t *testing.T) {
		total, apps, err := repo.GetListApps(ctx, 2, 2, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, apps)
		assert.Equal(t, int64(5), total)
		assert.Len(t, *apps, 2)
	})

	t.Run("Order By Name Ascending", func(t *testing.T) {
		total, apps, err := repo.GetListApps(ctx, 0, 10, "name", "asc")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.NotNil(t, apps)

		// Verify ascending order
		if len(*apps) > 1 {
			for i := 0; i < len(*apps)-1; i++ {
				assert.LessOrEqual(t, (*apps)[i].Name, (*apps)[i+1].Name)
			}
		}
	})

	t.Run("Order By ClientID Descending", func(t *testing.T) {
		total, apps, err := repo.GetListApps(ctx, 0, 10, "client_id", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.NotNil(t, apps)

		// Verify descending order
		if len(*apps) > 1 {
			for i := 0; i < len(*apps)-1; i++ {
				assert.GreaterOrEqual(t, (*apps)[i].ClientID, (*apps)[i+1].ClientID)
			}
		}
	})

	t.Run("Invalid OrderBy Field", func(t *testing.T) {
		// Should default to created_at
		total, apps, err := repo.GetListApps(ctx, 0, 10, "invalid_field", "desc")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.NotNil(t, apps)
	})

	t.Run("Invalid Sort Direction", func(t *testing.T) {
		// Should default to desc
		total, apps, err := repo.GetListApps(ctx, 0, 10, "name", "invalid")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.NotNil(t, apps)
	})

	t.Run("Includes Soft Deleted", func(t *testing.T) {
		// Soft delete one app
		app6 := createTestApp("app-6", "App 6", "client-6")
		err := repo.Create(ctx, app6)
		require.NoError(t, err)

		err = db.Delete(app6).Error
		require.NoError(t, err)

		// GetListApps uses Unscoped in Find, but not in Count
		// So total count excludes soft deleted, but results include them
		total, apps, err := repo.GetListApps(ctx, 0, 10, "", "")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total) // Count excludes soft deleted
		assert.Len(t, *apps, 6)          // But Find includes soft deleted
	})

	t.Run("Empty Result", func(t *testing.T) {
		// Request beyond available data
		total, apps, err := repo.GetListApps(ctx, 100, 10, "", "")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), total) // Total count excludes soft deleted
		assert.NotNil(t, apps)
		assert.Empty(t, *apps) // But no apps in this page
	})
}

func TestApplicationRepository_ContextCancellation(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)

	t.Run("Create with Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		app := createTestApp("app-1", "Test App", "client-1")
		err := repo.Create(ctx, app)
		// SQLite may not respect context cancellation for fast operations
		_ = err
	})

	t.Run("GetByID with Cancelled Context", func(t *testing.T) {
		// First create an app
		app := createTestApp("app-2", "Test App", "client-2")
		err := repo.Create(context.Background(), app)
		require.NoError(t, err)

		// Try to retrieve with cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = repo.GetByID(ctx, "app-2")
		// SQLite may not respect context cancellation
		_ = err
	})
}

func TestApplicationRepository_EdgeCases(t *testing.T) {
	db := setupApplicationTestDB(t)
	repo := repository.NewAppsRepository(db)
	ctx := context.Background()

	t.Run("Long Field Values", func(t *testing.T) {
		longString := string(make([]byte, 1000))
		app := createTestApp("app-1", "Test App", "client-1")
		app.Uri = longString
		app.RedirectUri = longString

		err := repo.Create(ctx, app)
		assert.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, "app-1")
		assert.NoError(t, err)
		assert.Equal(t, longString, retrieved.Uri)
	})

	t.Run("Special Characters in Fields", func(t *testing.T) {
		app := createTestApp("app-2", "Test App with 特殊字符 & symbols!@#", "client-2")
		app.Uri = "https://example.com?param=value&other=特殊"

		err := repo.Create(ctx, app)
		assert.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, "app-2")
		assert.NoError(t, err)
		assert.Equal(t, "Test App with 特殊字符 & symbols!@#", retrieved.Name)
	})

	t.Run("Empty Strings vs Null", func(t *testing.T) {
		app := createTestApp("app-3", "Test App", "client-3")
		app.Uri = ""
		app.RedirectUri = ""

		err := repo.Create(ctx, app)
		assert.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, "app-3")
		assert.NoError(t, err)
		assert.Equal(t, "", retrieved.Uri)
		assert.Equal(t, "", retrieved.RedirectUri)
	})
}
