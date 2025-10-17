package main

import "crypsis-backend/internal/config"

func main() {
	// Load environment
	properties := config.LoadProperties()

	// Initialize DB connection
	db, err := config.NewDatabase(config.Config{
		Host:     properties.DBHost,
		Port:     properties.DBPort,
		User:     properties.DBUser,
		Password: properties.DBPassword,
		DBName:   properties.DBName,
		SSLMode:  properties.DBSSLMode,
	})
	if err != nil {
		panic(err)
	}

	// Create AppConfig
	appConfig := &config.AppConfig{
		Properties: properties,
		DB:         db.Connection,
	}

	// Bootstrap the application
	config.BootstrapApp(appConfig)
}
