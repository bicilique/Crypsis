package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

// Validarion Error
var (
	ErrInvalidInput = errors.New("invalid input")
)

// Admin Error
var (
	AdminErrNotFound         = errors.New("admin account not found")
	AdminErrWrongCredentials = errors.New("wrong credentials")
	AdminErrAlreadyExists    = errors.New("admin already exists")
)

// Crypto Error
var (
	ErrHashCalculationFailed = errors.New("hash calculation failed")
	ErrKeyGenerationFailed   = errors.New("key generation failed")
	ErrFileEncryptionFailed  = errors.New("file encryption failed")
	ErrFileUidOrKeyInvalid   = errors.New("file UID or key is invalid")
	ErrHashNotMatch          = errors.New("hash value does not match")
)

// File Error
var (
	ErrUnauthorizedFileAccess = errors.New("unauthorized file access")
	ErrFileNotFound           = errors.New("file not found")
	ErrFileIsEmpty            = errors.New("file is empty")
	ErrFileAlreadyExists      = errors.New("file already exists")
	ErrFailedToReadFile       = errors.New("failed to read file")
	ErrFileUploadFailed       = errors.New("file upload failed")
	ErrFileDownloadFailed     = errors.New("file download failed")
)

// KM Error
var (
	ErrKeyNotFound                = errors.New("key not found")
	ErrKeyAlreadyExists           = errors.New("key already exists")
	ErrKeyDeletionFailed          = errors.New("key deletion failed")
	ErrKeyUpdateFailed            = errors.New("key update failed")
	ErrKeyRecoveryFailed          = errors.New("key recovery failed")
	ErrFailedToGenerateKeyFromKMS = errors.New("failed to generate key from KMS")
	ErrFailedToImportKeyFromKMS   = errors.New("failed to import key from KMS")
	ErrFailedToExportKeyToKMS     = errors.New("failed to export key to KMS")
)

// APP error
var (
	ErrAppAlreadyExists = errors.New("app already exists")
	ErrAppNotFound      = errors.New("app not found")
	ErrAppNotActive     = errors.New("app is inactive")
	ErrAppAlreadyActive = errors.New("app is already active")
)

func WrapValidationError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var msgs []string
		for _, e := range ve {
			msgs = append(msgs, fmt.Sprintf("%s is %s", e.Field(), e.Tag()))
		}
		return fmt.Errorf("%w: %s", ErrInvalidInput, strings.Join(msgs, ", "))
	}
	return err
}
