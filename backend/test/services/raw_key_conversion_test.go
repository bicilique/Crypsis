package services

import (
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/services"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestKeysetFromRawAES256GCM(t *testing.T) {
	service := services.NewCryptographicService()

	t.Run("Valid 32-byte key", func(t *testing.T) {
		// Generate a random 32-byte key
		rawKey := make([]byte, 32)
		_, err := rand.Read(rawKey)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		// Convert to Tink keyset
		handle, err := service.KeysetFromRawAES256GCM(rawKey)
		if err != nil {
			t.Fatalf("KeysetFromRawAES256GCM failed: %v", err)
		}

		if handle == nil {
			t.Fatal("Expected non-nil keyset handle")
		}

		t.Logf("Successfully created keyset handle from raw key")
	})

	t.Run("Invalid key length - too short", func(t *testing.T) {
		rawKey := make([]byte, 16) // 16 bytes (AES-128), should fail for AES-256
		_, err := rand.Read(rawKey)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		_, err = service.KeysetFromRawAES256GCM(rawKey)
		if err == nil {
			t.Fatal("Expected error for invalid key length, got nil")
		}

		t.Logf("Correctly rejected invalid key length: %v", err)
	})

	t.Run("Invalid key length - too long", func(t *testing.T) {
		rawKey := make([]byte, 64) // 64 bytes, should fail
		_, err := rand.Read(rawKey)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		_, err = service.KeysetFromRawAES256GCM(rawKey)
		if err == nil {
			t.Fatal("Expected error for invalid key length, got nil")
		}

		t.Logf("Correctly rejected invalid key length: %v", err)
	})
}

func TestImportRawKeyAsBase64(t *testing.T) {
	service := services.NewCryptographicService()

	t.Run("Valid raw key import", func(t *testing.T) {
		// Generate a random 32-byte key
		rawKey := make([]byte, 32)
		_, err := rand.Read(rawKey)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		// Import as base64 Tink keyset
		keyBase64, err := service.ImportRawKeyAsBase64(rawKey)
		if err != nil {
			t.Fatalf("ImportRawKeyAsBase64 failed: %v", err)
		}

		if keyBase64 == "" {
			t.Fatal("Expected non-empty base64 key")
		}

		// Decode and verify it's valid base64
		decoded, err := base64.StdEncoding.DecodeString(keyBase64)
		if err != nil {
			t.Fatalf("Failed to decode base64 key: %v", err)
		}

		// Tink keysets are typically around 100+ bytes
		if len(decoded) < 50 {
			t.Fatalf("Decoded key too short: got %d bytes", len(decoded))
		}

		t.Logf("Successfully imported raw key as base64 (length: %d bytes)", len(decoded))
	})

	t.Run("Invalid key length", func(t *testing.T) {
		rawKey := make([]byte, 24) // Invalid length
		_, err := rand.Read(rawKey)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		_, err = service.ImportRawKeyAsBase64(rawKey)
		if err == nil {
			t.Fatal("Expected error for invalid key length, got nil")
		}

		t.Logf("Correctly rejected invalid key: %v", err)
	})
}

func TestRawKeyEncryptionDecryption(t *testing.T) {
	service := services.NewCryptographicService()

	t.Run("Encrypt and decrypt with raw key", func(t *testing.T) {
		// Generate a random 32-byte raw key
		rawKey := make([]byte, 32)
		_, err := rand.Read(rawKey)
		if err != nil {
			t.Fatalf("Failed to generate random key: %v", err)
		}

		// Convert to Tink keyset format
		keyBase64, err := service.ImportRawKeyAsBase64(rawKey)
		if err != nil {
			t.Fatalf("ImportRawKeyAsBase64 failed: %v", err)
		}

		// Test data
		plaintext := []byte("Hello, this is a test file content!")

		// Encrypt
		ciphertext, err := service.EncryptFile(keyBase64, plaintext)
		if err != nil {
			t.Fatalf("EncryptFile failed: %v", err)
		}

		if len(ciphertext) == 0 {
			t.Fatal("Expected non-empty ciphertext")
		}

		t.Logf("Encrypted %d bytes to %d bytes", len(plaintext), len(ciphertext))

		// Decrypt
		decrypted, err := service.DecryptFile(keyBase64, ciphertext)
		if err != nil {
			t.Fatalf("DecryptFile failed: %v", err)
		}

		// Verify
		if string(decrypted) != string(plaintext) {
			t.Fatalf("Decrypted data mismatch.\nExpected: %s\nGot: %s", plaintext, decrypted)
		}

		t.Log("Successfully encrypted and decrypted with raw key")
	})

	t.Run("Encrypt with raw key, decrypt with wrong key", func(t *testing.T) {
		// Generate two different raw keys
		rawKey1 := make([]byte, 32)
		rawKey2 := make([]byte, 32)
		_, err := rand.Read(rawKey1)
		require.NoError(t, err)
		_, err = rand.Read(rawKey2)
		require.NoError(t, err)

		key1Base64, _ := service.ImportRawKeyAsBase64(rawKey1)
		key2Base64, _ := service.ImportRawKeyAsBase64(rawKey2)

		plaintext := []byte("Secret message")

		// Encrypt with key1
		ciphertext, err := service.EncryptFile(key1Base64, plaintext)
		if err != nil {
			t.Fatalf("EncryptFile failed: %v", err)
		}

		// Try to decrypt with key2 (should fail)
		_, err = service.DecryptFile(key2Base64, ciphertext)
		if err == nil {
			t.Fatal("Expected decryption to fail with wrong key, but it succeeded")
		}

		t.Logf("Correctly failed to decrypt with wrong key: %v", err)
	})
}

func TestRawKeyVsTinkKeyCompatibility(t *testing.T) {
	service := services.NewCryptographicService()

	t.Run("Compare raw key vs native Tink key", func(t *testing.T) {
		// Generate a native Tink key
		tinkKey, err := service.GenerateKey()
		if err != nil {
			t.Fatalf("GenerateKey failed: %v", err)
		}

		// Generate a raw key and convert it
		rawKey := make([]byte, 32)
		rand.Read(rawKey)
		rawKeyConverted, err := service.ImportRawKeyAsBase64(rawKey)
		if err != nil {
			t.Fatalf("ImportRawKeyAsBase64 failed: %v", err)
		}

		plaintext := []byte("Test data for compatibility check")

		// Both keys should be able to encrypt/decrypt
		// Test Tink key
		ciphertext1, err := service.EncryptFile(tinkKey, plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt with Tink key: %v", err)
		}

		decrypted1, err := service.DecryptFile(tinkKey, ciphertext1)
		if err != nil {
			t.Fatalf("Failed to decrypt with Tink key: %v", err)
		}

		if string(decrypted1) != string(plaintext) {
			t.Fatal("Tink key encryption/decryption failed")
		}

		// Test converted raw key
		ciphertext2, err := service.EncryptFile(rawKeyConverted, plaintext)
		if err != nil {
			t.Fatalf("Failed to encrypt with converted raw key: %v", err)
		}

		decrypted2, err := service.DecryptFile(rawKeyConverted, ciphertext2)
		if err != nil {
			t.Fatalf("Failed to decrypt with converted raw key: %v", err)
		}

		if string(decrypted2) != string(plaintext) {
			t.Fatal("Converted raw key encryption/decryption failed")
		}

		t.Log("Both native Tink key and converted raw key work correctly")
	})
}

func TestStringEncryptionWithRawKey(t *testing.T) {
	service := services.NewCryptographicService()

	t.Run("Encrypt and decrypt string with raw key", func(t *testing.T) {
		// Generate and convert a raw key
		rawKey := make([]byte, 32)
		rand.Read(rawKey)
		keyBase64, err := service.ImportRawKeyAsBase64(rawKey)
		if err != nil {
			t.Fatalf("ImportRawKeyAsBase64 failed: %v", err)
		}

		plaintext := "This is a secret message"

		// Encrypt
		ciphertext, err := service.EncryptString(keyBase64, plaintext)
		if err != nil {
			t.Fatalf("EncryptString failed: %v", err)
		}

		if ciphertext == "" {
			t.Fatal("Expected non-empty ciphertext")
		}

		// Decrypt
		decrypted, err := service.DecryptString(keyBase64, ciphertext)
		if err != nil {
			t.Fatalf("DecryptString failed: %v", err)
		}

		if decrypted != plaintext {
			t.Fatalf("Decrypted string mismatch.\nExpected: %s\nGot: %s", plaintext, decrypted)
		}

		t.Log("Successfully encrypted and decrypted string with raw key")
	})
}

func TestHexKeyFromKMS(t *testing.T) {
	hexKey := "4939BD67B68947A16EC5F90036C7379924C10114CFBBAC123235101D90A4B004"
	keyBytes, err := helper.HexToBytes(hexKey)
	if err != nil {
		t.Fatalf("Failed to convert hex to bytes: %v", err)
	}
	fmt.Printf("Decoded Key : %s\n", keyBytes)

	service := services.NewCryptographicService()
	rawkey, err := service.ImportRawKeyAsBase64(keyBytes)
	assert.NoError(t, err, "ImportRawKeyAsBase64 should not return an error")
	assert.NotEmpty(t, rawkey, "Imported raw key should not be empty")

	encrypted, err := service.EncryptFile(rawkey, []byte("Test Data"))
	assert.NoError(t, err, "EncryptFile should not return an error")
	assert.NotEmpty(t, encrypted, "Encrypted data should not be empty")

	decrypted, err := service.DecryptFile(rawkey, encrypted)
	assert.NoError(t, err, "DecryptFile should not return an error")
	assert.Equal(t, []byte("Test Data"), decrypted, "Decrypted data should match original")
}

func TestKMSKeyWrappingFlow(t *testing.T) {
	t.Run("Full KMS flow: KEK wraps DEK", func(t *testing.T) {
		service := services.NewCryptographicService()

		// Simulate KEK from KMS (like in bootstrapApp.go)
		kekHex := "4939BD67B68947A16EC5F90036C7379924C10114CFBBAC123235101D90A4B004"
		kekBytes, err := helper.HexToBytes(kekHex)
		assert.NoError(t, err, "HexToBytes should not fail")

		// Convert KEK to Tink keyset (this is what bootstrapApp.go does now)
		kek, err := service.ImportRawKeyAsBase64(kekBytes)
		assert.NoError(t, err, "ImportRawKeyAsBase64 for KEK should not fail")
		assert.NotEmpty(t, kek, "KEK should not be empty")

		// Simulate DEK from KMS (like in getEncryptionKey)
		dekHex := "A1B2C3D4E5F6071829304A5B6C7D8E9F0A1B2C3D4E5F6071829304A5B6C7D8E9"
		dekBytes, err := helper.HexToBytes(dekHex)
		assert.NoError(t, err, "HexToBytes for DEK should not fail")

		// Convert DEK to Tink keyset
		dek, err := service.ImportRawKeyAsBase64(dekBytes)
		assert.NoError(t, err, "ImportRawKeyAsBase64 for DEK should not fail")
		assert.NotEmpty(t, dek, "DEK should not be empty")

		// Wrap DEK with KEK (encrypt the DEK key string with KEK)
		wrappedDEK, err := service.EncryptString(kek, dek)
		assert.NoError(t, err, "EncryptString (wrapping DEK with KEK) should not fail")
		assert.NotEmpty(t, wrappedDEK, "Wrapped DEK should not be empty")

		// Unwrap DEK (decrypt to get back the DEK)
		unwrappedDEK, err := service.DecryptString(kek, wrappedDEK)
		assert.NoError(t, err, "DecryptString (unwrapping DEK) should not fail")
		assert.Equal(t, dek, unwrappedDEK, "Unwrapped DEK should match original DEK")

		t.Logf("Original DEK length: %d", len(dek))
		t.Logf("Unwrapped DEK length: %d", len(unwrappedDEK))
		t.Logf("Are they equal? %v", dek == unwrappedDEK)

		// Use the unwrapped DEK to encrypt file data
		plaintext := []byte("Sensitive file data that needs encryption")
		encrypted, err := service.EncryptFile(unwrappedDEK, plaintext)
		assert.NoError(t, err, "EncryptFile should not fail")
		assert.NotEmpty(t, encrypted, "Encrypted data should not be empty")

		// Decrypt the file data
		decrypted, err := service.DecryptFile(unwrappedDEK, encrypted)
		assert.NoError(t, err, "DecryptFile should not fail")
		assert.Equal(t, plaintext, decrypted, "Decrypted data should match original")

		t.Log("✅ Full KMS key wrapping flow works correctly")
		t.Logf("   KEK (base64): %d chars", len(kek))
		t.Logf("   DEK (base64): %d chars", len(dek))
		t.Logf("   Wrapped DEK (base64): %d chars", len(wrappedDEK))
		t.Logf("   Unwrapped DEK (base64): %d chars", len(unwrappedDEK))
		t.Logf("   Encrypted file: %d bytes", len(encrypted))
	})

	t.Run("KEK from file wraps KMS DEK", func(t *testing.T) {
		service := services.NewCryptographicService()

		// Simulate KEK from local file (native Tink keyset)
		kek, err := service.GenerateKey()
		assert.NoError(t, err, "GenerateKey for KEK should not fail")

		// Simulate DEK from KMS (raw key converted to Tink)
		dekHex := "B2C3D4E5F6071829304A5B6C7D8E9F0A1B2C3D4E5F6071829304A5B6C7D8E9F0"
		dekBytes, err := helper.HexToBytes(dekHex)
		assert.NoError(t, err, "HexToBytes should not fail")

		dek, err := service.ImportRawKeyAsBase64(dekBytes)
		assert.NoError(t, err, "ImportRawKeyAsBase64 should not fail")

		// Wrap and unwrap
		wrappedDEK, err := service.EncryptString(kek, dek)
		assert.NoError(t, err, "EncryptString should not fail")

		unwrappedDEK, err := service.DecryptString(kek, wrappedDEK)
		assert.NoError(t, err, "DecryptString should not fail")
		assert.Equal(t, dek, unwrappedDEK, "Unwrapped DEK should match")

		t.Log("✅ Native KEK can wrap KMS-sourced DEK correctly")
	})
}

func TestKEKValidation(t *testing.T) {
	service := services.NewCryptographicService()

	t.Run("Validate KEK is proper Tink keyset", func(t *testing.T) {
		// Simulate loading KEK from KMS (like your actual scenario)
		kekHex := "4939BD67B68947A16EC5F90036C7379924C10114CFBBAC123235101D90A4B004"
		kekBytes, err := helper.HexToBytes(kekHex)
		assert.NoError(t, err, "HexToBytes should not fail")
		assert.Equal(t, 32, len(kekBytes), "KEK should be 32 bytes")

		// Convert to Tink keyset
		kek, err := service.ImportRawKeyAsBase64(kekBytes)
		assert.NoError(t, err, "ImportRawKeyAsBase64 should not fail")
		assert.NotEmpty(t, kek, "KEK should not be empty")

		t.Logf("KEK base64 length: %d", len(kek))
		t.Logf("KEK first 50 chars: %.50s", kek)

		// Try to use the KEK to encrypt something (like wrapping a DEK)
		testPlaintext := "This is a test DEK that needs to be wrapped"
		wrapped, err := service.EncryptString(kek, testPlaintext)
		assert.NoError(t, err, "EncryptString with KEK should work")
		assert.NotEmpty(t, wrapped, "Wrapped text should not be empty")

		// Try to unwrap
		unwrapped, err := service.DecryptString(kek, wrapped)
		assert.NoError(t, err, "DecryptString with KEK should work")
		assert.Equal(t, testPlaintext, unwrapped, "Unwrapped should match original")

		t.Log("✅ KEK is valid and can be used for wrapping")
	})

	t.Run("Test with potentially invalid KEK - raw hex", func(t *testing.T) {
		// This simulates if KEK was NOT converted properly
		invalidKEK := "4939BD67B68947A16EC5F90036C7379924C10114CFBBAC123235101D90A4B004"

		// Try to use it directly (this should fail)
		_, err := service.EncryptString(invalidKEK, "test")
		assert.Error(t, err, "EncryptString with raw hex should fail")
		t.Logf("Expected error: %v", err)
	})

	t.Run("Test with potentially invalid KEK - raw bytes base64", func(t *testing.T) {
		// This simulates if KEK was base64-encoded raw bytes (not Tink keyset)
		kekBytes := make([]byte, 32)
		rand.Read(kekBytes)
		invalidKEK := base64.StdEncoding.EncodeToString(kekBytes)

		// Try to use it directly (this should fail with "invalid keyset")
		_, err := service.EncryptString(invalidKEK, "test")
		assert.Error(t, err, "EncryptString with base64 raw bytes should fail")
		assert.Contains(t, err.Error(), "invalid keyset", "Should return 'invalid keyset' error")
		t.Logf("Expected error: %v", err)
	})
}
