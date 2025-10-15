package services

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/awnumar/memguard"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/prf"
)

const (
	HashSHA256 = "SHA-256"
	HashMD5    = "MD5"
)

type CryptographicService struct{}

// NewCryptographicService creates a new instance of CryptographicService
func NewCryptographicService() CryptographicInterface {
	return &CryptographicService{}
}

// GenerateKey generates a new AES-GCM key using Tink and returns it as a base64-encoded string
func (c *CryptographicService) GenerateKey() (string, error) {
	// Create a new keyset handle using AES256-GCM
	handle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		slog.Error("Failed to generate key", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	// Serialize the keyset to bytes
	buf := new(bytes.Buffer)
	writer := keyset.NewBinaryWriter(buf)
	if err := insecurecleartextkeyset.Write(handle, writer); err != nil {
		slog.Error("Failed to serialize key", slog.Any("error", err))
		return "", fmt.Errorf("failed to serialize key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// EncryptString encrypts the provided plaintext string using Tink AEAD
func (c *CryptographicService) EncryptString(key, text string) (string, error) {
	memguardKey := memguard.NewBufferFromBytes([]byte(key))
	defer memguardKey.Destroy()

	keyBytes, err := base64.StdEncoding.DecodeString(memguardKey.String())
	if err != nil {
		return "", fmt.Errorf("invalid base64 key: %w", err)
	}

	secureKeyBytes := memguard.NewBufferFromBytes(keyBytes)
	defer secureKeyBytes.Destroy()

	reader := keyset.NewBinaryReader(bytes.NewReader(secureKeyBytes.Bytes()))
	handle, err := insecurecleartextkeyset.Read(reader)
	if err != nil {
		slog.Error("Failed to read keyset", slog.Any("error", err))
		return "", fmt.Errorf("failed to read keyset: %w", err)
	}

	primitive, err := aead.New(handle)
	if err != nil {
		slog.Error("Failed to get AEAD primitive", slog.Any("error", err))
		return "", fmt.Errorf("failed to get AEAD primitive: %w", err)
	}

	// Encrypt the plaintext
	ciphertext, err := primitive.Encrypt([]byte(text), nil)
	if err != nil {
		slog.Error("Failed to encrypt", slog.Any("error", err))
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString decrypts the provided ciphertext string using Tink AEAD
func (c *CryptographicService) DecryptString(key, text string) (string, error) {
	memguardKey := memguard.NewBufferFromBytes([]byte(key))
	defer memguardKey.Destroy()

	keyBytes, err := base64.StdEncoding.DecodeString(memguardKey.String())
	if err != nil {
		slog.Error("Invalid base64 key", slog.Any("error", err))
		return "", fmt.Errorf("invalid base64 key: %w", err)
	}

	secureKeyBytes := memguard.NewBufferFromBytes(keyBytes)
	defer secureKeyBytes.Destroy()

	reader := keyset.NewBinaryReader(bytes.NewReader(secureKeyBytes.Bytes()))
	handle, err := insecurecleartextkeyset.Read(reader)
	if err != nil {
		slog.Error("Failed to read keyset", slog.Any("error", err))
		return "", fmt.Errorf("failed to read keyset: %w", err)
	}

	primitive, err := aead.New(handle)
	if err != nil {
		slog.Error("Failed to get AEAD primitive", slog.Any("error", err))
		return "", fmt.Errorf("failed to get AEAD primitive: %w", err)
	}

	// Decode the ciphertext from base64
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		slog.Error("Invalid base64 encrypted text", slog.Any("error", err))
		return "", fmt.Errorf("invalid base64 encrypted text: %w", err)
	}

	// Decrypt the ciphertext
	plaintext, err := primitive.Decrypt(ciphertext, nil)
	if err != nil {
		slog.Error("Failed to decrypt", slog.Any("error", err))
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	slog.Info("Decryption successful")
	return string(plaintext), nil
}

// HashString generates a hash of the given text using the specified hash method
func (c *CryptographicService) HashString(hashMethod, text string) (string, error) {
	var hashBytes []byte
	slog.Info("Hashing string")
	switch hashMethod {
	case HashSHA256:
		hash := sha256.Sum256([]byte(text))
		hashBytes = hash[:]
	case HashMD5:
		hash := md5.Sum([]byte(text))
		hashBytes = hash[:]
	default:
		slog.Error("Unsupported hash algorithm", slog.Any("algorithm", hashMethod))
		return "", errors.New("unsupported hash algorithm")
	}
	slog.Info("Hashing successful")
	return base64.StdEncoding.EncodeToString(hashBytes), nil
}

// CompareHash compares a hash with the hash of the given text
func (c *CryptographicService) CompareHash(hashMethod, text, hash string) bool {
	hashedText, err := c.HashString(hashMethod, text)
	if err != nil {
		return false
	}
	return hashedText == hash
}

// EncryptFile encrypts a file using Tink AEAD
func (c *CryptographicService) EncryptFile(key string, file []byte) ([]byte, error) {
	// Secure key string with memguard
	memguardKey := memguard.NewBufferFromBytes([]byte(key))
	defer memguardKey.Destroy()

	keyBytes, err := base64.StdEncoding.DecodeString(memguardKey.String())
	if err != nil {
		return nil, fmt.Errorf("invalid base64 key: %w", err)
	}

	// Secure keyBytes with memguard
	secureKeyBytes := memguard.NewBufferFromBytes(keyBytes)
	defer secureKeyBytes.Destroy()

	reader := keyset.NewBinaryReader(bytes.NewReader(secureKeyBytes.Bytes()))
	handle, err := insecurecleartextkeyset.Read(reader)
	if err != nil {
		slog.Error("Failed to read keyset", slog.Any("error", err))
		return nil, fmt.Errorf("failed to read keyset: %w", err)
	}

	primitive, err := aead.New(handle)
	if err != nil {
		slog.Error("Failed to get AEAD primitive", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get AEAD primitive: %w", err)
	}

	ciphertext, err := primitive.Encrypt(file, nil)
	if err != nil {
		slog.Error("Failed to encrypt file", slog.Any("error", err))
		return nil, fmt.Errorf("failed to encrypt file: %w", err)
	}
	return ciphertext, nil
}

// DecryptFile decrypts a file using Tink AEAD
func (c *CryptographicService) DecryptFile(key string, encryptedFile []byte) ([]byte, error) {
	memguardKey := memguard.NewBufferFromBytes([]byte(key))
	defer memguardKey.Destroy()

	keyBytes, err := base64.StdEncoding.DecodeString(memguardKey.String())
	if err != nil {
		return nil, fmt.Errorf("invalid base64 key: %w", err)
	}

	secureKeyBytes := memguard.NewBufferFromBytes(keyBytes)
	defer secureKeyBytes.Destroy()

	reader := keyset.NewBinaryReader(bytes.NewReader(secureKeyBytes.Bytes()))
	handle, err := insecurecleartextkeyset.Read(reader)
	if err != nil {
		slog.Error("Failed to read keyset", slog.Any("error", err))
		return nil, fmt.Errorf("failed to read keyset: %w", err)
	}

	primitive, err := aead.New(handle)
	if err != nil {
		slog.Error("Failed to get AEAD primitive", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get AEAD primitive: %w", err)
	}

	// Decrypt the file
	plaintext, err := primitive.Decrypt(encryptedFile, nil)
	if err != nil {
		slog.Error("Failed to decrypt file", slog.Any("error", err))
		return nil, fmt.Errorf("failed to decrypt file: %w", err)
	}

	slog.Info("File decryption successful")
	return plaintext, nil
}

// HashFile generates a hash of the file content
func (c *CryptographicService) HashFile(hashMethod string, file []byte) (string, error) {
	slog.Info("Hashing file content")
	var hashBytes []byte
	switch hashMethod {
	case HashSHA256:
		hash := sha256.Sum256(file)
		hashBytes = hash[:]
	case HashMD5:
		hash := md5.Sum(file)
		hashBytes = hash[:]
	default:
		slog.Error("Unsupported hash algorithm", slog.Any("algorithm", hashMethod))
		return "", errors.New("unsupported hash algorithm")
	}

	slog.Info("File hash generated successfully")
	return base64.StdEncoding.EncodeToString(hashBytes), nil
}

// CompareHashFile compares a hash with the hash of the given file
func (c *CryptographicService) CompareHashFile(hashMethod string, file []byte, hash string) bool {
	slog.Info("Comparing file hash")
	fileHash, err := c.HashFile(hashMethod, file)
	if err != nil {
		return false
	}
	return fileHash == hash
}

// KeyDerivationFunction derives a key from the input and salt using Tink PRF
func (c *CryptographicService) KeyDerivationFunction(input string, salt []byte) ([]byte, error) {
	// Create a new PRF keyset
	handle, err := keyset.NewHandle(prf.HKDFSHA256PRFKeyTemplate())
	if err != nil {
		slog.Error("Failed to create PRF keyset", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create PRF keyset: %w", err)
	}

	// Get PRF set
	prfSet, err := prf.NewPRFSet(handle)
	if err != nil {
		slog.Error("Failed to get PRF set", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get PRF set: %w", err)
	}

	// Get primary PRF - PrimaryID is a field, not a method
	primaryID := prfSet.PrimaryID

	// PRFs is a field, not a method
	prfs := prfSet.PRFs
	primaryPRF, ok := prfs[primaryID]
	if !ok {
		return nil, errors.New("primary PRF not found")
	}

	// Compute PRF with input and salt combined
	inputData := append([]byte(input), salt...)
	derivedKey, err := primaryPRF.ComputePRF(inputData, 32) // 32 bytes = 256 bits
	if err != nil {
		slog.Error("Failed to compute PRF", slog.Any("error", err))
		return nil, fmt.Errorf("failed to compute PRF: %w", err)
	}

	slog.Info("Key derivation successful")
	return derivedKey, nil
}
