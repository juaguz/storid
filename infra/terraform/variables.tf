# Variable for the database password
variable "db_password" {
  description = "The password for the RDS database"
  type        = string
  sensitive   = true
}

# Variable for SMTP username
variable "smtp_username" {
  description = "Username for SMTP"
  type        = string
  sensitive   = true
}

# Variable for SMTP password
variable "smtp_password" {
  description = "Password for SMTP"
  type        = string
  sensitive   = true
}

# Variable for SMTP host
variable "smtp_host" {
    description = "The SMTP host for sending emails"
    type        = string
    default     = "smtp.gmail.com"
}

variable "smtp_port" {
    description = "The SMTP port for sending emails"
    type        = string
    default     = "587"
}

# Variable for environment configuration
variable "environment" {
  description = "The environment where the application is running (e.g., local, production)"
  type        = string
  default     = "production"
}

# AWS region configuration
variable "aws_region" {
  description = "The AWS region where resources will be deployed"
  type        = string
  default     = "us-east-1"
}

# Database connection variables
variable "db_instance_class" {
  description = "The instance class for the RDS database"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "The allocated storage for the RDS database in GB"
  type        = number
  default     = 20
}

variable "db_engine_version" {
  description = "The engine version for the RDS database"
  type        = string
  default     = "13.3"
}

variable "db_name" {
  description = "The name of the database"
  type        = string
  default     = "mydatabase"
}

variable "db_username" {
  description = "The username for the database"
  type        = string
  default     = "admin"
}

# Variables for the S3 data bucket
variable "data_bucket_name" {
  description = "The name of the S3 bucket for storing additional data"
  type        = string
  default     = "my-application-data-bucket"
}

# Variables for Lambda deployment S3 bucket
variable "lambda_deployment_bucket_name" {
  description = "The name of the S3 bucket for storing Lambda function code"
  type        = string
  default     = "my-lambda-deployment-bucket-3f7b29d2"
}

# Variables for Lambda function ZIP files
variable "lambda_importer_zip" {
  description = "The ZIP file for the Importer Lambda function code"
  type        = string
  default     = "lambda_importer.zip"
}

variable "lambda_sender_zip" {
  description = "The ZIP file for the Sender Lambda function code"
  type        = string
  default     = "lambda_sender.zip"
}

variable "allowed_ip" {
  description = "The IP address allowed to access the RDS instance"
  type        = string
  default     = "0.0.0.0/0"
}
