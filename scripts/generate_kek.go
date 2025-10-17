package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/insecurecleartextkeyset"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func main() {
	// Generate a new Tink AES-256-GCM keyset
	handle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate key: %v\n", err)
		os.Exit(1)
	}

	// Serialize the keyset to bytes
	buf := new(bytes.Buffer)
	writer := keyset.NewBinaryWriter(buf)
	if err := insecurecleartextkeyset.Write(handle, writer); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to serialize key: %v\n", err)
		os.Exit(1)
	}

	// Write raw bytes to file
	keyBytes := buf.Bytes()
	if err := os.WriteFile("../resources/sample.key", keyBytes, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write key file: %v\n", err)
		os.Exit(1)
	}

	// Also print base64 for reference
	keyBase64 := base64.StdEncoding.EncodeToString(keyBytes)
	fmt.Printf("✓ New Tink keyset generated and saved to resources/sample.key\n")
	fmt.Printf("✓ Key length: %d bytes\n", len(keyBytes))
	fmt.Printf("✓ Base64 encoded key:\n%s\n", keyBase64)
}
