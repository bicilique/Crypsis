package services

import (
	"bytes"
	"context"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioService implements the StorageInterface for MinIO object storage.
// It provides methods for file upload, download, update, deletion, restoration,
// and metadata operations in a MinIO-compatible object storage system.
type MinioService struct {
	client *minio.Client
}

// NewMinioService creates a new MinIO service instance with the provided configuration.
// It initializes a MinIO client with the specified endpoint, credentials, and SSL settings.
// Returns nil if client initialization fails.
func NewMinioService(input model.MinIOConfig) StorageInterface {
	client, err := minio.New(input.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(input.AccessKeyID, input.SecretAccessKey, ""),
		Secure: input.UseSSL,
	})
	if err != nil {
		log.Fatalf("Error initialize minio client: %v", err)
		return nil
	}

	return &MinioService{
		client: client,
	}
}

// UploadFile uploads a file to the specified bucket with the given name and size.
// Returns a StorageTransactionResponse containing transaction metadata including version ID, last modified time,
// expiration, location, checksum, and latest status.
func (s *MinioService) UploadFile(ctx context.Context, bucketName, fileName string, file multipart.File, fileSize int64) (*model.StorageTransactionResponse, error) {
	// Start tracing span
	tracer := helper.GetTracingHelper()
	ctx, span := tracer.StartStorageSpan(ctx, "PutObject", bucketName, fileName)
	defer span.End()

	helper.AddAttributes(span, map[string]interface{}{
		"file.size": fileSize,
	})

	resp, err := s.client.PutObject(ctx, bucketName, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		slog.Error("failed to upload file",
			"error", err,
			"bucket", bucketName,
			"file", fileName,
		)
		helper.RecordError(span, err)
		return nil, fmt.Errorf("upload failed for file %s: %w", fileName, err)
	}

	// Construct location if not provided by MinIO
	location := resp.Location
	if location == "" {
		// Get endpoint from client
		endpoint := s.client.EndpointURL().String()
		location = fmt.Sprintf("%s/%s/%s", endpoint, bucketName, fileName)
	}

	// Version ID may be empty if versioning is not enabled on the bucket
	versionID := resp.VersionID
	if versionID == "" {
		// For non-versioned buckets, we can use "null" as a placeholder
		// or leave it empty - tests should handle this gracefully
		versionID = "null"
	}

	return &model.StorageTransactionResponse{
		VersionID:      versionID,
		LastModified:   resp.LastModified.String(),
		Expiration:     resp.Expiration.String(),
		Location:       location,
		ChecksumSHA256: resp.ChecksumSHA256,
		IsLatest:       true,
		IsDeleteMarker: false,
	}, nil
}

// DownloadFile downloads a file from the specified bucket and returns its contents as bytes.
// Returns an error if the file does not exist or cannot be read.
func (s *MinioService) DownloadFile(ctx context.Context, bucketName, fileName string) ([]byte, error) {
	isExist, _, err := s.Exists(ctx, bucketName, fileName)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf("file %s does not exist", fileName)
	}

	object, err := s.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		slog.Error("failed to download file",
			"error", err,
			"bucket", bucketName,
			"file", fileName,
		)
		return nil, fmt.Errorf("download failed for file %s: %w", fileName, err)
	}
	defer object.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object)
	if err != nil {
		slog.Error("failed to read file", slog.Any("error", err))
		return nil, fmt.Errorf("error reading file %s: %w", fileName, err)
	}

	return buf.Bytes(), nil
}

// UpdateFile updates an existing file in the specified bucket by uploading a new version.
// This is effectively an alias for UploadFile since MinIO handles versioning automatically.
func (s *MinioService) UpdateFile(ctx context.Context, bucketName, fileName string, file multipart.File, fileSize int64) (*model.StorageTransactionResponse, error) {
	_, _, err := s.Exists(ctx, bucketName, fileName)
	if err != nil {
		return nil, err
	}
	return s.UploadFile(ctx, bucketName, fileName, file, fileSize)
}

// Exists checks if a file exists in the specified bucket and returns its metadata.
// Returns true if the file exists (even if marked for deletion), along with transaction metadata.
func (s *MinioService) Exists(ctx context.Context, bucketName, fileName string) (bool, *model.StorageTransactionResponse, error) {
	// Validate input
	if fileName == "" {
		return false, nil, fmt.Errorf("object name cannot be empty")
	}

	resp, err := s.client.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		// Convert MinIO "not found" errors to a friendly (false, nil) result
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "the specified key does not exist") ||
			strings.Contains(errMsg, "no such key") ||
			strings.Contains(errMsg, "not found") ||
			strings.Contains(errMsg, "nosuchkey") {
			return false, nil, nil
		}

		slog.Error("failed to check file existence",
			"error", err,
			"bucket", bucketName,
			"file", fileName,
		)
		return false, nil, fmt.Errorf("error checking file %s existence: %w", fileName, err)
	}

	// Check if file is marked for deletion
	if resp.IsDeleteMarker {
		result := &model.StorageTransactionResponse{
			VersionID:      resp.VersionID,
			IsDeleteMarker: true,
		}
		slog.Warn("file is marked as deleted", slog.Any("result", result))
		return true, result, nil
	}

	return true, &model.StorageTransactionResponse{
		VersionID:      resp.VersionID,
		LastModified:   resp.LastModified.String(),
		Expiration:     resp.Expiration.String(),
		ChecksumSHA256: resp.ChecksumSHA256,
		IsLatest:       resp.IsLatest,
		IsDeleteMarker: resp.IsDeleteMarker,
	}, nil
}

// ListFiles lists all files in the specified bucket recursively.
// Returns a slice of file keys (object names).
func (s *MinioService) ListFiles(ctx context.Context, bucketName string) ([]string, error) {
	files := []string{}
	objectCh := s.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			slog.Error("failed to list files",
				"error", object.Err,
				"bucket", bucketName,
			)
			return nil, fmt.Errorf("error listing files: %w", object.Err)
		}
		files = append(files, object.Key)
	}
	return files, nil
}

// GetFileMetadata retrieves metadata of a specific file in the specified bucket.
// Returns a map containing content-type and size information.
func (s *MinioService) GetFileMetadata(ctx context.Context, bucketName, fileName string) (map[string]string, error) {
	info, err := s.client.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		slog.Error("failed to get file metadata",
			"error", err,
			"bucket", bucketName,
			"file", fileName,
		)
		return nil, fmt.Errorf("error getting metadata for file %s: %w", fileName, err)
	}
	return map[string]string{
		"content-type": info.ContentType,
		"size":         fmt.Sprintf("%d", info.Size),
	}, nil
}

// DeleteFile performs a soft delete of a file in the specified bucket.
// The file is marked for deletion but can be restored using its version ID.
func (s *MinioService) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	err := s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{VersionID: ""})
	if err != nil {
		slog.Error("failed to soft delete object",
			"error", err,
			"bucket", bucketName,
			"object", objectName,
		)
		return fmt.Errorf("soft delete failed for object %s: %w", objectName, err)
	}
	return nil
}

// RestoreFile restores a soft-deleted file in the specified bucket using its version ID.
// The file is restored for 30 days with Standard tier access.
func (s *MinioService) RestoreFile(ctx context.Context, bucketName, fileName, versionID string) error {
	days := 30
	typeMinio := minio.RestoreSelect
	tierMinio := minio.TierStandard

	err := s.client.RestoreObject(ctx, bucketName, fileName, versionID, minio.RestoreRequest{
		Days: &days,
		Type: &typeMinio,
		Tier: &tierMinio,
	})
	if err != nil {
		slog.Error("failed to restore object",
			"error", err,
			"bucket", bucketName,
			"file", fileName,
			"version", versionID,
		)
		return fmt.Errorf("restore failed for object %s: %w", fileName, err)
	}
	return nil
}

// ListFileVersion lists all versions of a specific file in the specified bucket.
// Returns a slice of version keys for the given file prefix.
func (s *MinioService) ListFileVersion(ctx context.Context, bucketName, fileName string) ([]string, error) {
	versions := []string{}
	objectCh := s.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix: fileName,
	})
	for object := range objectCh {
		if object.Err != nil {
			slog.Error("failed to list object versions",
				"error", object.Err,
				"bucket", bucketName,
				"file", fileName,
			)
			return nil, fmt.Errorf("error listing object versions: %w", object.Err)
		}
		versions = append(versions, object.Key)
	}
	return versions, nil
}
