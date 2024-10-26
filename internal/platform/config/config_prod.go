//go:build !local

package config

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.uber.org/zap"
)

// Obtener un secreto desde AWS Secrets Manager
func fetchSecretFromAWS(secretID string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return "", err
	}

	client := secretsmanager.NewFromConfig(cfg)
	secretOutput, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return "", err
	}

	return *secretOutput.SecretString, nil
}

func LoadConfig(logger *zap.Logger) *Config {
	logger.Info("Loading config")

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.Error("Invalid DB_PORT", zap.Error(err))
	}

	var dbPassword string
	var smtpUsername, smtpPassword string

	dbSecretID := os.Getenv("DB_PASSWORD_SECRET_ID")
	if dbSecretID == "" {
		logger.Error("DB_PASSWORD_SECRET_ID is not set")
	}

	dbSecretValue, err := fetchSecretFromAWS(dbSecretID)
	if err != nil {
		logger.Error("Failed to fetch DB secret from AWS Secrets Manager: %v", zap.Error(err))
	}

	var dbSecretData map[string]string
	if err := json.Unmarshal([]byte(dbSecretValue), &dbSecretData); err != nil {
		logger.Error("Failed to parse DB secret", zap.Error(err))
	}
	dbPassword = dbSecretData["DB_PASSWORD"]

	smtpSecretID := os.Getenv("SMTP_CREDENTIALS_SECRET_ID")
	if smtpSecretID == "" {
		logger.Fatal("SMTP_CREDENTIALS_SECRET_ID is not set")
	}

	smtpSecretValue, err := fetchSecretFromAWS(smtpSecretID)
	if err != nil {
		logger.Fatal("Failed to fetch SMTP secret from AWS Secrets Manager", zap.Error(err))
	}

	var smtpSecretData map[string]string
	if err := json.Unmarshal([]byte(smtpSecretValue), &smtpSecretData); err != nil {
		logger.Error("Failed to parse SMTP secret", zap.Error(err))
	}
	smtpUsername = smtpSecretData["SMTP_USERNAME"]
	smtpPassword = smtpSecretData["SMTP_PASSWORD"]
	logger.Info("Config loaded", zap.String("SMTP_PORT", smtpSecretData["SMTP_PORT"]))
	dbConfig := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: dbPassword,
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		logger.Error("Error loading AWS config", zap.Error(err))
	}

	s3Config := &S3Config{
		Client: s3.NewFromConfig(awsCfg),
	}

	smtpConfig := &SMTPConfig{
		Host:     smtpSecretData["SMTP_HOST"],
		Port:     smtpSecretData["SMTP_PORT"],
		Username: smtpUsername,
		Password: smtpPassword,
	}

	return &Config{
		DBConfig:   dbConfig,
		S3Config:   s3Config,
		SMTPConfig: smtpConfig,
	}
}
