package repository

import (
	"context"
	"crypsis-backend/internal/entity"

	"gorm.io/gorm"
)

// FileRepository defines the contract for file data access operations.
// It provides methods for creating, retrieving, updating, deleting, restoring files, handling metadata, and batch operations.
type FileRepository interface {
	// Create adds a new file record.
	Create(ctx context.Context, file *entity.Files) error
	// GetByID retrieves a file by its ID.
	GetByID(ctx context.Context, id string) (*entity.Files, error)
	// GetByHash retrieves a file by its hash value.
	GetByHash(ctx context.Context, hash string) (*entity.Files, error)
	// GetListFiles returns a paginated list of files for an app.
	GetListFiles(ctx context.Context, appID string, offset, limit int, orderBy, sort string) (int64, []entity.Files, error)
	// GetListFilesForAdmin returns a paginated list of files for admin users.
	GetListFilesForAdmin(ctx context.Context, offset, limit int, orderBy, sort string) (int64, []entity.Files, error)
	// GetAll retrieves all files.
	GetAll(ctx context.Context) ([]entity.Files, error)
	// Update modifies an existing file record.
	Update(ctx context.Context, file *entity.Files) error
	// Delete removes a file by its ID.
	Delete(ctx context.Context, id string) error
	// RestoreFile restores a deleted file by its ID.
	RestoreFile(ctx context.Context, fileID string) error
	// WithTransaction executes a function within a database transaction.
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	// CreateFileWithMetadata adds a new file and its metadata.
	CreateFileWithMetadata(ctx context.Context, file *entity.Files, metadata *entity.Metadata) error
	// GetMetadataByFileID retrieves metadata for a file by its ID.
	GetMetadataByFileID(ctx context.Context, fileID string) (*entity.Metadata, error)
	// GetMetadataByAppIDAndFileID retrieves metadata by app and file ID.
	GetMetadataByAppIDAndFileID(ctx context.Context, appID, fileID string) (*entity.Metadata, error)
	// GetDeletedMetadataByAppIDAndFileID retrieves deleted metadata by app and file ID.
	GetDeletedMetadataByAppIDAndFileID(ctx context.Context, appID, fileID string) (*entity.Metadata, error)
	// GetMetadataByEncHash retrieves metadata by encrypted hash.
	GetMetadataByEncHash(ctx context.Context, encHash string) (*entity.Metadata, error)
	// GetAllMetadata retrieves all metadata records.
	GetAllMetadata(ctx context.Context) ([]entity.Metadata, error)
	// GetAllKeyUIDs retrieves all key UIDs.
	GetAllKeyUIDs(ctx context.Context) ([]string, error)
	// UpdateFileAndMetadata updates both file and metadata records.
	UpdateFileAndMetadata(ctx context.Context, file *entity.Files, metadata *entity.Metadata) error
	// UpdateEncKeyByKeyUID updates the encryption key for a given key UID.
	UpdateEncKeyByKeyUID(ctx context.Context, keyUID string, newEncKey string) error
	// BatchUpdateEncKeys updates multiple encryption keys in batch.
	BatchUpdateEncKeys(ctx context.Context, updates map[string]string) error
}

// AdminRepository defines the contract for admin data access operations.
// It provides methods for creating, retrieving, updating, deleting, and verifying admin users.
type AdminRepository interface {
	// GetByUsername retrieves an admin by username.
	GetByUsername(ctx context.Context, username string) (*entity.Admins, error)
	// GetByID retrieves an admin by ID.
	GetByID(ctx context.Context, id string) (*entity.Admins, error)
	// GetByClientID retrieves an admin by client ID.
	GetByClientID(ctx context.Context, id string) (*entity.Admins, error)
	// GetList returns a paginated list of admins.
	GetList(ctx context.Context, offset, limit int, orderBy, sort string) (*[]entity.Admins, error)
	// Create adds a new admin record.
	Create(ctx context.Context, admin *entity.Admins) error
	// Update modifies an existing admin record.
	Update(ctx context.Context, admin *entity.Admins) error
	// Delete removes an admin by ID.
	Delete(ctx context.Context, id string) error
	// LoadAdminIDs loads all admin IDs into memory or cache.
	LoadAdminIDs(ctx context.Context) error
	// IsAdmin checks if the given ID belongs to an admin user.
	IsAdmin(ctx context.Context, id string) bool
}

// AppsRepository defines the contract for application data access operations.
// It provides methods for creating, retrieving, updating, deleting, restoring, and listing applications.
type AppsRepository interface {
	// Create adds a new application record.
	Create(ctx context.Context, app *entity.Apps) error
	// GetByID retrieves an application by its ID.
	GetByID(ctx context.Context, id string) (*entity.Apps, error)
	// GetByClientID retrieves an application by client ID.
	GetByClientID(ctx context.Context, clientID string) (*entity.Apps, error)
	// GetByClientIDLimited retrieves a limited application record by client ID.
	GetByClientIDLimited(ctx context.Context, clientID string) (*entity.Apps, error)
	// GetByName retrieves an application by its name.
	GetByName(ctx context.Context, name string) (*entity.Apps, error)
	// Update modifies an existing application record.
	Update(ctx context.Context, app *entity.Apps) error
	// Delete removes an application by its ID.
	Delete(ctx context.Context, id string) error
	// Restore restores a deleted application by its ID.
	Restore(ctx context.Context, id string) error
	// GetListApps returns a paginated list of applications.
	GetListApps(ctx context.Context, offset, limit int, orderBy, sort string) (int64, *[]entity.Apps, error)
	// GetAll retrieves all applications.
	GetAll(ctx context.Context) ([]entity.Apps, error)
}

// FileLogsRepository defines the contract for file log data access operations.
// It provides methods for creating, listing, retrieving, and deleting file logs.
type FileLogsRepository interface {
	// Create adds a new file log record.
	Create(ctx context.Context, log *entity.FileLogs) error
	// List returns a paginated list of file logs.
	List(ctx context.Context, offset, limit int, orderBy, sort string) (int64, *[]entity.FileLogs, error)
	// GetByFileID retrieves file logs by file ID.
	GetByFileID(ctx context.Context, fileID string) (*[]entity.FileLogs, error)
	// GetByAction retrieves file logs by action type.
	GetByAction(ctx context.Context, action string) (*[]entity.FileLogs, error)
	// DeleteOldLogs deletes file logs older than the specified number of days.
	DeleteOldLogs(ctx context.Context, days int) error
}
