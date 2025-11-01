package services

import (
	"bytes"
	"context"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var (
	// ErrInvalidInput is returned when input parameters are invalid
	ErrInvalidInput = errors.New("invalid input parameter")
	// ErrKMSRequest is returned when KMS request fails
	ErrKMSRequest = errors.New("KMS request failed")
	// ErrKMSResponse is returned when KMS response parsing fails
	ErrKMSResponse = errors.New("KMS response parsing failed")
	// ErrKeyNotFound is returned when key cannot be found
	ErrKeyNotFound = errors.New("key not found")
)

// KmsService provides cryptographic key management operations using KMIP protocol.
// It handles symmetric keys, key pairs, encryption, decryption, and key lifecycle management.
type KmsService struct {
	secureClient *http.Client
	kmsURL       string
}

// NewKmsService creates a new instance of KmsService with the provided HTTPS client and KMS URL.
//
// Parameters:
//   - secureClient: HTTP client configured with TLS/SSL certificates for secure communication
//   - kmsURL: Base URL of the KMS server (e.g., "https://kms.example.com")
//
// Returns:
//   - KMSInterface: Interface for KMS operations
func NewKmsService(secureClient *http.Client, kmsURL string) KMSInterface {
	return &KmsService{
		secureClient: secureClient,
		kmsURL:       kmsURL,
	}
}

// GenerateSymetricKey creates a new symmetric AES-256 key in the KMS with the specified tag name.
//
// This method generates a symmetric key suitable for encryption/decryption operations.
// The key is created with AES-256 algorithm and stored in the KMS with the provided tag name.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - name: Tag name to identify the key (must not be empty)
//
// Returns:
//   - string: Unique identifier (UID) of the generated key
//   - error: Error if key generation fails
//
// Example:
//
//	keyUID, err := kmsService.GenerateSymetricKey(ctx, "app-encryption-key")
func (s *KmsService) GenerateSymetricKey(ctx context.Context, name string) (string, error) {
	// Start tracing span
	tracer := helper.GetTracingHelper()
	ctx, span := tracer.StartKMSSpan(ctx, "GenerateSymmetricKey", name)
	defer span.End()

	// Validate input
	if strings.TrimSpace(name) == "" {
		err := fmt.Errorf("%w: key name cannot be empty", ErrInvalidInput)
		helper.RecordError(span, err)
		return "", err
	}

	// Generate key template
	jsonBody, err := helper.GenerateKeyTemplate(name)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate key template", slog.String("name", name), slog.Any("error", err))
		helper.RecordError(span, err)
		return "", fmt.Errorf("failed to generate key template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		helper.RecordError(span, err)
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		helper.RecordError(span, err)
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract UniqueIdentifier
	uniqueIdentifier, err := extractUniqueIdentifier(kmsResp)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract UniqueIdentifier", slog.Any("error", err))
		helper.RecordError(span, err)
		return "", err
	}

	helper.RecordSuccess(span, "Key generated successfully")
	// slog.InfoContext(ctx, "Successfully generated symmetric key", slog.String("keyUID", uniqueIdentifier), slog.String("name", name))
	return uniqueIdentifier, nil
}

// LocateKey finds all keys with the specified tag name and returns their unique identifiers.
//
// This method searches for keys in the KMS that match the provided tag name.
// It can return multiple key UIDs if multiple keys exist with the same tag.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - name: Tag name to search for (must not be empty)
//
// Returns:
//   - []string: List of unique identifiers for matching keys
//   - error: Error if search fails or no keys are found
//
// Example:
//
//	keyUIDs, err := kmsService.LocateKey(ctx, "app-encryption-key")
func (s *KmsService) LocateKey(ctx context.Context, name string) ([]string, error) {
	// Validate input
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("%w: key name cannot be empty", ErrInvalidInput)
	}

	// Generate locate template
	jsonBody, err := helper.GenerateLocateKeyTemplate(name)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate locate template", slog.String("name", name), slog.Any("error", err))
		return nil, fmt.Errorf("failed to generate locate template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return nil, fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract LocatedItems and UniqueIdentifiers
	var locatedItems int
	var uniqueIdentifiers []string

	for _, v := range kmsResp.Value {
		switch v.Tag {
		case "LocatedItems":
			// Ensure correct type assertion for integer
			if count, ok := v.Value.(float64); ok {
				locatedItems = int(count)
			}
		case "UniqueIdentifier":
			// If UniqueIdentifier is an array
			if list, ok := v.Value.([]interface{}); ok {
				for _, item := range list {
					if obj, valid := item.(map[string]interface{}); valid {
						if id, exists := obj["value"].(string); exists {
							uniqueIdentifiers = append(uniqueIdentifiers, id)
						}
					}
				}
			} else if single, ok := v.Value.(string); ok {
				// If UniqueIdentifier is a single string
				uniqueIdentifiers = append(uniqueIdentifiers, single)
			}
		}
	}

	if locatedItems != len(uniqueIdentifiers) {
		slog.WarnContext(ctx, "Mismatch between located items and extracted identifiers",
			slog.Int("LocatedItems", locatedItems),
			slog.Int("Extracted", len(uniqueIdentifiers)),
		)
	}

	if len(uniqueIdentifiers) == 0 {
		slog.ErrorContext(ctx, "No UniqueIdentifiers found in response", slog.String("name", name))
		return nil, fmt.Errorf("%w: no keys found with name '%s'", ErrKeyNotFound, name)
	}

	slog.InfoContext(ctx, "Successfully located keys", slog.Int("count", len(uniqueIdentifiers)), slog.String("name", name))
	return uniqueIdentifiers, nil
}

// ExportKey exports the key material for the specified key UID.
//
// This method retrieves the raw key material in hexadecimal format from the KMS.
// The key must be exportable and the caller must have appropriate permissions.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the key to export (must not be empty)
//
// Returns:
//   - string: Hexadecimal representation of the key material
//   - error: Error if export fails
//
// Example:
//
//	keyMaterial, err := kmsService.ExportKey(ctx, "key-uid-12345")
func (s *KmsService) ExportKey(ctx context.Context, keyUID string) (string, error) {
	// Start tracing span
	tracer := helper.GetTracingHelper()
	ctx, span := tracer.StartKMSSpan(ctx, "ExportKey", keyUID)
	defer span.End()

	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		err := fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
		helper.RecordError(span, err)
		return "", err
	}

	// Generate export template
	jsonBody, err := helper.GenerateExportTemplate(keyUID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate export template", slog.String("keyUID", keyUID), slog.Any("error", err))
		helper.RecordError(span, err)
		return "", fmt.Errorf("failed to generate export template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		helper.RecordError(span, err)
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		helper.RecordError(span, err)
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract key material from nested structure
	keyMaterial, err := extractKeyMaterial(kmsResp)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract key material", slog.String("keyUID", keyUID), slog.Any("error", err))
		helper.RecordError(span, err)
		return "", err
	}

	helper.RecordSuccess(span, "Key exported successfully")
	// slog.InfoContext(ctx, "Successfully exported key", slog.String("keyUID", keyUID))
	return keyMaterial, nil
}

// Encrypt encrypts the provided plaintext using the specified key with AES-GCM mode.
//
// This method performs authenticated encryption and returns the encrypted data along with
// the initialization vector (IV) and authentication tag required for decryption.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the encryption key (must not be empty)
//   - text: Hexadecimal string representation of the plaintext to encrypt (must not be empty)
//
// Returns:
//   - encryptedData: Hexadecimal representation of encrypted data
//   - iv: Hexadecimal representation of the initialization vector
//   - authTag: Hexadecimal representation of the authentication tag
//   - error: Error if encryption fails
//
// Example:
//
//	encrypted, iv, authTag, err := kmsService.Encrypt(ctx, "key-uid-12345", "48656c6c6f")
func (s *KmsService) Encrypt(ctx context.Context, keyUID string, text string) (string, string, string, error) {
	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		return "", "", "", fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(text) == "" {
		return "", "", "", fmt.Errorf("%w: text cannot be empty", ErrInvalidInput)
	}

	// Generate encrypt template
	jsonBody, err := helper.GenerateEncryptTemplate(keyUID, text)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate encrypt template", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", "", "", fmt.Errorf("failed to generate encrypt template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", "", "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return "", "", "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract encryption values with type safety
	var encryptedData, iv, authTag string
	for _, v := range kmsResp.Value {
		switch v.Tag {
		case "Data":
			if val, ok := v.Value.(string); ok {
				encryptedData = val
			}
		case "IvCounterNonce":
			if val, ok := v.Value.(string); ok {
				iv = val
			}
		case "AuthenticatedEncryptionTag":
			if val, ok := v.Value.(string); ok {
				authTag = val
			}
		}
	}

	// Check if all required fields were found
	if encryptedData == "" || iv == "" || authTag == "" {
		slog.ErrorContext(ctx, "Missing encryption response data", slog.Any("response", kmsResp))
		return "", "", "", fmt.Errorf("%w: missing encryption response fields (data=%v, iv=%v, authTag=%v)",
			ErrKMSResponse, encryptedData != "", iv != "", authTag != "")
	}

	slog.InfoContext(ctx, "Successfully encrypted data", slog.String("keyUID", keyUID))
	return encryptedData, iv, authTag, nil
}

// Decrypt decrypts the encrypted data using the specified key and authentication parameters.
//
// This method performs authenticated decryption using AES-GCM mode. It requires the encrypted data,
// initialization vector, and authentication tag that were returned by the Encrypt method.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the decryption key (must not be empty)
//   - encryptedData: Hexadecimal representation of encrypted data (must not be empty)
//   - ivCounterNonce: Hexadecimal representation of the initialization vector (must not be empty)
//   - authTag: Hexadecimal representation of the authentication tag (must not be empty)
//
// Returns:
//   - string: Hexadecimal representation of the decrypted plaintext
//   - error: Error if decryption fails or authentication check fails
//
// Example:
//
//	decrypted, err := kmsService.Decrypt(ctx, "key-uid-12345", encData, iv, authTag)
func (s *KmsService) Decrypt(ctx context.Context, keyUID, encryptedData, ivCounterNonce, authTag string) (string, error) {
	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		return "", fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(encryptedData) == "" {
		return "", fmt.Errorf("%w: encryptedData cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(ivCounterNonce) == "" {
		return "", fmt.Errorf("%w: ivCounterNonce cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(authTag) == "" {
		return "", fmt.Errorf("%w: authTag cannot be empty", ErrInvalidInput)
	}

	// Generate decrypt template
	jsonBody, err := helper.GenerateDecryptTemplate(keyUID, encryptedData, ivCounterNonce, authTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate decrypt template", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", fmt.Errorf("failed to generate decrypt template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract decrypted data
	var decryptedData string
	for _, v := range kmsResp.Value {
		if v.Tag == "Data" {
			if data, ok := v.Value.(string); ok {
				decryptedData = data
				break
			} else {
				slog.ErrorContext(ctx, "Unexpected type for Data field", slog.Any("value", v.Value))
				return "", fmt.Errorf("%w: unexpected type for Data field", ErrKMSResponse)
			}
		}
	}

	// Check if decrypted data was found
	if decryptedData == "" {
		slog.ErrorContext(ctx, "Decrypted data not found in response")
		return "", fmt.Errorf("%w: decrypted data not found in response", ErrKMSResponse)
	}

	slog.InfoContext(ctx, "Successfully decrypted data", slog.String("keyUID", keyUID))
	return decryptedData, nil
}

// DestroyKey permanently deletes the specified key from the KMS.
//
// This operation is irreversible. Once a key is destroyed, it cannot be recovered
// and any data encrypted with this key cannot be decrypted.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the key to destroy (must not be empty)
//
// Returns:
//   - string: UID of the destroyed key (for confirmation)
//   - error: Error if destruction fails
//
// Example:
//
//	destroyedUID, err := kmsService.DestroyKey(ctx, "key-uid-12345")
func (s *KmsService) DestroyKey(ctx context.Context, keyUID string) (string, error) {
	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		return "", fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
	}

	// Generate destroy template
	jsonBody, err := helper.GenerateDestroyTemplate(keyUID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate destroy template", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", fmt.Errorf("failed to generate destroy template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract UniqueIdentifier
	destroyedKeyUID, err := extractUniqueIdentifier(kmsResp)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract UniqueIdentifier", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", err
	}

	slog.InfoContext(ctx, "Successfully destroyed key", slog.String("keyUID", destroyedKeyUID))
	return destroyedKeyUID, nil
}

// RevokeKey marks the specified key as revoked in the KMS.
//
// Unlike DestroyKey, revoked keys are not deleted but marked as compromised or no longer trusted.
// This allows for audit trails while preventing future use of the key.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the key to revoke (must not be empty)
//
// Returns:
//   - string: UID of the revoked key (for confirmation)
//   - error: Error if revocation fails
//
// Example:
//
//	revokedUID, err := kmsService.RevokeKey(ctx, "key-uid-12345")
func (s *KmsService) RevokeKey(ctx context.Context, keyUID string) (string, error) {
	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		return "", fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
	}

	// Generate revoke template
	jsonBody, err := helper.GenerateRevokeTemplate(keyUID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate revoke template", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", fmt.Errorf("failed to generate revoke template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract UniqueIdentifier
	revokedKeyUID, err := extractUniqueIdentifier(kmsResp)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract UniqueIdentifier", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", err
	}

	slog.InfoContext(ctx, "Successfully revoked key", slog.String("keyUID", revokedKeyUID))
	return revokedKeyUID, nil
}

// ReKey creates a new version of the specified key for key rotation purposes.
//
// This operation generates a new key that can be used to replace the old key,
// supporting key rotation policies and security best practices.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the key to rotate (must not be empty)
//
// Returns:
//   - string: UID of the new key version
//   - error: Error if re-keying fails
//
// Example:
//
//	newKeyUID, err := kmsService.ReKey(ctx, "key-uid-12345")
func (s *KmsService) ReKey(ctx context.Context, keyUID string) (string, error) {
	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		return "", fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
	}

	// Generate rekey template
	jsonBody, err := helper.GenerateReKeyTemplate(keyUID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate rekey template", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", fmt.Errorf("failed to generate rekey template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract UniqueIdentifier
	newKeyUID, err := extractUniqueIdentifier(kmsResp)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract UniqueIdentifier", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", err
	}

	slog.InfoContext(ctx, "Successfully re-keyed", slog.String("oldKeyUID", keyUID), slog.String("newKeyUID", newKeyUID))
	return newKeyUID, nil
}

// Covercrypt encrypts data using CoverCrypt algorithm for policy-based encryption.
//
// CoverCrypt is an advanced encryption scheme that allows fine-grained access control
// based on encryption policies and attributes.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - keyUID: Unique identifier of the CoverCrypt public key (must not be empty)
//   - text: Hexadecimal string representation of the plaintext to encrypt (must not be empty)
//
// Returns:
//   - string: Hexadecimal representation of the encrypted data
//   - error: Error if encryption fails
//
// Example:
//
//	encrypted, err := kmsService.Covercrypt(ctx, "covercrypt-key-uid", "48656c6c6f")
func (s *KmsService) Covercrypt(ctx context.Context, keyUID string, text string) (string, error) {
	// Validate input
	if strings.TrimSpace(keyUID) == "" {
		return "", fmt.Errorf("%w: keyUID cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("%w: text cannot be empty", ErrInvalidInput)
	}

	// Generate encryption request template
	jsonBody, err := helper.GenerateCoverCryptEncryptTemplate(text, keyUID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate encryption template", slog.String("keyUID", keyUID), slog.Any("error", err))
		return "", fmt.Errorf("failed to generate encryption template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err))
		return "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract encrypted data
	var encryptedData string
	for _, v := range kmsResp.Value {
		if v.Tag == "Data" {
			if data, ok := v.Value.(string); ok {
				encryptedData = data
				break
			} else {
				slog.ErrorContext(ctx, "Unexpected type for Data")
				return "", fmt.Errorf("%w: unexpected type for Data field", ErrKMSResponse)
			}
		}
	}

	// Ensure Data was found
	if encryptedData == "" {
		slog.ErrorContext(ctx, "Failed to extract encrypted data from response")
		return "", fmt.Errorf("%w: encrypted data not found in response", ErrKMSResponse)
	}

	slog.InfoContext(ctx, "Successfully encrypted data using CoverCrypt", slog.String("keyUID", keyUID))
	return encryptedData, nil
}

// GenerateKeyPair creates a new asymmetric key pair (ECDH) with the specified tag name.
//
// This method generates an ECDH key pair using CURVE25519 for key exchange operations.
// Both private and public keys are stored in the KMS with the same tag name.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - name: Tag name to identify the key pair (must not be empty)
//
// Returns:
//   - privateKeyUID: Unique identifier of the generated private key
//   - publicKeyUID: Unique identifier of the generated public key
//   - error: Error if key pair generation fails
//
// Example:
//
//	privUID, pubUID, err := kmsService.GenerateKeyPair(ctx, "ecdh-keypair")
func (s *KmsService) GenerateKeyPair(ctx context.Context, name string) (string, string, error) {
	// Validate input
	if strings.TrimSpace(name) == "" {
		return "", "", fmt.Errorf("%w: key name cannot be empty", ErrInvalidInput)
	}

	// Generate key template
	jsonBody, err := helper.GenerateKeyPairTemplate(name)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate key template", slog.String("name", name), slog.Any("error", err))
		return "", "", fmt.Errorf("failed to generate key template: %w", err)
	}

	// Send request
	body, err := s.sendRequest(ctx, jsonBody)
	if err != nil {
		return "", "", err
	}

	// Parse JSON response
	var kmsResp model.KmsResponse
	if err := json.Unmarshal(body, &kmsResp); err != nil {
		slog.ErrorContext(ctx, "Failed to parse JSON response", slog.Any("error", err), slog.Any("body", string(body)))
		return "", "", fmt.Errorf("%w: failed to parse JSON response: %v", ErrKMSResponse, err)
	}

	// Extract PrivateKey and PublicKey UniqueIdentifier
	var privateKeyUID, publicKeyUID string
	for _, v := range kmsResp.Value {
		switch v.Tag {
		case "PrivateKeyUniqueIdentifier":
			if id, ok := v.Value.(string); ok {
				privateKeyUID = id
			} else {
				slog.ErrorContext(ctx, "Unexpected type for PrivateKeyUniqueIdentifier", slog.Any("value", v.Value))
				return "", "", fmt.Errorf("%w: unexpected type for PrivateKeyUniqueIdentifier", ErrKMSResponse)
			}
		case "PublicKeyUniqueIdentifier":
			if id, ok := v.Value.(string); ok {
				publicKeyUID = id
			} else {
				slog.ErrorContext(ctx, "Unexpected type for PublicKeyUniqueIdentifier", slog.Any("value", v.Value))
				return "", "", fmt.Errorf("%w: unexpected type for PublicKeyUniqueIdentifier", ErrKMSResponse)
			}
		}
	}

	// Ensure both keys were found
	if privateKeyUID == "" || publicKeyUID == "" {
		slog.ErrorContext(ctx, "Failed to extract key identifiers from response", slog.Any("response", kmsResp))
		return "", "", fmt.Errorf("%w: failed to extract key identifiers (privateKey=%v, publicKey=%v)",
			ErrKMSResponse, privateKeyUID != "", publicKeyUID != "")
	}

	// slog.InfoContext(ctx, "Successfully generated key pair",
	// 	slog.String("privateKeyUID", privateKeyUID),
	// 	slog.String("publicKeyUID", publicKeyUID),
	// 	slog.String("name", name))
	return privateKeyUID, publicKeyUID, nil
}

// sendRequest sends an HTTP POST request to the KMS server with the provided JSON body.
//
// This is a private helper method that handles HTTP communication, error handling,
// and response validation for all KMS operations.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - jsonBody: JSON-encoded request body
//
// Returns:
//   - []byte: Response body from the KMS server
//   - error: Error if request fails or KMS returns non-200 status
func (s *KmsService) sendRequest(ctx context.Context, jsonBody string) ([]byte, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.kmsURL+"/kmip/2_1", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", slog.Any("error", err))
		return nil, fmt.Errorf("%w: failed to create request: %v", ErrKMSRequest, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.secureClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send request to KMS", slog.Any("error", err))
		return nil, fmt.Errorf("%w: failed to send request: %v", ErrKMSRequest, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", slog.Any("error", err))
		return nil, fmt.Errorf("%w: failed to read response body: %v", ErrKMSRequest, err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "KMS returned an error",
			slog.Int("status_code", resp.StatusCode),
			slog.String("response", string(body)))
		return nil, fmt.Errorf("%w: status=%d, response=%s", ErrKMSRequest, resp.StatusCode, string(body))
	}

	return body, nil
}

// Helper functions for parsing KMS responses

// extractUniqueIdentifier extracts the UniqueIdentifier from a KMS response.
func extractUniqueIdentifier(kmsResp model.KmsResponse) (string, error) {
	for _, v := range kmsResp.Value {
		if v.Tag == "UniqueIdentifier" {
			if id, ok := v.Value.(string); ok {
				return id, nil
			}
			return "", fmt.Errorf("%w: unexpected type for UniqueIdentifier", ErrKMSResponse)
		}
	}
	return "", fmt.Errorf("%w: UniqueIdentifier not found in response", ErrKMSResponse)
}

// extractKeyMaterial extracts the key material from a nested KMS export response.
func extractKeyMaterial(kmsResp model.KmsResponse) (string, error) {
	// Loop through values and extract key material
	for _, item := range kmsResp.Value {
		if item.Tag == "Object" {
			objectValue, ok := item.Value.([]interface{})
			if !ok {
				return "", fmt.Errorf("%w: invalid Object structure", ErrKMSResponse)
			}

			// Find "KeyBlock" in Object
			for _, obj := range objectValue {
				objMap, ok := obj.(map[string]interface{})
				if !ok || objMap["tag"] != "KeyBlock" {
					continue
				}

				keyBlockValue, ok := objMap["value"].([]interface{})
				if !ok {
					return "", fmt.Errorf("%w: invalid KeyBlock structure", ErrKMSResponse)
				}

				// Find "KeyValue" in KeyBlock
				for _, keyBlockItem := range keyBlockValue {
					keyBlockMap, ok := keyBlockItem.(map[string]interface{})
					if !ok || keyBlockMap["tag"] != "KeyValue" {
						continue
					}

					keyValueList, ok := keyBlockMap["value"].([]interface{})
					if !ok {
						return "", fmt.Errorf("%w: invalid KeyValue structure", ErrKMSResponse)
					}

					// Find "KeyMaterial" in KeyValue
					for _, keyValueItem := range keyValueList {
						keyValueMap, ok := keyValueItem.(map[string]interface{})
						if !ok || keyValueMap["tag"] != "KeyMaterial" {
							continue
						}

						keyMaterialList, ok := keyValueMap["value"].([]interface{})
						if !ok {
							return "", fmt.Errorf("%w: invalid KeyMaterial structure", ErrKMSResponse)
						}

						// Find "ByteString" in KeyMaterial
						for _, keyMaterialItem := range keyMaterialList {
							keyMaterialMap, ok := keyMaterialItem.(map[string]interface{})
							if !ok || keyMaterialMap["tag"] != "ByteString" {
								continue
							}

							byteString, ok := keyMaterialMap["value"].(string)
							if !ok {
								return "", fmt.Errorf("%w: ByteString not found or invalid", ErrKMSResponse)
							}

							return byteString, nil
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("%w: key material not found in response", ErrKMSResponse)
}
