package model

type KmsResponse struct {
	Tag   string          `json:"tag"`
	Type  string          `json:"type"`
	Value []ValueResponse `json:"value"` // Use slice to handle the array
}

// ValueResponse holds key-related details
type ValueResponse struct {
	Tag   string      `json:"tag"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"` // Make this flexible (string, int, or list)
}

// ResponseItem catches each item in the array
type ResponseItem struct {
	Tag   string      `json:"tag"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// KeyAttributes catches key attributes
type KeyAttributes struct {
	CryptographicAlgorithm string `json:"cryptographic_algorithm"`
	CryptographicLength    int    `json:"cryptographic_length"`
	CryptographicUsageMask int    `json:"cryptographic_usage_mask"`
	KeyFormatType          string `json:"key_format_type"`
	ObjectType             string `json:"object_type"`
	Sensitive              bool   `json:"sensitive"`
}

// KeyMaterial catches key material
type KeyMaterial struct {
	ByteString string `json:"byte_string"`
}

// KeyBlock catches key block information
type KeyBlock struct {
	KeyFormatType          string        `json:"key_format_type"`
	KeyMaterial            KeyMaterial   `json:"key_material"`
	Attributes             KeyAttributes `json:"attributes"`
	CryptographicAlgorithm string        `json:"cryptographic_algorithm"`
	CryptographicLength    int           `json:"cryptographic_length"`
}

// LocateResponse represents the response for locating keys
type LocateResponse struct {
	LocatedItems      int      `json:"located_items"`
	UniqueIdentifiers []string `json:"unique_identifiers"`
}
