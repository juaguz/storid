//go:build local

package db

import (
	_ "embed"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed postgres.sql
var rawSql string

//go:embed accounts.sql
var accounts string

func NewDB(config *Config, logger *zap.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	db.Exec(rawSql)
	db.Exec(accounts)

	return db
}
