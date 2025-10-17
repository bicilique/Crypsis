package main

import (
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/services"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("KEK Converter - Convert raw hex AES-256 key to Tink keyset format")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  go run scripts/convert_kek.go <hex_key> [output_file]")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run scripts/convert_kek.go 4939BD67B68947A16EC5F90036C7379924C10114CFBBAC123235101D90A4B004")
		fmt.Println("  go run scripts/convert_kek.go 4939BD67B68947A16EC5F90036C7379924C10114CFBBAC123235101D90A4B004 resources/sample.key")
		fmt.Println("")
		os.Exit(1)
	}

	hexKey := os.Args[1]
	outputFile := "resources/sample.key"
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	}

	fmt.Printf("Converting KEK from hex to Tink keyset...\n")
	fmt.Printf("Input hex: %.60s...\n", hexKey)
	fmt.Printf("Output file: %s\n\n", outputFile)

	// Convert hex to bytes
	keyBytes, err := helper.HexToBytes(hexKey)
	if err != nil {
		fmt.Printf("❌ Error: Invalid hex key: %v\n", err)
		os.Exit(1)
	}

	if len(keyBytes) != 32 {
		fmt.Printf("❌ Error: Invalid key length: expected 32 bytes for AES-256, got %d bytes\n", len(keyBytes))
		fmt.Printf("   Hex string should be 64 characters (2 hex chars per byte)\n")
		os.Exit(1)
	}

	fmt.Printf("✓ Hex key decoded: %d bytes\n", len(keyBytes))

	// Convert to Tink keyset
	service := services.NewCryptographicService()
	kek, err := service.ImportRawKeyAsBase64(keyBytes)
	if err != nil {
		fmt.Printf("❌ Error: Failed to convert to Tink keyset: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Converted to Tink keyset: %d chars (base64)\n", len(kek))

	// Decode base64 to binary
	keysetBytes, err := base64.StdEncoding.DecodeString(kek)
	if err != nil {
		fmt.Printf("❌ Error: Failed to decode base64: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Decoded to binary: %d bytes\n", len(keysetBytes))

	// Write to file
	err = os.WriteFile(outputFile, keysetBytes, 0600)
	if err != nil {
		fmt.Printf("❌ Error: Failed to write to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Saved to file: %s\n\n", outputFile)

	// Verify the conversion
	fmt.Println("Verifying the conversion...")
	loadedKek, err := helper.FileToBase64(outputFile)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to read back the file: %v\n", err)
	} else if loadedKek != kek {
		fmt.Printf("⚠️  Warning: Loaded KEK doesn't match original\n")
	} else {
		fmt.Printf("✓ Verification successful\n")
	}

	// Test encryption/decryption
	fmt.Println("\nTesting KEK...")
	testText := "Test encryption with KEK"
	encrypted, err := service.EncryptString(kek, testText)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to encrypt test text: %v\n", err)
	} else {
		decrypted, err := service.DecryptString(kek, encrypted)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to decrypt test text: %v\n", err)
		} else if decrypted != testText {
			fmt.Printf("⚠️  Warning: Decrypted text doesn't match original\n")
		} else {
			fmt.Printf("✓ KEK encryption/decryption test passed\n")
		}
	}

	fmt.Println("\n✅ KEK conversion completed successfully!")
	fmt.Println("   You can now use this KEK in your application.")
	fmt.Println("")
}
