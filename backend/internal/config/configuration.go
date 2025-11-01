package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Properties struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	MKeyPath          string
	StorageEndpoint   string
	StrorageAccessID  string
	StrorageSecretKey string
	StorageSSL        bool
	BucketName        string

	HashMethod        string
	EncMethod         string
	HashEncryptedFile bool

	HydraPublicURL string
	HydraAdminURL  string

	KMSEnable bool
	KMSKeyUID string
	KMSUrl    string
	KeyPath   string
	CertPath  string
	CAPath    string

	// OpenTelemetry
	OTELEnable     bool
	OTELEndpoint   string
	ServiceName    string
	ServiceVersion string
	Environment    string
}

func LoadProperties() *Properties {
	log.Println("Load configuration from .env using gotodotenv")
	var err error

	if os.Getenv("DEVELOPER_HOST") == "true" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	properties := &Properties{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		DBSSLMode:         getEnvWithDefault("DB_SSLMODE", "disable"),
		MKeyPath:          os.Getenv("MKEY_PATH"),
		StorageEndpoint:   os.Getenv("STORAGE_ENDPOINT"),
		StrorageAccessID:  os.Getenv("STORAGE_ACCESS_KEY"),
		StrorageSecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		StorageSSL:        os.Getenv("STRORAGE_SSL") == "true",
		BucketName:        os.Getenv("BUCKET_NAME"),
		HashMethod:        os.Getenv("HASH_METHOD"),
		HydraPublicURL:    os.Getenv("HYDRA_PUBLIC_URL"),
		HydraAdminURL:     os.Getenv("HYDRA_ADMIN_URL"),
		KMSEnable:         os.Getenv("KMS_ENABLE") == "true",
		KMSKeyUID:         os.Getenv("KMS_KEY_UID"),
		KMSUrl:            os.Getenv("KMS_URL"),
		KeyPath:           os.Getenv("KEY_PATH"),
		CertPath:          os.Getenv("CERT_PATH"),
		CAPath:            os.Getenv("CA_PATH"),
		EncMethod:         os.Getenv("ENC_METHOD"),
		HashEncryptedFile: os.Getenv("HASH_ENCRYPTED_FILE") == "true",
		OTELEnable:        getEnvWithDefault("OTEL_ENABLE", "false") == "true",
		OTELEndpoint:      getEnvWithDefault("OTEL_ENDPOINT", "localhost:4318"),
		ServiceName:       getEnvWithDefault("SERVICE_NAME", "crypsis-backend"),
		ServiceVersion:    getEnvWithDefault("SERVICE_VERSION", "1.0.0"),
		Environment:       getEnvWithDefault("ENVIRONMENT", "development"),
	}

	return properties
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
