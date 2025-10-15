package model

type MinIOConfig struct {
	Endpoint        string `json:"endpoint" validate:"required,url"`
	AccessKeyID     string `json:"access_key_id" validate:"required"`
	SecretAccessKey string `json:"secret_access_key" validate:"required"`
	UseSSL          bool   `json:"use_ssl"`
	BucketName      string `json:"bucket_name" validate:"required"`
	Region          string `json:"region" validate:"required"`
}

type StorageTransactionResponse struct {
	VersionID      string
	LastModified   string
	Expiration     string
	Location       string
	ChecksumSHA256 string
	IsLatest       bool
	IsDeleteMarker bool
}
