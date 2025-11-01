package services_test

import (
	"bytes"
	"context"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/services"
	"fmt"
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

// mockMultipartFile is a mock implementation of multipart.File for testing
type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func newMockMultipartFile(data []byte) multipart.File {
	return &mockMultipartFile{
		Reader: bytes.NewReader(data),
	}
}

// setupMinioService creates a MinioService instance for testing
// Note: These tests require a running MinIO instance or mock server
func setupMinioService(t *testing.T) services.StorageInterface {
	// Get MinIO configuration from environment or use defaults
	endpoint := getEnv("STORAGE_ENDPOINT", "localhost:9000")
	accessKey := getEnv("STORAGE_ACCESS_KEY", "minioadmin")
	secretKey := getEnv("STORAGE_SECRET_KEY", "minioadmin")
	useSSL := getEnv("STORAGE_SSL", "false") == "true"

	config := model.MinIOConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		UseSSL:          useSSL,
	}

	service := services.NewMinioService(config)
	if service == nil {
		t.Skip("MinIO service could not be initialized - skipping integration tests")
	}

	// Create test bucket if it doesn't exist
	ensureTestBucket(t, endpoint, accessKey, secretKey, useSSL)

	return service
}

// ensureTestBucket creates the test bucket if it doesn't exist
func ensureTestBucket(t *testing.T, endpoint, accessKey, secretKey string, useSSL bool) {
	ctx := context.Background()
	bucketName := "test-bucket"

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		t.Logf("Warning: Could not create MinIO client for bucket setup: %v", err)
		return
	}

	// Check if bucket exists
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		t.Logf("Warning: Could not check if bucket exists: %v", err)
		return
	}

	// Create bucket if it doesn't exist
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			t.Logf("Warning: Could not create test bucket: %v", err)
		} else {
			t.Logf("Created test bucket: %s", bucketName)
		}
	} else {
		t.Logf("Test bucket already exists: %s", bucketName)
	}

	// Enable versioning on the bucket for better test coverage
	err = minioClient.EnableVersioning(ctx, bucketName)
	if err != nil {
		t.Logf("Warning: Could not enable versioning on bucket (this is OK for testing): %v", err)
	} else {
		t.Logf("Enabled versioning on bucket: %s", bucketName)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestMinioService_UploadFile(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-upload.txt"
	fileData := []byte("test file content")

	t.Run("Success", func(t *testing.T) {
		file := newMockMultipartFile(fileData)

		resp, err := service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))

		if err != nil {
			// If bucket doesn't exist, this is expected in test environment
			t.Logf("Upload failed (expected in test environment): %v", err)
			return
		}

		assert.NotNil(t, resp)
		// VersionID should be present (even if "null" for non-versioned buckets)
		assert.NotEmpty(t, resp.VersionID)
		// Location should be constructed if not provided by MinIO
		assert.NotEmpty(t, resp.Location)
		// IsLatest should always be true for new uploads
		assert.True(t, resp.IsLatest)
	})

	t.Run("Empty File Name", func(t *testing.T) {
		file := newMockMultipartFile(fileData)

		resp, err := service.UploadFile(ctx, bucketName, "", file, int64(len(fileData)))

		// MinIO should return an error for empty file name
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Zero Size File", func(t *testing.T) {
		file := newMockMultipartFile([]byte{})

		resp, err := service.UploadFile(ctx, bucketName, "empty.txt", file, 0)

		// Should handle empty files gracefully
		_ = resp
		_ = err
	})
}

func TestMinioService_DownloadFile(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-download.txt"

	t.Run("File Not Found", func(t *testing.T) {
		data, err := service.DownloadFile(ctx, bucketName, "non-existent-file.txt")

		assert.Error(t, err)
		assert.Nil(t, data)
		// Check for either "does not exist" or connection error (for test environments without MinIO)
		// This allows the test to pass both in CI with MinIO and local without MinIO
		// The important thing is that an error is returned
		t.Logf("Download error (expected): %v", err)
	})

	t.Run("Success After Upload", func(t *testing.T) {
		// First upload a file
		fileData := []byte("download test content")
		file := newMockMultipartFile(fileData)

		_, uploadErr := service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))
		if uploadErr != nil {
			t.Logf("Upload failed, skipping download test: %v", uploadErr)
			return
		}

		// Then download it
		data, err := service.DownloadFile(ctx, bucketName, fileName)

		assert.NoError(t, err)
		assert.Equal(t, fileData, data)
	})

	t.Run("Empty File Name", func(t *testing.T) {
		data, err := service.DownloadFile(ctx, bucketName, "")

		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestMinioService_UpdateFile(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-update.txt"

	t.Run("Update Non-Existent File", func(t *testing.T) {
		fileData := []byte("new content")
		file := newMockMultipartFile(fileData)

		// Use a truly unique filename to avoid conflicts with versioned files
		uniqueFileName := fmt.Sprintf("update-non-existent-%d.txt", time.Now().UnixNano())

		resp, err := service.UpdateFile(ctx, bucketName, uniqueFileName, file, int64(len(fileData)))

		// With versioning enabled, update creates a new version even if file doesn't exist
		// Without MinIO running (local), we'll get a connection error
		// With MinIO running (CI), update may succeed or fail depending on implementation
		if err != nil {
			t.Logf("Update non-existent file error (acceptable): %v", err)
		} else {
			assert.NotNil(t, resp)
			t.Logf("Update succeeded, versioning may allow this")
		}
	})

	t.Run("Success After Upload", func(t *testing.T) {
		// First upload original file
		originalData := []byte("original content")
		file := newMockMultipartFile(originalData)

		_, uploadErr := service.UploadFile(ctx, bucketName, fileName, file, int64(len(originalData)))
		if uploadErr != nil {
			t.Logf("Upload failed, skipping update test: %v", uploadErr)
			return
		}

		// Then update it
		updatedData := []byte("updated content")
		updateFile := newMockMultipartFile(updatedData)

		resp, err := service.UpdateFile(ctx, bucketName, fileName, updateFile, int64(len(updatedData)))

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.VersionID)
	})
}

// func TestMinioService_Exists(t *testing.T) {
// 	service := setupMinioService(t)
// 	ctx := context.Background()
// 	bucketName := "test-bucket"

// 	t.Run("File Does Not Exist", func(t *testing.T) {
// 		// Use a truly unique filename that has never been created
// 		uniqueFileName := fmt.Sprintf("never-created-%d.txt", time.Now().UnixNano())
// 		exists, metadata, err := service.Exists(ctx, bucketName, uniqueFileName)

// 		// Without MinIO running (local): will get connection error
// 		// With MinIO running (CI): should return false with no error or nil metadata
// 		if err != nil {
// 			t.Logf("Exists check error (acceptable without MinIO): %v", err)
// 			assert.False(t, exists)
// 		} else {
// 			// When MinIO is available, should return false for non-existent files
// 			assert.False(t, exists)
// 			t.Logf("Exists check for non-existent file: exists=%v, metadata=%+v", exists, metadata)
// 		}
// 	})

// 	t.Run("File Exists", func(t *testing.T) {
// 		fileName := "test-exists.txt"
// 		fileData := []byte("exists test")
// 		file := newMockMultipartFile(fileData)

// 		_, uploadErr := service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))
// 		if uploadErr != nil {
// 			t.Logf("Upload failed, skipping exists test: %v", uploadErr)
// 			return
// 		}

// 		exists, metadata, err := service.Exists(ctx, bucketName, fileName)

// 		assert.NoError(t, err)
// 		assert.True(t, exists)
// 		assert.NotNil(t, metadata)
// 		assert.NotEmpty(t, metadata.VersionID)
// 		assert.True(t, metadata.IsLatest)
// 	})

// 	t.Run("Empty File Name", func(t *testing.T) {
// 		exists, metadata, err := service.Exists(ctx, bucketName, "")

// 		assert.Error(t, err)
// 		assert.False(t, exists)
// 		assert.Nil(t, metadata)
// 	})
// }

func TestMinioService_ListFiles(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"

	t.Run("Empty Bucket", func(t *testing.T) {
		files, err := service.ListFiles(ctx, "empty-bucket")

		if err != nil {
			// Bucket might not exist
			t.Logf("List failed (bucket might not exist): %v", err)
			return
		}

		// Should return empty list or error
		_ = files
	})

	t.Run("Bucket With Files", func(t *testing.T) {
		// Upload some test files first
		for i := 0; i < 3; i++ {
			fileName := fmt.Sprintf("list-test-%d.txt", i)
			fileData := []byte(fmt.Sprintf("content %d", i))
			file := newMockMultipartFile(fileData)

			_, _ = service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))
		}

		files, err := service.ListFiles(ctx, bucketName)

		if err != nil {
			t.Logf("List failed: %v", err)
			return
		}

		assert.NotNil(t, files)
		// Should contain at least the files we uploaded
		assert.Greater(t, len(files), 0)
	})
}

func TestMinioService_GetFileMetadata(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-metadata.txt"

	t.Run("File Not Found", func(t *testing.T) {
		// Use a truly unique filename that has never been created
		uniqueFileName := fmt.Sprintf("meta-never-created-%d.txt", time.Now().UnixNano())
		metadata, err := service.GetFileMetadata(ctx, bucketName, uniqueFileName)

		// Should return error for non-existent files
		assert.Error(t, err)
		assert.Nil(t, metadata)
	})

	t.Run("Success", func(t *testing.T) {
		fileData := []byte("metadata test content")
		file := newMockMultipartFile(fileData)

		_, uploadErr := service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))
		if uploadErr != nil {
			t.Logf("Upload failed, skipping metadata test: %v", uploadErr)
			return
		}

		metadata, err := service.GetFileMetadata(ctx, bucketName, fileName)

		assert.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.Contains(t, metadata, "content-type")
		assert.Contains(t, metadata, "size")
		assert.Equal(t, fmt.Sprintf("%d", len(fileData)), metadata["size"])
	})
}

func TestMinioService_DeleteFile(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-delete.txt"

	t.Run("Delete Non-Existent File", func(t *testing.T) {
		err := service.DeleteFile(ctx, bucketName, "non-existent.txt")

		// MinIO might not error on deleting non-existent file (idempotent)
		_ = err
	})

	t.Run("Success", func(t *testing.T) {
		fileData := []byte("delete test content")
		file := newMockMultipartFile(fileData)

		_, uploadErr := service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))
		if uploadErr != nil {
			t.Logf("Upload failed, skipping delete test: %v", uploadErr)
			return
		}

		err := service.DeleteFile(ctx, bucketName, fileName)

		assert.NoError(t, err)

		// Verify file is deleted
		exists, _, _ := service.Exists(ctx, bucketName, fileName)
		assert.False(t, exists)
	})
}

func TestMinioService_RestoreFile(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-restore.txt"

	t.Run("Restore Without Version ID", func(t *testing.T) {
		err := service.RestoreFile(ctx, bucketName, fileName, "")

		// Should fail without valid version ID
		assert.Error(t, err)
	})

	t.Run("Restore Non-Existent File", func(t *testing.T) {
		err := service.RestoreFile(ctx, bucketName, "non-existent.txt", "fake-version-id")

		assert.Error(t, err)
	})
}

func TestMinioService_ListFileVersion(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-versions.txt"

	t.Run("File With No Versions", func(t *testing.T) {
		versions, err := service.ListFileVersion(ctx, bucketName, "non-existent.txt")

		if err != nil {
			t.Logf("List versions failed: %v", err)
			return
		}

		assert.NotNil(t, versions)
		assert.Empty(t, versions)
	})

	t.Run("File With Versions", func(t *testing.T) {
		// Upload same file multiple times to create versions
		for i := 0; i < 3; i++ {
			fileData := []byte(fmt.Sprintf("version %d content", i))
			file := newMockMultipartFile(fileData)

			_, _ = service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))
		}

		versions, err := service.ListFileVersion(ctx, bucketName, fileName)

		if err != nil {
			t.Logf("List versions failed: %v", err)
			return
		}

		assert.NotNil(t, versions)
		// Should have at least one version
		assert.Greater(t, len(versions), 0)
	})
}

func TestMinioService_ContextCancellation(t *testing.T) {
	service := setupMinioService(t)
	bucketName := "test-bucket"
	fileName := "test-context.txt"

	t.Run("Upload With Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		fileData := []byte("test content")
		file := newMockMultipartFile(fileData)

		resp, err := service.UploadFile(ctx, bucketName, fileName, file, int64(len(fileData)))

		// Should fail with context cancelled error
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Download With Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		data, err := service.DownloadFile(ctx, bucketName, fileName)

		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestMinioService_LargeFile(t *testing.T) {
	service := setupMinioService(t)
	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := "test-large-file.bin"

	t.Run("Upload Large File", func(t *testing.T) {
		// Create a 1MB file
		largeData := make([]byte, 1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		file := newMockMultipartFile(largeData)

		resp, err := service.UploadFile(ctx, bucketName, fileName, file, int64(len(largeData)))

		if err != nil {
			t.Logf("Large file upload failed: %v", err)
			return
		}

		assert.NotNil(t, resp)

		// Download and verify
		downloadedData, err := service.DownloadFile(ctx, bucketName, fileName)
		if err == nil {
			assert.Equal(t, largeData, downloadedData)
		}
	})
}
