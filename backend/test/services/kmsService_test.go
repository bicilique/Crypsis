package services_test

import (
	"context"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// mockKMSServer creates a test HTTP server that simulates KMS responses
type mockKMSServer struct {
	server       *httptest.Server
	responses    map[string]interface{}
	requestCount int
}

func newMockKMSServer() *mockKMSServer {
	mock := &mockKMSServer{
		responses: make(map[string]interface{}),
	}

	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.requestCount++

		// Verify it's a POST request
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Read and parse request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Determine operation type from request tag
		tag, _ := reqBody["tag"].(string)

		// Generate appropriate response based on operation
		var response interface{}
		switch tag {
		case "Create":
			response = createSymmetricKeyResponse()
		case "Locate":
			response = locateKeyResponse()
		case "Export":
			response = exportKeyResponse()
		case "Encrypt":
			response = encryptResponse()
		case "Decrypt":
			response = decryptResponse()
		case "Destroy":
			response = destroyKeyResponse()
		case "Revoke":
			response = revokeKeyResponse()
		case "ReKey":
			response = rekeyResponse()
		case "CreateKeyPair":
			response = createKeyPairResponse()
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))

	return mock
}

func (m *mockKMSServer) close() {
	m.server.Close()
}

// Mock response generators
func createSymmetricKeyResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "CreateResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "test-key-uid-12345"},
			{Tag: "ObjectType", Type: "Enumeration", Value: "SymmetricKey"},
		},
	}
}

func locateKeyResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "LocateResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "LocatedItems", Type: "Integer", Value: float64(2)},
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "key-uid-1"},
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "key-uid-2"},
		},
	}
}

func exportKeyResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "ExportResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "test-key-uid-12345"},
			{
				Tag:  "Object",
				Type: "Structure",
				Value: []interface{}{
					map[string]interface{}{
						"tag": "KeyBlock",
						"value": []interface{}{
							map[string]interface{}{
								"tag": "KeyValue",
								"value": []interface{}{
									map[string]interface{}{
										"tag": "KeyMaterial",
										"value": []interface{}{
											map[string]interface{}{
												"tag":   "ByteString",
												"value": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func encryptResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "EncryptResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "test-key-uid-12345"},
			{Tag: "Data", Type: "ByteString", Value: "encrypted-data-hex"},
			{Tag: "IvCounterNonce", Type: "ByteString", Value: "iv-hex-value"},
			{Tag: "AuthenticatedEncryptionTag", Type: "ByteString", Value: "auth-tag-hex"},
		},
	}
}

func decryptResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "DecryptResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "test-key-uid-12345"},
			{Tag: "Data", Type: "ByteString", Value: "decrypted-plaintext-hex"},
		},
	}
}

func destroyKeyResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "DestroyResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "test-key-uid-12345"},
		},
	}
}

func revokeKeyResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "RevokeResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "test-key-uid-12345"},
		},
	}
}

func rekeyResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "ReKeyResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "UniqueIdentifier", Type: "TextString", Value: "new-key-uid-67890"},
		},
	}
}

func createKeyPairResponse() model.KmsResponse {
	return model.KmsResponse{
		Tag:  "CreateKeyPairResponse",
		Type: "Structure",
		Value: []model.ValueResponse{
			{Tag: "PrivateKeyUniqueIdentifier", Type: "TextString", Value: "private-key-uid"},
			{Tag: "PublicKeyUniqueIdentifier", Type: "TextString", Value: "public-key-uid"},
		},
	}
}

// Tests

func TestNewKmsService(t *testing.T) {
	client := &http.Client{}
	kmsURL := "https://test-kms.example.com"

	service := services.NewKmsService(client, kmsURL)

	if service == nil {
		t.Fatal("Expected non-nil service instance")
	}
}

func TestGenerateSymmetricKey(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key generation", func(t *testing.T) {
		keyUID, err := service.GenerateSymetricKey(ctx, "test-app-key")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if keyUID != "test-key-uid-12345" {
			t.Errorf("Expected keyUID 'test-key-uid-12345', got: %s", keyUID)
		}
	})

	t.Run("empty key name", func(t *testing.T) {
		_, err := service.GenerateSymetricKey(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key name")
		}

		if !strings.Contains(err.Error(), "key name cannot be empty") {
			t.Errorf("Expected 'key name cannot be empty' error, got: %v", err)
		}
	})

	t.Run("whitespace key name", func(t *testing.T) {
		_, err := service.GenerateSymetricKey(ctx, "   ")
		if err == nil {
			t.Error("Expected error for whitespace key name")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := service.GenerateSymetricKey(ctx, "test-key")
		if err == nil {
			t.Error("Expected error for cancelled context")
		}
	})
}

func TestLocateKey(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key location", func(t *testing.T) {
		keyUIDs, err := service.LocateKey(ctx, "test-app-key")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(keyUIDs) != 2 {
			t.Errorf("Expected 2 key UIDs, got: %d", len(keyUIDs))
		}

		expectedUIDs := map[string]bool{"key-uid-1": true, "key-uid-2": true}
		for _, uid := range keyUIDs {
			if !expectedUIDs[uid] {
				t.Errorf("Unexpected key UID: %s", uid)
			}
		}
	})

	t.Run("empty key name", func(t *testing.T) {
		_, err := service.LocateKey(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key name")
		}
	})
}

func TestExportKey(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key export", func(t *testing.T) {
		keyMaterial, err := service.ExportKey(ctx, "test-key-uid-12345")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expectedMaterial := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
		if keyMaterial != expectedMaterial {
			t.Errorf("Expected key material %s, got: %s", expectedMaterial, keyMaterial)
		}
	})

	t.Run("empty key UID", func(t *testing.T) {
		_, err := service.ExportKey(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key UID")
		}
	})
}

func TestEncrypt(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful encryption", func(t *testing.T) {
		keyUID := "test-key-uid-12345"
		plaintext := "48656c6c6f" // "Hello" in hex

		encryptedData, iv, authTag, err := service.Encrypt(ctx, keyUID, plaintext)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if encryptedData == "" {
			t.Error("Expected non-empty encrypted data")
		}
		if iv == "" {
			t.Error("Expected non-empty IV")
		}
		if authTag == "" {
			t.Error("Expected non-empty auth tag")
		}

		if encryptedData != "encrypted-data-hex" {
			t.Errorf("Expected encrypted-data-hex, got: %s", encryptedData)
		}
		if iv != "iv-hex-value" {
			t.Errorf("Expected iv-hex-value, got: %s", iv)
		}
		if authTag != "auth-tag-hex" {
			t.Errorf("Expected auth-tag-hex, got: %s", authTag)
		}
	})

	t.Run("empty key UID", func(t *testing.T) {
		_, _, _, err := service.Encrypt(ctx, "", "48656c6c6f")
		if err == nil {
			t.Error("Expected error for empty key UID")
		}
	})

	t.Run("empty plaintext", func(t *testing.T) {
		_, _, _, err := service.Encrypt(ctx, "test-key-uid", "")
		if err == nil {
			t.Error("Expected error for empty plaintext")
		}
	})
}

func TestDecrypt(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful decryption", func(t *testing.T) {
		keyUID := "test-key-uid-12345"
		encryptedData := "encrypted-data-hex"
		iv := "iv-hex-value"
		authTag := "auth-tag-hex"

		decrypted, err := service.Decrypt(ctx, keyUID, encryptedData, iv, authTag)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if decrypted == "" {
			t.Error("Expected non-empty decrypted data")
		}

		if decrypted != "decrypted-plaintext-hex" {
			t.Errorf("Expected decrypted-plaintext-hex, got: %s", decrypted)
		}
	})

	t.Run("empty parameters", func(t *testing.T) {
		tests := []struct {
			name          string
			keyUID        string
			encryptedData string
			iv            string
			authTag       string
		}{
			{"empty keyUID", "", "data", "iv", "tag"},
			{"empty encryptedData", "uid", "", "iv", "tag"},
			{"empty iv", "uid", "data", "", "tag"},
			{"empty authTag", "uid", "data", "iv", ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := service.Decrypt(ctx, tt.keyUID, tt.encryptedData, tt.iv, tt.authTag)
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				}
			})
		}
	})
}

func TestDestroyKey(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key destruction", func(t *testing.T) {
		keyUID := "test-key-uid-12345"

		destroyedUID, err := service.DestroyKey(ctx, keyUID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if destroyedUID != keyUID {
			t.Errorf("Expected destroyed UID %s, got: %s", keyUID, destroyedUID)
		}
	})

	t.Run("empty key UID", func(t *testing.T) {
		_, err := service.DestroyKey(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key UID")
		}
	})
}

func TestRevokeKey(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key revocation", func(t *testing.T) {
		keyUID := "test-key-uid-12345"

		revokedUID, err := service.RevokeKey(ctx, keyUID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if revokedUID != keyUID {
			t.Errorf("Expected revoked UID %s, got: %s", keyUID, revokedUID)
		}
	})

	t.Run("empty key UID", func(t *testing.T) {
		_, err := service.RevokeKey(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key UID")
		}
	})
}

func TestReKey(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key rotation", func(t *testing.T) {
		oldKeyUID := "test-key-uid-12345"

		newKeyUID, err := service.ReKey(ctx, oldKeyUID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if newKeyUID == "" {
			t.Error("Expected non-empty new key UID")
		}

		if newKeyUID != "new-key-uid-67890" {
			t.Errorf("Expected new-key-uid-67890, got: %s", newKeyUID)
		}

		if newKeyUID == oldKeyUID {
			t.Error("New key UID should be different from old key UID")
		}
	})

	t.Run("empty key UID", func(t *testing.T) {
		_, err := service.ReKey(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key UID")
		}
	})
}

func TestGenerateKeyPair(t *testing.T) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	t.Run("successful key pair generation", func(t *testing.T) {
		name := "test-keypair"

		privateKeyUID, publicKeyUID, err := service.GenerateKeyPair(ctx, name)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if privateKeyUID == "" {
			t.Error("Expected non-empty private key UID")
		}
		if publicKeyUID == "" {
			t.Error("Expected non-empty public key UID")
		}

		if privateKeyUID == publicKeyUID {
			t.Error("Private and public key UIDs should be different")
		}

		if privateKeyUID != "private-key-uid" {
			t.Errorf("Expected private-key-uid, got: %s", privateKeyUID)
		}
		if publicKeyUID != "public-key-uid" {
			t.Errorf("Expected public-key-uid, got: %s", publicKeyUID)
		}
	})

	t.Run("empty key name", func(t *testing.T) {
		_, _, err := service.GenerateKeyPair(ctx, "")
		if err == nil {
			t.Error("Expected error for empty key name")
		}
	})
}

func TestContextTimeout(t *testing.T) {
	// Create a server that delays response
	delayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer delayServer.Close()

	client := &http.Client{}
	service := services.NewKmsService(client, delayServer.URL)

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := service.GenerateSymetricKey(ctx, "test-key")
	if err == nil {
		t.Error("Expected timeout error")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "KMS request failed") {
		t.Errorf("Expected context deadline or request failed error, got: %v", err)
	}
}

func TestKMSServerError(t *testing.T) {
	// Create a server that returns errors
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Internal server error"}`)
	}))
	defer errorServer.Close()

	client := &http.Client{}
	service := services.NewKmsService(client, errorServer.URL)
	ctx := context.Background()

	t.Run("server error on key generation", func(t *testing.T) {
		_, err := service.GenerateSymetricKey(ctx, "test-key")
		if err == nil {
			t.Error("Expected error for server error response")
		}

		if !strings.Contains(err.Error(), "500") {
			t.Errorf("Expected error to contain status code 500, got: %v", err)
		}
	})
}

func TestInvalidJSONResponse(t *testing.T) {
	// Create a server that returns invalid JSON
	invalidJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{invalid json response`)
	}))
	defer invalidJSONServer.Close()

	client := &http.Client{}
	service := services.NewKmsService(client, invalidJSONServer.URL)
	ctx := context.Background()

	_, err := service.GenerateSymetricKey(ctx, "test-key")
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}

	if !strings.Contains(err.Error(), "parse") && !strings.Contains(err.Error(), "JSON") {
		t.Errorf("Expected JSON parsing error, got: %v", err)
	}
}

// Benchmark tests
func BenchmarkGenerateSymmetricKey(b *testing.B) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateSymetricKey(ctx, fmt.Sprintf("bench-key-%d", i))
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	mockServer := newMockKMSServer()
	defer mockServer.close()

	client := &http.Client{}
	service := services.NewKmsService(client, mockServer.server.URL)
	ctx := context.Background()

	keyUID := "bench-key-uid"
	plaintext := "48656c6c6f20576f726c64" // "Hello World" in hex

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encData, iv, authTag, err := service.Encrypt(ctx, keyUID, plaintext)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}

		_, err = service.Decrypt(ctx, keyUID, encData, iv, authTag)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}
