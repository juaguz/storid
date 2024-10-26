package config

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type S3Config struct {
	Client *s3.Client
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Config struct {
	DBConfig   *DBConfig
	S3Config   *S3Config
	SMTPConfig *SMTPConfig
}
