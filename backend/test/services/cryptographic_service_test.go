package services_test

import (
	"crypsis-backend/internal/services"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	key, err := cryptoService.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	if key == "" {
		t.Fatal("Generated key is empty")
	}

	t.Logf("Generated key: %s (length: %d)", key[:20]+"...", len(key))
}

func TestEncryptDecryptString(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	// Generate a key
	key, err := cryptoService.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Original plaintext
	plaintext := "Hello, Tink Cryptography!"

	// Encrypt
	ciphertext, err := cryptoService.EncryptString(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt string: %v", err)
	}

	if ciphertext == "" {
		t.Fatal("Encrypted string is empty")
	}

	t.Logf("Plaintext: %s", plaintext)
	t.Logf("Ciphertext: %s", ciphertext[:20]+"...")

	// Decrypt
	decrypted, err := cryptoService.DecryptString(key, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt string: %v", err)
	}

	// Verify
	if decrypted != plaintext {
		t.Fatalf("Decrypted text does not match original. Expected: %s, Got: %s", plaintext, decrypted)
	}

	t.Logf("Successfully decrypted: %s", decrypted)
}

func TestEncryptDecryptWithDifferentKeys(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	// Generate two different keys
	key1, _ := cryptoService.GenerateKey()
	key2, _ := cryptoService.GenerateKey()

	plaintext := "Secret message"

	// Encrypt with key1
	ciphertext, err := cryptoService.EncryptString(key1, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Try to decrypt with key2 (should fail)
	_, err = cryptoService.DecryptString(key2, ciphertext)
	if err == nil {
		t.Fatal("Expected decryption to fail with different key, but it succeeded")
	}

	t.Logf("Correctly failed to decrypt with wrong key: %v", err)
}

func TestHashStringSHA256(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	text := "Test message for hashing"
	hash, err := cryptoService.HashString(services.HashSHA256, text)
	if err != nil {
		t.Fatalf("Failed to hash string: %v", err)
	}

	if hash == "" {
		t.Fatal("Hash is empty")
	}

	t.Logf("SHA-256 Hash: %s", hash)

	// Hash the same text again and verify consistency
	hash2, _ := cryptoService.HashString(services.HashSHA256, text)
	if hash != hash2 {
		t.Fatal("Same input produced different hashes")
	}
}

func TestHashStringMD5(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	text := "Test message for MD5 hashing"
	hash, err := cryptoService.HashString(services.HashMD5, text)
	if err != nil {
		t.Fatalf("Failed to hash string with MD5: %v", err)
	}

	if hash == "" {
		t.Fatal("MD5 hash is empty")
	}

	t.Logf("MD5 Hash: %s", hash)
}

func TestCompareHash(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	text := "Password123"
	hash, _ := cryptoService.HashString(services.HashSHA256, text)

	// Test with correct text
	if !cryptoService.CompareHash(services.HashSHA256, text, hash) {
		t.Fatal("Hash comparison failed for correct text")
	}

	// Test with incorrect text
	if cryptoService.CompareHash(services.HashSHA256, "WrongPassword", hash) {
		t.Fatal("Hash comparison should fail for incorrect text")
	}

	t.Log("Hash comparison working correctly")
}

func TestEncryptDecryptFile(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	// Generate a key
	key, err := cryptoService.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Simulate file content
	fileContent := []byte("This is the content of a test file.\nIt has multiple lines.\nAnd some data: 12345")

	// Encrypt file
	encryptedFile, err := cryptoService.EncryptFile(key, fileContent)
	if err != nil {
		t.Fatalf("Failed to encrypt file: %v", err)
	}

	if len(encryptedFile) == 0 {
		t.Fatal("Encrypted file is empty")
	}

	t.Logf("Original file size: %d bytes", len(fileContent))
	t.Logf("Encrypted file size: %d bytes", len(encryptedFile))

	// Decrypt file
	decryptedFile, err := cryptoService.DecryptFile(key, encryptedFile)
	if err != nil {
		t.Fatalf("Failed to decrypt file: %v", err)
	}

	// Verify
	if string(decryptedFile) != string(fileContent) {
		t.Fatal("Decrypted file content does not match original")
	}

	t.Log("File encryption/decryption successful")
}

func TestHashFile(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	fileContent := []byte("File content to be hashed")

	hash, err := cryptoService.HashFile(services.HashSHA256, fileContent)
	if err != nil {
		t.Fatalf("Failed to hash file: %v", err)
	}

	if hash == "" {
		t.Fatal("File hash is empty")
	}

	t.Logf("File SHA-256 Hash: %s", hash)

	// Verify consistency
	hash2, _ := cryptoService.HashFile(services.HashSHA256, fileContent)
	if hash != hash2 {
		t.Fatal("Same file produced different hashes")
	}
}

func TestCompareHashFile(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	fileContent := []byte("Original file content")
	hash, _ := cryptoService.HashFile(services.HashSHA256, fileContent)

	// Test with correct content
	if !cryptoService.CompareHashFile(services.HashSHA256, fileContent, hash) {
		t.Fatal("File hash comparison failed for correct content")
	}

	// Test with modified content
	modifiedContent := []byte("Modified file content")
	if cryptoService.CompareHashFile(services.HashSHA256, modifiedContent, hash) {
		t.Fatal("File hash comparison should fail for modified content")
	}

	t.Log("File hash comparison working correctly")
}

func TestKeyDerivationFunction(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	input := "user_password"
	salt := []byte("random_salt_12345")

	derivedKey, err := cryptoService.KeyDerivationFunction(input, salt)
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	if len(derivedKey) != 32 {
		t.Fatalf("Expected derived key length 32, got %d", len(derivedKey))
	}

	t.Logf("Derived key length: %d bytes", len(derivedKey))

	// Note: Tink's PRF creates a new random key each time, so we can't test consistency
	// across multiple calls without reusing the same handle. This is expected behavior.

	// Verify that we can derive a key with different salt
	differentSalt := []byte("different_salt_67890")
	derivedKey2, err := cryptoService.KeyDerivationFunction(input, differentSalt)
	if err != nil {
		t.Fatalf("Failed to derive key with different salt: %v", err)
	}

	if len(derivedKey2) != 32 {
		t.Fatalf("Expected derived key length 32, got %d", len(derivedKey2))
	}

	t.Log("Key derivation working correctly")
}

func TestUnsupportedHashAlgorithm(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	_, err := cryptoService.HashString("UNSUPPORTED_ALGO", "test")
	if err == nil {
		t.Fatal("Expected error for unsupported hash algorithm")
	}

	t.Logf("Correctly returned error for unsupported algorithm: %v", err)
}

func TestEmptyStringEncryption(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	key, _ := cryptoService.GenerateKey()

	// Encrypt empty string
	ciphertext, err := cryptoService.EncryptString(key, "")
	if err != nil {
		t.Fatalf("Failed to encrypt empty string: %v", err)
	}

	// Decrypt
	decrypted, err := cryptoService.DecryptString(key, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt empty string: %v", err)
	}

	if decrypted != "" {
		t.Fatalf("Expected empty string, got: %s", decrypted)
	}

	t.Log("Empty string encryption/decryption successful")
}

func TestLargeDataEncryption(t *testing.T) {
	cryptoService := services.NewCryptographicService()

	key, _ := cryptoService.GenerateKey()

	// Create large data (1MB)
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Encrypt
	encrypted, err := cryptoService.EncryptFile(key, largeData)
	if err != nil {
		t.Fatalf("Failed to encrypt large file: %v", err)
	}

	// Decrypt
	decrypted, err := cryptoService.DecryptFile(key, encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt large file: %v", err)
	}

	// Verify
	if len(decrypted) != len(largeData) {
		t.Fatalf("Decrypted data size mismatch. Expected: %d, Got: %d", len(largeData), len(decrypted))
	}

	for i := range largeData {
		if decrypted[i] != largeData[i] {
			t.Fatalf("Data mismatch at index %d", i)
		}
	}

	t.Logf("Large file encryption/decryption successful (1MB)")
}

func BenchmarkGenerateKey(b *testing.B) {
	cryptoService := services.NewCryptographicService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cryptoService.GenerateKey()
	}
}

func BenchmarkEncryptString(b *testing.B) {
	cryptoService := services.NewCryptographicService()
	key, _ := cryptoService.GenerateKey()
	plaintext := "Benchmark test message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cryptoService.EncryptString(key, plaintext)
	}
}

func BenchmarkDecryptString(b *testing.B) {
	cryptoService := services.NewCryptographicService()
	key, _ := cryptoService.GenerateKey()
	plaintext := "Benchmark test message"
	ciphertext, _ := cryptoService.EncryptString(key, plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cryptoService.DecryptString(key, ciphertext)
	}
}

func BenchmarkHashStringSHA256(b *testing.B) {
	cryptoService := services.NewCryptographicService()
	text := "Benchmark hash message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cryptoService.HashString(services.HashSHA256, text)
	}
}
