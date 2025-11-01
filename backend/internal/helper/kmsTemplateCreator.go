package helper

import (
	"encoding/hex"
	"encoding/json"
)

// Attribute represents a key attribute (both Cryptographic and Vendor attributes)
type Attribute struct {
	Tag   string      `json:"tag"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// BodyRequest represents the JSON structure for key operations
type BodyRequest struct {
	Tag   string        `json:"tag"`
	Type  string        `json:"type"`
	Value []interface{} `json:"value"`
}

// GenerateKeyTemplate creates a JSON request to generate a symmetric key
func GenerateKeyTemplate(keyName string) (string, error) {
	keyNames := []string{keyName}
	jsonArray, _ := json.Marshal(keyNames)
	keyNameHex := hex.EncodeToString(jsonArray)

	createKeyTemplate := BodyRequest{
		Tag:  "Create",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "ObjectType", Type: "Enumeration", Value: "SymmetricKey"},
			Attribute{
				Tag:  "Attributes",
				Type: "Structure",
				Value: []Attribute{
					{Tag: "CryptographicAlgorithm", Type: "Enumeration", Value: "AES"},
					{Tag: "CryptographicLength", Type: "Integer", Value: 256},
					{Tag: "CryptographicUsageMask", Type: "Integer", Value: 2108},
					{Tag: "KeyFormatType", Type: "Enumeration", Value: "TransparentSymmetricKey"},
					{Tag: "ObjectType", Type: "Enumeration", Value: "SymmetricKey"},
					{
						Tag:  "VendorAttributes",
						Type: "Structure",
						Value: []Attribute{
							{
								Tag:  "VendorAttributes",
								Type: "Structure",
								Value: []Attribute{
									{Tag: "VendorIdentification", Type: "TextString", Value: "cosmian"},
									{Tag: "AttributeName", Type: "TextString", Value: "tag"},
									{Tag: "AttributeValue", Type: "ByteString", Value: keyNameHex}, // 🔥 Sesuai dengan format curl
								},
							},
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(createKeyTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateLocateKeyTemplate creates a JSON request to locate a symmetric key
func GenerateLocateKeyTemplate(keyName string) (string, error) {
	keyNames := []string{keyName}
	jsonArray, _ := json.Marshal(keyNames)
	keyNameHex := hex.EncodeToString(jsonArray)

	locateKeyTemplate := BodyRequest{
		Tag:  "Locate",
		Type: "Structure",
		Value: []interface{}{
			Attribute{
				Tag:  "Attributes",
				Type: "Structure",
				Value: []Attribute{
					{
						Tag:  "VendorAttributes",
						Type: "Structure",
						Value: []Attribute{
							{
								Tag:  "VendorAttributes",
								Type: "Structure",
								Value: []Attribute{
									{Tag: "VendorIdentification", Type: "TextString", Value: "cosmian"},
									{Tag: "AttributeName", Type: "TextString", Value: "tag"},
									{Tag: "AttributeValue", Type: "ByteString", Value: keyNameHex},
								},
							},
						},
					},
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(locateKeyTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateEncryptTemplate creates a JSON request for encryption
func GenerateEncryptTemplate(keyUID string, plaintext string) (string, error) {
	// Convert plaintext to hex representation
	encryptTemplate := BodyRequest{
		Tag:  "Encrypt",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
			Attribute{Tag: "Data", Type: "ByteString", Value: plaintext},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(encryptTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateCoverCryptEncryptTemplate creates a JSON request for CoverCrypt encryption
func GenerateCoverCryptEncryptTemplate(keyUID, plaintext string) (string, error) {
	encryptTemplate := BodyRequest{
		Tag:  "Encrypt",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
			Attribute{
				Tag:  "CryptographicParameters",
				Type: "Structure",
				Value: []interface{}{
					Attribute{Tag: "CryptographicAlgorithm", Type: "Enumeration", Value: "CoverCrypt"},
				},
			},
			Attribute{Tag: "Data", Type: "ByteString", Value: plaintext},
		},
	}

	jsonData, err := json.MarshalIndent(encryptTemplate, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// GenerateDecryptTemplate creates a JSON request for decryption
func GenerateDecryptTemplate(keyUID, encryptedData, ivCounterNonce, authTag string) (string, error) {
	decryptTemplate := BodyRequest{
		Tag:  "Decrypt",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
			Attribute{Tag: "Data", Type: "ByteString", Value: encryptedData},
			Attribute{Tag: "IvCounterNonce", Type: "ByteString", Value: ivCounterNonce},
			Attribute{Tag: "AuthenticatedEncryptionTag", Type: "ByteString", Value: authTag},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(decryptTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateExportTemplate creates a JSON request for key export
func GenerateExportTemplate(keyUID string) (string, error) {
	exportTemplate := BodyRequest{
		Tag:  "Export",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
			Attribute{Tag: "KeyWrapType", Type: "Enumeration", Value: "AsRegistered"},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(exportTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateDestroyTemplate creates a JSON request for key deletion
func GenerateDestroyTemplate(keyUID string) (string, error) {
	exportTemplate := BodyRequest{
		Tag:  "Destroy",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
			Attribute{Tag: "Remove", Type: "Boolean", Value: true},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(exportTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateRevokeTemplate creates a JSON request for key revocation
func GenerateRevokeTemplate(keyUID string) (string, error) {
	exportTemplate := BodyRequest{
		Tag:  "Revoke",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
			Attribute{Tag: "RevocationReason", Type: "TextString", Value: "key was compromised"},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(exportTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateReKeyTemplate creates a JSON request for key rekey
func GenerateReKeyTemplate(keyUID string) (string, error) {
	exportTemplate := BodyRequest{
		Tag:  "ReKey",
		Type: "Structure",
		Value: []interface{}{
			Attribute{Tag: "UniqueIdentifier", Type: "TextString", Value: keyUID},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(exportTemplate)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// GenerateReKeyTemplate creates a JSON request for generating a key pair with ECDH
func GenerateKeyPairTemplate(keyName string) (string, error) {
	keyNames := []string{keyName}
	jsonArray, _ := json.Marshal(keyNames)
	keyNameHex := hex.EncodeToString(jsonArray)

	keyPairRequest := BodyRequest{
		Tag:  "CreateKeyPair",
		Type: "Structure",
		Value: []interface{}{
			Attribute{
				Tag:  "CommonAttributes",
				Type: "Structure",
				Value: []interface{}{
					Attribute{Tag: "CryptographicAlgorithm", Type: "Enumeration", Value: "ECDH"},
					Attribute{Tag: "CryptographicLength", Type: "Integer", Value: 253},
					Attribute{
						Tag:  "CryptographicDomainParameters",
						Type: "Structure",
						Value: []interface{}{
							Attribute{Tag: "QLength", Type: "Integer", Value: 253},
							Attribute{Tag: "RecommendedCurve", Type: "Enumeration", Value: "CURVE25519"},
						},
					},
					Attribute{Tag: "CryptographicUsageMask", Type: "Integer", Value: 2108},
					Attribute{Tag: "KeyFormatType", Type: "Enumeration", Value: "ECPrivateKey"},
					Attribute{Tag: "ObjectType", Type: "Enumeration", Value: "PrivateKey"},
					Attribute{
						Tag:  "VendorAttributes",
						Type: "Structure",
						Value: []interface{}{
							Attribute{
								Tag:  "VendorAttributes",
								Type: "Structure",
								Value: []interface{}{
									Attribute{Tag: "VendorIdentification", Type: "TextString", Value: "cosmian"},
									Attribute{Tag: "AttributeName", Type: "TextString", Value: "tag"},
									Attribute{Tag: "AttributeValue", Type: "ByteString", Value: keyNameHex},
								},
							},
						},
					},
				},
			},
		},
	}

	// Convert struct to JSON
	jsonData, err := json.MarshalIndent(keyPairRequest, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
