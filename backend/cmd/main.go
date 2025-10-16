package main

import "crypsis-backend/internal/config"

func main() {
	// startApp()

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

	appConfig := &config.AppConfig{
		Properties: properties,
		DB:         db.Connection,
	}

	config.BootstrapApp(appConfig)

}
