//go:build !local

package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(config *Config, logger *zap.Logger) *gorm.DB {
	logger.Info("Connecting to database", zap.String("dsn", config.DSN()))
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	logger.Info("Connected to database")

	return db
}
