package config

import (
	"context"
	delivery "crypsis-backend/internal/delivery/http"
	"crypsis-backend/internal/delivery/middlewere"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/repository"
	"crypsis-backend/internal/services"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppConfig struct {
	Properties *Properties
	DB         *gorm.DB
}

func BootstrapApp(config *AppConfig) {
	// initialize repositories
	repos := initRepositories(config.DB)

	// initialize services
	services := initServices(config.Properties, repos, config.DB)

	// initialize http server
	httpServer := initHttpServer(services, config.Properties, repos.adminRepository)

	// Start HTTP server with graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startHTTPServer(ctx, httpServer)

}

// startHTTPServer starts the HTTP server with graceful shutdown
func startHTTPServer(ctx context.Context, server *http.Server) {
	go func() {
		log.Printf("🌐 Starting HTTP server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("🛑 Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("✅ HTTP server stopped gracefully")
	}
}

func initHttpServer(services Services, config *Properties, adminRepo repository.AdminRepository) *http.Server {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	tokenMiddlewereConfig := middlewere.TokenMiddlewareConfig{
		HydraAdminURL: config.HydraAdminURL + "/admin/oauth2/introspect",
		RequiredScope: "offline",
		AdminRepo:     adminRepo,
	}

	router := delivery.RouterConfig{
		Router:          gin.Default(),
		ClientHandler:   delivery.NewClientHandler(services.fileService),
		AdminHandler:    delivery.NewAdminHandler(services.applicationService, services.adminService, services.fileService),
		HydraAdminURL:   config.HydraAdminURL,
		TokenMiddlewere: tokenMiddlewereConfig,
	}
	router.Setup()

	return &http.Server{
		Addr:    ":" + port,
		Handler: router.Router,
	}
}

func initServices(config *Properties, repos Repositories, db *gorm.DB) Services {
	// Load admin IDs into cache
	repos.adminRepository.LoadAdminIDs(context.Background())

	minIOService := services.NewMinioService(model.MinIOConfig{
		Endpoint:        config.StorageEndpoint,
		AccessKeyID:     config.StrorageAccessID,
		SecretAccessKey: config.StrorageSecretKey,
		BucketName:      config.BucketName,
		UseSSL:          config.StorageSSL,
	})

	cryptographicService := services.NewCryptographicService()

	keyConfig := &model.KeyConfig{
		KMSEnable: config.KMSEnable,
	}

	var kmsService services.KMSInterface
	if config.KMSEnable {
		// Load key from KMS
		secureClient := helper.CreateHTTPSClient(config.CertPath, config.KeyPath, config.CAPath)
		kmsService = services.NewKmsService(secureClient, config.KMSUrl)
		keyHex, err := kmsService.ExportKey(context.Background(), config.KMSKeyUID)
		if err != nil {
			log.Fatalf("Failed to export key: %v", err)
		}
		slog.Info("Successfully exported KEK from KMS", slog.String("keyUID", config.KMSKeyUID), slog.Int("hex_length", len(keyHex)))

		// Convert hex to bytes
		keyBytes, err := helper.HexToBytes(keyHex)
		if err != nil {
			log.Fatalf("Failed to decode hex key: %v", err)
		}
		slog.Info("Converted KEK hex to bytes", slog.Int("bytes_length", len(keyBytes)))

		// Convert raw key bytes to Tink keyset format
		key, err := cryptographicService.ImportRawKeyAsBase64(keyBytes)
		if err != nil {
			log.Fatalf("Failed to convert raw key to Tink keyset: %v", err)
		}
		slog.Info("Successfully converted KEK to Tink keyset", slog.Int("base64_length", len(key)))

		keyConfig.UID = config.KMSKeyUID
		keyConfig.KEK = key
	} else {
		// Load key from file
		key, err := helper.FileToBase64(config.MKeyPath)
		if err != nil {
			log.Fatalf("Failed to decode key: %v", err)
		}
		keyConfig.KEK = key
	}

	oauth2Service := services.NewHydraService(config.HydraAdminURL, config.HydraPublicURL)
	adminService := services.NewAdminService(oauth2Service, repos.adminRepository, repos.fileLogRepository, cryptographicService)
	applicationService := services.NewApplicationService(oauth2Service, repos.applicationRepository, repos.fileLogRepository)

	fileServiceParams := services.FileServiceParams{
		CryptoService:         cryptographicService,
		StorageService:        minIOService,
		KMSService:            kmsService,
		FileRepository:        repos.fileRepository,
		FileLogsRepository:    repos.fileLogRepository,
		ApplicationRepository: repos.applicationRepository,
		AdminRepository:       repos.adminRepository,
		DB:                    db,
		KeyConfig:             keyConfig,
		BucketName:            config.BucketName,
		HashMethod:            config.HashMethod,
		HashEncryptedFile:     config.HashEncryptedFile,
		EncryptionMethod:      config.EncMethod,
	}

	fileService := services.NewFileService(fileServiceParams)

	return Services{
		adminService:         adminService,
		applicationService:   applicationService,
		cryptographicService: cryptographicService,
		fileService:          fileService,
		oauth2Service:        oauth2Service,
		storageService:       minIOService,
		kmsService:           kmsService,
	}

}

func initRepositories(db *gorm.DB) Repositories {
	return Repositories{
		applicationRepository: repository.NewAppsRepository(db),
		adminRepository:       repository.NewAdminRepository(db),
		fileRepository:        repository.NewFileRepository(db),
		fileLogRepository:     repository.NewFileLogRepository(db),
	}

}

type Services struct {
	adminService         services.AdminInterface
	applicationService   services.ApplicationInterface
	cryptographicService services.CryptographicInterface
	fileService          services.FileInterface
	oauth2Service        services.OAuth2Interface
	storageService       services.StorageInterface
	kmsService           services.KMSInterface
}

type Repositories struct {
	applicationRepository repository.ApplicationRepository
	adminRepository       repository.AdminRepository
	fileRepository        repository.FileRepository
	fileLogRepository     repository.FileLogsRepository
}
