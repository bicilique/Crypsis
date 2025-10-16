package model

import (
	"mime/multipart"
	"time"
)

// UploadFileRequest represents the request body for uploading a file.
type UploadFileRequest struct {
	FileName string                `form:"file_name" binding:"required"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
}

// EncryptFileRequest represents the request body for encrypting a file.
type EncryptFileRequest struct {
	FileName string                `form:"file_name" binding:"required"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
}

// DecryptFileRequest represents the request body for decrypting a file.
type DecryptFileRequest struct {
	FileName string                `form:"file_name" binding:"required"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
}

// UploadFileResponse represents the response body after a file upload.
type UploadFileResponse struct {
	FileName string `json:"file_name"`
	FileID   string `json:"file_id"`
}

// DownloadFileResponse represents the response body when a file is downloaded.
type DownloadFileResponse struct {
	FileName string `json:"file_name"`
	Content  []byte `json:"content"`
}

// EncryptFileResponse represents the response body after encrypting a file.
type EncryptFileResponse struct {
	FileName      string `json:"file_name"`
	EncryptedData []byte `json:"encrypted_data"`
}

// DecryptFileResponse represents the response body after decrypting a file.
type DecryptFileResponse struct {
	FileName      string `json:"file_name"`
	DecryptedData []byte `json:"decrypted_data"`
}

type ListFilesResponse struct {
	Files []FileResponse `json:"files"`
	Count int64          `json:"count"`
}

type FileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"file_name"`
	Size      int64  `json:"file_size"`
	OwnerID   string `json:"app_id,omitempty"`
	MimeType  string `json:"file_type"`
	UpdatedAt string `json:"updated_at"`
	Deleted   bool   `json:"deleted,omitempty"`
}

type FileLogResponse struct {
	ID        string                 `json:"file_id"`
	ActorID   string                 `json:"actor_id"`
	ActorType string                 `json:"actor_type"`
	Action    string                 `json:"action"`
	IP        string                 `json:"ip"`
	Timestamp time.Time              `json:"timestamp"`
	UserAgent string                 `json:"user_agent"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type FileMetadataResponse struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"file_name,omitempty"`
	Size       int64  `json:"file_size,omitempty"`
	MimeType   string `json:"file_type,omitempty"`
	VersionID  string `json:"version_id,omitempty"`
	Hash       string `json:"hash,omitempty"`
	BucketName string `json:"bucket,omitempty"`
	Location   string `json:"location,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}
