//go:build local

package config

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadConfig(logger *zap.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		logger.Info("Failed to load .env file", zap.Error(err))
	}

	logger.Info("Loading config")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.Error("Invalid DB_PORT", zap.Error(err))
	}

	dbConfig := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		logger.Error("Error loading AWS config", zap.Error(err))
	}

	s3Config := &S3Config{
		Client: s3.NewFromConfig(awsCfg, func(options *s3.Options) {
			options.BaseEndpoint = aws.String(os.Getenv("LOCAL_STACK_ENDPOINT"))
			options.UsePathStyle = true
		}),
	}

	smtpConfig := &SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	return &Config{
		DBConfig:   dbConfig,
		S3Config:   s3Config,
		SMTPConfig: smtpConfig,
	}
}
