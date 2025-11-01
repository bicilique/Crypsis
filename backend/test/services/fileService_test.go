package services

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/model"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

// Mock implementations for dependencies
type MockCryptographicService struct {
	mock.Mock
}

func (m *MockCryptographicService) GenerateKey() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockCryptographicService) KeyDerivationFunction(input string, salt []byte) ([]byte, error) {
	args := m.Called(input, salt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptographicService) EncryptString(key, plainText string) (string, error) {
	args := m.Called(key, plainText)
	return args.String(0), args.Error(1)
}

func (m *MockCryptographicService) DecryptString(key, cipherText string) (string, error) {
	args := m.Called(key, cipherText)
	return args.String(0), args.Error(1)
}

func (m *MockCryptographicService) EncryptFile(key string, fileBytes []byte) ([]byte, error) {
	args := m.Called(key, fileBytes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptographicService) DecryptFile(key string, encryptedFileBytes []byte) ([]byte, error) {
	args := m.Called(key, encryptedFileBytes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptographicService) ImportRawKeyAsBase64(keyBytes []byte) (string, error) {
	args := m.Called(keyBytes)
	return args.String(0), args.Error(1)
}

func (m *MockCryptographicService) KeysetFromRawAES256GCM(rawKey []byte) (*keyset.Handle, error) {
	args := m.Called(rawKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*keyset.Handle), args.Error(1)
}

func (m *MockCryptographicService) HashString(hashMethod, text string) (string, error) {
	args := m.Called(hashMethod, text)
	return args.String(0), args.Error(1)
}

func (m *MockCryptographicService) CompareHash(hashMethod, text, hash string) bool {
	args := m.Called(hashMethod, text, hash)
	return args.Bool(0)
}

func (m *MockCryptographicService) HashFile(hashMethod string, file []byte) (string, error) {
	args := m.Called(hashMethod, file)
	return args.String(0), args.Error(1)
}

func (m *MockCryptographicService) CompareHashFile(hashMethod string, file []byte, hash string) bool {
	args := m.Called(hashMethod, file, hash)
	return args.Bool(0)
}

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) UploadFile(ctx context.Context, bucketName, fileName string, file multipart.File, size int64) (*model.StorageTransactionResponse, error) {
	args := m.Called(ctx, bucketName, fileName, file, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StorageTransactionResponse), args.Error(1)
}

func (m *MockStorageService) DownloadFile(ctx context.Context, bucketName, fileName string) ([]byte, error) {
	args := m.Called(ctx, bucketName, fileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageService) Exists(ctx context.Context, bucketName, fileName string) (bool, *model.StorageTransactionResponse, error) {
	args := m.Called(ctx, bucketName, fileName)
	if args.Get(1) == nil {
		return args.Bool(0), nil, args.Error(2)
	}
	return args.Bool(0), args.Get(1).(*model.StorageTransactionResponse), args.Error(2)
}

func (m *MockStorageService) DeleteFile(ctx context.Context, bucketName, fileName string) error {
	args := m.Called(ctx, bucketName, fileName)
	return args.Error(0)
}

func (m *MockStorageService) UpdateFile(ctx context.Context, bucketName, fileName string, file multipart.File, size int64) (*model.StorageTransactionResponse, error) {
	args := m.Called(ctx, bucketName, fileName, file, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StorageTransactionResponse), args.Error(1)
}

func (m *MockStorageService) ListFiles(ctx context.Context, bucketName string) ([]string, error) {
	args := m.Called(ctx, bucketName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorageService) GetFileMetadata(ctx context.Context, bucketName, fileName string) (map[string]string, error) {
	args := m.Called(ctx, bucketName, fileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockStorageService) RestoreFile(ctx context.Context, bucketName, fileName, versionID string) error {
	args := m.Called(ctx, bucketName, fileName, versionID)
	return args.Error(0)
}

func (m *MockStorageService) ListFileVersion(ctx context.Context, bucketName, fileName string) ([]string, error) {
	args := m.Called(ctx, bucketName, fileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type MockKMSService struct {
	mock.Mock
}

func (m *MockKMSService) GenerateSymetricKey(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockKMSService) GenerateKeyPair(ctx context.Context, name string) (string, string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockKMSService) ExportKey(ctx context.Context, keyUID string) (string, error) {
	args := m.Called(ctx, keyUID)
	return args.String(0), args.Error(1)
}

func (m *MockKMSService) LocateKey(ctx context.Context, name string) ([]string, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockKMSService) Encrypt(ctx context.Context, keyUID string, text string) (string, string, string, error) {
	args := m.Called(ctx, keyUID, text)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockKMSService) Decrypt(ctx context.Context, keyUID, encryptedData, ivCounterNonce, authTag string) (string, error) {
	args := m.Called(ctx, keyUID, encryptedData, ivCounterNonce, authTag)
	return args.String(0), args.Error(1)
}

func (m *MockKMSService) DestroyKey(ctx context.Context, keyUID string) (string, error) {
	args := m.Called(ctx, keyUID)
	return args.String(0), args.Error(1)
}

func (m *MockKMSService) RevokeKey(ctx context.Context, keyUID string) (string, error) {
	args := m.Called(ctx, keyUID)
	return args.String(0), args.Error(1)
}

func (m *MockKMSService) ReKey(ctx context.Context, keyUID string) (string, error) {
	args := m.Called(ctx, keyUID)
	return args.String(0), args.Error(1)
}

func (m *MockKMSService) Covercrypt(ctx context.Context, keyUID string, text string) (string, error) {
	args := m.Called(ctx, keyUID, text)
	return args.String(0), args.Error(1)
}

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) GetListFiles(ctx context.Context, appID string, offset, limit int, sortBy, order string) (int64, []entity.Files, error) {
	args := m.Called(ctx, appID, offset, limit, sortBy, order)
	if args.Get(1) == nil {
		return args.Get(0).(int64), nil, args.Error(2)
	}
	return args.Get(0).(int64), args.Get(1).([]entity.Files), args.Error(2)
}

func (m *MockFileRepository) GetListFilesForAdmin(ctx context.Context, offset, limit int, sortBy, order string) (int64, []entity.Files, error) {
	args := m.Called(ctx, offset, limit, sortBy, order)
	if args.Get(1) == nil {
		return args.Get(0).(int64), nil, args.Error(2)
	}
	return args.Get(0).(int64), args.Get(1).([]entity.Files), args.Error(2)
}

func (m *MockFileRepository) GetMetadataByAppIDAndFileID(ctx context.Context, appID, fileID string) (*entity.Metadata, error) {
	args := m.Called(ctx, appID, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Metadata), args.Error(1)
}

func (m *MockFileRepository) CreateFileWithMetadata(ctx context.Context, file *entity.Files, metadata *entity.Metadata) error {
	args := m.Called(ctx, file, metadata)
	return args.Error(0)
}

func (m *MockFileRepository) GetByID(ctx context.Context, fileID string) (*entity.Files, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Files), args.Error(1)
}

func (m *MockFileRepository) Delete(ctx context.Context, fileID string) error {
	args := m.Called(ctx, fileID)
	return args.Error(0)
}

func (m *MockFileRepository) Recover(ctx context.Context, fileID string) error {
	args := m.Called(ctx, fileID)
	return args.Error(0)
}

func (m *MockFileRepository) BatchUpdateEncKeys(ctx context.Context, updates map[string]string) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockFileRepository) Create(ctx context.Context, file *entity.Files) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) GetAll(ctx context.Context) ([]entity.Files, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Files), args.Error(1)
}

func (m *MockFileRepository) GetAllKeyUIDs(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockFileRepository) GetAllMetadata(ctx context.Context) ([]entity.Metadata, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Metadata), args.Error(1)
}

func (m *MockFileRepository) GetByHash(ctx context.Context, hash string) (*entity.Files, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Files), args.Error(1)
}

func (m *MockFileRepository) GetByName(ctx context.Context, appID, name string) (entity.Metadata, error) {
	args := m.Called(ctx, appID, name)
	return args.Get(0).(entity.Metadata), args.Error(1)
}

func (m *MockFileRepository) GetDeletedMetadataByAppIDAndFileID(ctx context.Context, appID, fileID string) (*entity.Metadata, error) {
	args := m.Called(ctx, appID, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Metadata), args.Error(1)
}

type MockFileLogsRepository struct {
	mock.Mock
}

func (m *MockFileLogsRepository) Create(ctx context.Context, log *entity.FileLogs) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockFileLogsRepository) GetListLogs(ctx context.Context, offset, limit int, sortBy, order string) (int64, *[]entity.FileLogs, error) {
	args := m.Called(ctx, offset, limit, sortBy, order)
	if args.Get(1) == nil {
		return args.Get(0).(int64), nil, args.Error(2)
	}
	return args.Get(0).(int64), args.Get(1).(*[]entity.FileLogs), args.Error(2)
}

func (m *MockFileLogsRepository) DeleteOldLogs(ctx context.Context, days int) error {
	args := m.Called(ctx, days)
	return args.Error(0)
}

func (m *MockFileLogsRepository) GetByAction(ctx context.Context, action string) (*[]entity.FileLogs, error) {
	args := m.Called(ctx, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]entity.FileLogs), args.Error(1)
}

func (m *MockFileLogsRepository) GetByFileID(ctx context.Context, fileID string) (*[]entity.FileLogs, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]entity.FileLogs), args.Error(1)
}

func (m *MockFileLogsRepository) List(ctx context.Context, offset, limit int, sortBy, order string) (int64, *[]entity.FileLogs, error) {
	args := m.Called(ctx, offset, limit, sortBy, order)
	if args.Get(1) == nil {
		return args.Get(0).(int64), nil, args.Error(2)
	}
	return args.Get(0).(int64), args.Get(1).(*[]entity.FileLogs), args.Error(2)
}

type MockApplicationRepository struct {
	mock.Mock
}

func (m *MockApplicationRepository) GetByClientID(ctx context.Context, clientID string) (*entity.Apps, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Apps), args.Error(1)
}

func (m *MockApplicationRepository) GetByID(ctx context.Context, appID string) (*entity.Apps, error) {
	args := m.Called(ctx, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Apps), args.Error(1)
}

func (m *MockApplicationRepository) Create(ctx context.Context, app *entity.Apps) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockApplicationRepository) Update(ctx context.Context, app *entity.Apps) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockApplicationRepository) Delete(ctx context.Context, appID string) error {
	args := m.Called(ctx, appID)
	return args.Error(0)
}

func (m *MockApplicationRepository) Recover(ctx context.Context, appID string) error {
	args := m.Called(ctx, appID)
	return args.Error(0)
}

func (m *MockApplicationRepository) GetList(ctx context.Context, offset, limit int, sortBy, order string) (int64, *[]entity.Apps, error) {
	args := m.Called(ctx, offset, limit, sortBy, order)
	if args.Get(1) == nil {
		return args.Get(0).(int64), nil, args.Error(2)
	}
	return args.Get(0).(int64), args.Get(1).(*[]entity.Apps), args.Error(2)
}

func (m *MockApplicationRepository) GetAll(ctx context.Context) ([]entity.Apps, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Apps), args.Error(1)
}

func (m *MockApplicationRepository) GetByClientIDLimited(ctx context.Context, clientID string) (*entity.Apps, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Apps), args.Error(1)
}

func (m *MockApplicationRepository) Restore(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCryptographicService mocks the CryptographicService interface
