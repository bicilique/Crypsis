package helper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/xtgo/uuid"
)

// GetFileSize returns the size of a multipart.File.
func GetFileSize(file multipart.File) (int64, error) {
	size, err := file.Seek(0, io.SeekEnd)
	if err == nil {
		_, _ = file.Seek(0, io.SeekStart)
		return size, nil
	}

	buf := new(bytes.Buffer)
	size, err = io.Copy(buf, file)
	if err != nil {
		slog.Error("Failed to read file into buffer", slog.Any("error", err))
		return 0, errors.New("failed to read file into buffer")
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		slog.Error("Failed to reset file cursor", slog.Any("error", err))
		return 0, errors.New("failed to reset file cursor")
	}

	return size, nil
}

// GetFileBytesFromMultipart reads all bytes from a multipart.File and returns the bytes, size, and MIME type.
func GetFileBytesFromMultipart(file multipart.File) ([]byte, int64, string, error) {
	// Default MIME type
	mimeType := "application/octet-stream"

	// Read first 512 bytes for MIME detection
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err == nil || err == io.EOF { // If reading succeeds or reaches EOF
		mimeType = http.DetectContentType(header[:n])
	} else {
		slog.Error("Failed to read file header", slog.Any("error", err))
	}

	// Seek back before reading the whole file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		slog.Error("Failed to reset file cursor", slog.Any("error", err))
		return nil, 0, "", errors.New("failed to reset file cursor")
	}

	// Read full file
	buf, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file bytes", slog.Any("error", err))
		return nil, 0, "", errors.New("failed to read file bytes")
	}

	size := int64(len(buf))

	return buf, size, mimeType, nil
}

// GetFileBytesFromFile reads all bytes from an os.File and returns the bytes and size.
func GetFileBytesFromFile(file *os.File) ([]byte, int64, error) {
	buf, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file bytes", slog.Any("error", err))
		return nil, 0, errors.New("failed to read file bytes")
	}

	size := int64(len(buf))

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		slog.Error("Failed to reset file cursor", slog.Any("error", err))
		return nil, 0, errors.New("failed to reset file cursor")
	}

	return buf, size, nil
}

// CreateMultipartFileFromBytes creates a multipart.File from byte content.
func CreateMultipartFileFromBytes(fileContent []byte, fileName string) (multipart.File, int64, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, 0, errors.New("failed to create form file: " + err.Error())
	}

	size, err := io.Copy(part, bytes.NewReader(fileContent))
	if err != nil {
		return nil, 0, errors.New("failed to copy file content: " + err.Error())
	}

	if err := writer.Close(); err != nil {
		return nil, 0, errors.New("failed to close multipart writer: " + err.Error())
	}

	form, err := multipart.NewReader(&buf, writer.Boundary()).ReadForm(10 << 20)
	if err != nil {
		return nil, 0, errors.New("failed to read form: " + err.Error())
	}

	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return nil, 0, errors.New("file not found in form")
	}

	multipartFile, err := fileHeaders[0].Open()
	if err != nil {
		return nil, 0, errors.New("failed to open file: " + err.Error())
	}

	return multipartFile, size, nil
}

// EncodeToBase64 encodes bytes to a base64 string.
func EncodeToBase64(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("data is empty")
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodeFromBase64 decodes a base64-encoded string to bytes.
func DecodeFromBase64(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, errors.New("encoded string is empty")
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("failed to decode base64: " + err.Error())
	}
	return decoded, nil
}

// HexToBase64 converts a hex-encoded string to a base64-encoded string.
func HexToBase64(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(bytes)
	return base64Str, nil
}

// FileToBase64 reads a file from the given path and encodes its content to a base64 string.
func FileToBase64(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(fileBytes)
	return base64Str, nil
}

func GenerateCustomUUID() uuid.UUID {
	uuidBytes := make([]byte, 16)

	// Read 16 random bytes
	_, err := rand.Read(uuidBytes)
	if err != nil {
		return uuid.UUID{}
	}

	// Set version (4) and variant (RFC 4122) bits
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40 // Version 4
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // Variant is RFC 4122

	// Convert bytes to UUID type
	secureUUID := uuid.FromBytes(uuidBytes)
	return secureUUID
}
