package config

import (
	"crypsis-backend/internal/entity"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Databse struct {
	Connection *gorm.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(config Config) (*Databse, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Databse{
		Connection: db,
	}, nil

}

func (d *Databse) RunMigrations() error {
	// Step 1: Migrate Role, Permission , and Category first
	if err := d.Connection.AutoMigrate(
		&entity.Apps{},
		&entity.Admins{},
		&entity.Files{},
		&entity.FileLogs{},
	); err != nil {
		return fmt.Errorf("failed to migrate core tables: %w", err)
	}

	// Step 2: Migrate remaining tables
	if err := d.Connection.AutoMigrate(
		&entity.Metadata{},
	); err != nil {
		return fmt.Errorf("failed to migrate remaining tables: %w", err)
	}
	return nil
}

func (d *Databse) Close() error {
	sqlDB, err := d.Connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
