package services

import (
	"context"
	"crypsis-backend/internal/model"
	"mime/multipart"
)

// ApplicationInterface defines the contract for application management operations.
// It provides methods for adding, updating, retrieving, listing, deleting, recovering applications, and rotating secrets.
type ApplicationInterface interface {
	// AddApp creates a new application with the specified name, URI, and redirect URI.
	AddApp(ctx context.Context, appName, appUri, redirectUri string) (map[string]interface{}, error)
	// UpdateApp updates an existing application's details using its UID.
	UpdateApp(ctx context.Context, appUID string, appName, appUri, redirectUri string) (map[string]interface{}, error)
	// GetInfo retrieves information about an application by its UID.
	GetInfo(ctx context.Context, appUID string) (map[string]interface{}, error)
	// ListApps returns a paginated and sorted list of applications.
	ListApps(ctx context.Context, limit, offset int, sortBy, order string) (int64, map[string]interface{}, error)
	// DeleteApp removes an application by its UID.
	DeleteApp(ctx context.Context, appUID string) error
	// RecoverApp restores a deleted application by its UID.
	RecoverApp(ctx context.Context, appUID string) (*string, error)
	// RotateSecret rotates the secret for the specified application UID.
	RotateSecret(ctx context.Context, appUID string) (map[string]interface{}, error)
}

// AdminInterface defines the contract for administrative user management operations.
// It provides methods for authentication, registration, token management, updating credentials, deleting admins, and listing admin users.
type AdminInterface interface {
	// Login authenticates an admin user and returns an access token.
	Login(ctx context.Context, username string, password string) (string, error)
	// Register creates a new admin user and returns an access token.
	Register(ctx context.Context, username string, password string) (string, error)
	// RefreshToken generates a new access token for the admin user.
	RefreshToken(ctx context.Context, adminID, accessToken string) (string, error)
	// RevokeToken revokes the access token for the admin user.
	RevokeToken(ctx context.Context, adminID, accessToken string) (string, error)
	// UpdateUsername changes the username for the specified admin user.
	UpdateUsername(ctx context.Context, adminID, newUsername string) error
	// UpdatePassword changes the password for the specified admin user.
	UpdatePassword(ctx context.Context, adminID, newPassword string) error
	// DeleteAdmin removes the admin user identified by requestID.
	DeleteAdmin(ctx context.Context, requestID string) (string, error)
	// GetAdminList returns a paginated list of admin users.
	GetAdminList(ctx context.Context, offset int, limit int) (map[string]interface{}, error)
}

// FileInterface defines the contract for file management operations.
// It provides methods for uploading, downloading, encrypting, decrypting, updating, deleting, recovering files, managing file metadata, re-keying, and listing files and logs.
type FileInterface interface {
	// Uploads a file and returns a unique file UID that can be used to download the file
	UploadFile(ctx context.Context, clientID, fileName string, input multipart.File) (fileUID string, err error)
	// Downloads a file and returns its decrypted form and its name
	DownloadFile(ctx context.Context, clientID, fileUID string) ([]byte, string, error)
	// Encrypts a file and returns encrypted form and its name
	EncryptFile(ctx context.Context, clientID, filename string, input multipart.File) ([]byte, string, error)
	// Decrypts a file and returns decrypted form
	DecryptFile(ctx context.Context, clientID, fileUID string, input multipart.File) ([]byte, error)
	// Returns metadata of a file
	GetFileMetadata(ctx context.Context, clientID, fileUID string) (map[string]interface{}, error)
	// Updates a file in storage
	UpdateFile(ctx context.Context, appID, clientID, fileName string, input multipart.File) (string, error)
	// Deletes a file from storage
	DeleteFile(ctx context.Context, clientID, fileUID string) error
	// Recovers a file from storage
	RecoverFile(ctx context.Context, clientID, fileUID string) (string, error)
	// Generates a new key and re-encrypts all files with the new key
	ReKey(ctx context.Context, clientID, keyUID string) (string, error)
	// Return a list of files
	ListFiles(ctx context.Context, clientID string, limit, offset int, sortBy, order string) (int64, []map[string]interface{}, error)
	// Return a list of files for admin only
	ListFilesForAdmin(ctx context.Context, adminID, appID string, limit, offset int, sortBy, order string) (int64, []map[string]interface{}, error)
	// Return a list of logs for admin only
	ListLogs(ctx context.Context, limit, offset int, sortBy, order string) (int64, []map[string]interface{}, error)
}

// KMSInterface defines the contract for Key Management Service operations.
// It provides methods for key generation, encryption, decryption, and key lifecycle management.
type KMSInterface interface {
	// GenerateSymetricKey creates a new symmetric key with the given name.
	GenerateSymetricKey(ctx context.Context, name string) (string, error)
	// GenerateKeyPair creates a new asymmetric key pair with the given name.
	GenerateKeyPair(ctx context.Context, name string) (string, string, error)
	// ExportKey exports the key identified by keyUID.
	ExportKey(ctx context.Context, keyUID string) (string, error)
	// LocateKey finds keys by name and returns their identifiers.
	LocateKey(ctx context.Context, name string) ([]string, error)
	// Encrypt encrypts the given text using the specified key.
	Encrypt(ctx context.Context, keyUID string, text string) (string, string, string, error)
	// Decrypt decrypts the encrypted data using the specified key and parameters.
	Decrypt(ctx context.Context, keyUID, encryptedData, ivCounterNonce, authTag string) (string, error)
	// DestroyKey permanently deletes the key identified by keyUID.
	DestroyKey(ctx context.Context, keyUID string) (string, error)
	// RevokeKey revokes the key identified by keyUID.
	RevokeKey(ctx context.Context, keyUID string) (string, error)
	// ReKey rotates the key identified by keyUID.
	ReKey(ctx context.Context, keyUID string) (string, error)
	// Covercrypt performs covercrypt operation using the specified key and text.
	Covercrypt(ctx context.Context, keyUID string, text string) (string, error)
}

// OAuth2Interface defines the contract for OAuth2 client and token management.
// It provides generic methods for client CRUD operations and token handling.
type OAuth2Interface interface {
	// CreateClient creates a new OAuth2 client with the provided input.
	CreateClient(ctx context.Context, input *model.ApplicationRequest) (*model.OAuth2ClientResponse, error)
	// GetClient retrieves details of an OAuth2 client by clientId.
	GetClient(ctx context.Context, clientId string) (*model.OAuth2ClientResponse, error)
	// ListClients returns a paginated list of OAuth2 clients.
	ListClients(ctx context.Context, pageToken string, perPage int64) ([]*model.OAuth2ClientResponse, error)
	// UpdateClient updates an OAuth2 client with the specified parameters.
	UpdateClient(ctx context.Context, clientName string, operations string, path string, value interface{}) (*model.OAuth2ClientResponse, error)
	// DeleteClient removes an OAuth2 client by clientId.
	DeleteClient(ctx context.Context, clientId string) error
	// IntrospectToken checks the validity and details of a token for a given scope.
	IntrospectToken(ctx context.Context, token, scope string) (*model.OAuth2TokenResponse, error)
	// TokenRequest requests a new token using the provided input.
	TokenRequest(ctx context.Context, input *model.TokenRequest) (*model.TokenResponse, error)
	// RevokeToken revokes a refresh token for the specified client.
	RevokeToken(ctx context.Context, clientId, clientSecret, refreshToken string) (string, error)
}

// StorageInterface defines the contract for file storage operations.
// It provides methods for uploading, downloading, updating, deleting, and managing files and their metadata.
type StorageInterface interface {
	// UploadFile uploads a file to the specified bucket with the given name and size.
	UploadFile(ctx context.Context, bucketName string, fileName string, file multipart.File, fileSize int64) (*model.StorageTransactionResponse, error)
	// DownloadFile retrieves the contents of a file from the specified bucket.
	DownloadFile(ctx context.Context, bucketName string, fileName string) ([]byte, error)
	// DeleteFile removes a file from the specified bucket.
	DeleteFile(ctx context.Context, bucketName string, fileName string) error
	// UpdateFile replaces an existing file in the bucket with a new file and size.
	UpdateFile(ctx context.Context, bucketName, fileName string, file multipart.File, fileSize int64) (*model.StorageTransactionResponse, error)
	// Exists checks if a file exists in the specified bucket and returns its metadata if present.
	Exists(ctx context.Context, bucketName string, fileName string) (bool, *model.StorageTransactionResponse, error)
	// ListFiles returns a list of all file names in the specified bucket.
	ListFiles(ctx context.Context, bucketName string) ([]string, error)
	// GetFileMetadata retrieves metadata for a file in the specified bucket.
	GetFileMetadata(ctx context.Context, bucketName string, fileName string) (map[string]string, error)
	// RestoreFile restores a file to a previous version using the version ID.
	RestoreFile(ctx context.Context, bucketName, fileName, versionID string) error
	// ListFileVersion lists all version IDs for a file in the specified bucket.
	ListFileVersion(ctx context.Context, bucketName, fileName string) ([]string, error)
}

// CryptographicInterface defines the contract for cryptographic operations.
// It provides methods for key generation, encryption, decryption, hashing, and key derivation for both strings and files.
type CryptographicInterface interface {
	// GenerateKey creates a new cryptographic key.
	GenerateKey() (string, error)
	// EncryptString encrypts a string using the provided key.
	EncryptString(key, text string) (string, error)
	// DecryptString decrypts a string using the provided key.
	DecryptString(key, text string) (string, error)
	// HashString generates a hash of the given text using the specified hash method.
	HashString(hashMethod, text string) (string, error)
	// CompareHash compares a hash with the hash of the given text using the specified method.
	CompareHash(hashMethod, text, hash string) bool
	// EncryptFile encrypts a file (as bytes) using the provided key.
	EncryptFile(key string, file []byte) ([]byte, error)
	// DecryptFile decrypts a file (as bytes) using the provided key.
	DecryptFile(key string, file []byte) ([]byte, error)
	// HashFile generates a hash of the file using the specified hash method.
	HashFile(hashMethod string, file []byte) (string, error)
	// CompareHashFile compares a hash with the hash of the given file using the specified method.
	CompareHashFile(hashMethod string, file []byte, hash string) bool
	// KeyDerivationFunction derives a key from the input and salt.
	KeyDerivationFunction(input string, salt []byte) ([]byte, error)
}
