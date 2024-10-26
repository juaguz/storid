provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "lambda_bucket" {
  bucket = "my-lambda-deployment-bucket-3f7b29d2"
}

resource "aws_s3_bucket" "data_bucket" {
  bucket = var.data_bucket_name
}

output "importer_api_url" {
  description = "The API Gateway URL for the Importer Lambda function"
  value       = "${aws_api_gateway_rest_api.my_api.execution_arn}/prod/importer"
}

output "sender_api_url" {
  description = "The API Gateway URL for the Sender Lambda function"
  value       = "${aws_api_gateway_rest_api.my_api.execution_arn}/prod/sender"
}

# Output for the Importer Lambda ARN
output "importer_lambda_arn" {
  description = "The ARN of the Importer Lambda function"
  value       = aws_lambda_function.importer_lambda.arn
}

# Output for the Sender Lambda ARN
output "sender_lambda_arn" {
  description = "The ARN of the Sender Lambda function"
  value       = aws_lambda_function.sender_lambda.arn
}

# Output for the RDS Endpoint
output "rds_endpoint" {
  description = "The endpoint of the RDS instance"
  value       = aws_db_instance.postgres.endpoint
}

# Output for the Data S3 Bucket Name
output "data_bucket_name" {
  description = "The name of the additional S3 bucket"
  value       = aws_s3_bucket.data_bucket.bucket
}


resource "aws_s3_object" "lambda_importer_zip" {
  bucket = aws_s3_bucket.lambda_bucket.bucket
  key    = var.lambda_importer_zip
  source = var.lambda_importer_zip
  etag   = filemd5(var.lambda_importer_zip)
}

resource "aws_s3_object" "lambda_sender_zip" {
  bucket = aws_s3_bucket.lambda_bucket.bucket
  key    = var.lambda_sender_zip
  source = var.lambda_sender_zip
  etag   = filemd5(var.lambda_sender_zip)
}

resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_exec_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "importer_lambda" {
  filename         = var.lambda_importer_zip
  function_name    = "importer_lambda"
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "main"
  source_code_hash = filebase64sha256(var.lambda_importer_zip)
  runtime          = "provided.al2"
  timeout          = 300
  environment {
    variables = {
      DB_HOST                   = aws_db_instance.postgres.address
      DB_PORT                   = "5432"
      DB_USER                   = var.db_username
      DB_NAME                   = var.db_name
      ENVIRONMENT               = var.environment
      DB_PASSWORD_SECRET_ID     = aws_secretsmanager_secret.db_password_secret.id
      SMTP_CREDENTIALS_SECRET_ID = aws_secretsmanager_secret.smtp_credentials_secret.id
      S3_BUCKET_NAME = aws_s3_bucket.data_bucket.bucket
    }
  }
}

resource "aws_lambda_function" "sender_lambda" {
  filename         = var.lambda_sender_zip
  function_name    = "sender_lambda"
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "main"
  source_code_hash = filebase64sha256(var.lambda_sender_zip)
  runtime          = "provided.al2"
  timeout          = 300
  environment {
    variables = {
      DB_HOST                   = aws_db_instance.postgres.address
      DB_PORT                   = "5432"
      DB_USER                   = var.db_username
      DB_NAME                   = var.db_name
      ENVIRONMENT               = var.environment
      DB_PASSWORD_SECRET_ID     = aws_secretsmanager_secret.db_password_secret.id
      SMTP_CREDENTIALS_SECRET_ID = aws_secretsmanager_secret.smtp_credentials_secret.id
      S3_BUCKET_NAME = aws_s3_bucket.data_bucket.bucket
    }
  }
}

resource "aws_db_instance" "postgres" {
  allocated_storage    = var.db_allocated_storage
  engine               = "postgres"
  engine_version       = "16.4"
  instance_class       = var.db_instance_class
  username             = var.db_username
  password             = var.db_password
  publicly_accessible  = true
  skip_final_snapshot  = true
}

resource "aws_secretsmanager_secret" "db_password_secret" {
  name        = "db-password-secret"
  description = "Password for the RDS database"
}

resource "aws_secretsmanager_secret_version" "db_password_secret_value" {
  secret_id     = aws_secretsmanager_secret.db_password_secret.id
  secret_string = jsonencode({
    DB_PASSWORD = var.db_password
  })
}

resource "aws_secretsmanager_secret" "smtp_credentials_secret" {
  name        = "smtp-credentials-secret"
  description = "SMTP Credentials"
}

resource "aws_secretsmanager_secret_version" "smtp_credentials_secret_value" {
  secret_id     = aws_secretsmanager_secret.smtp_credentials_secret.id
  secret_string = jsonencode({
    SMTP_USERNAME = var.smtp_username,
    SMTP_PASSWORD = var.smtp_password
    SMTP_HOST = var.smtp_host
    SMTP_PORT = var.smtp_port
  })
}

resource "aws_iam_policy" "secretsmanager_access" {
  name        = "SecretsManagerAccessPolicy"
  description = "Allow lambdas to access Secrets Manager"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": [
        "${aws_secretsmanager_secret.db_password_secret.arn}",
        "${aws_secretsmanager_secret.smtp_credentials_secret.arn}"
      ]
    }
  ]
}
EOF
}


resource "aws_iam_policy" "lambda_s3_access" {
  name        = "LambdaS3AccessPolicy"
  description = "Policy to allow Lambda to read objects from S3 bucket"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::${var.data_bucket_name}/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_attach_s3_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_s3_access.arn
}

resource "aws_iam_role_policy_attachment" "attach_secrets_policy" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.secretsmanager_access.arn
}

resource "aws_security_group" "rds_security_group" {
  name        = "rds-security-group"
  description = "Allow inbound access to RDS"

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.allowed_ip]  # Usa la variable de IP aquí
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
# API Gateway Rest API
resource "aws_api_gateway_rest_api" "my_api" {
  name        = "multi-lambda-api"
  description = "API Gateway for multiple Lambda functions"
}

# Create root resources
resource "aws_api_gateway_resource" "importer" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  parent_id   = aws_api_gateway_rest_api.my_api.root_resource_id
  path_part   = "importer"
}

resource "aws_api_gateway_resource" "sender" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  parent_id   = aws_api_gateway_rest_api.my_api.root_resource_id
  path_part   = "sender"
}

# Method for Importer Lambda
resource "aws_api_gateway_method" "importer_post" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.importer.id
  http_method   = "POST"
  authorization = "NONE"
}

# Method for Sender Lambda
resource "aws_api_gateway_method" "sender_post" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.sender.id
  http_method   = "POST"
  authorization = "NONE"
}

# Integration for Importer Lambda
resource "aws_api_gateway_integration" "importer_integration" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.importer.id
  http_method             = aws_api_gateway_method.importer_post.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.importer_lambda.invoke_arn
}

# Integration for Sender Lambda
resource "aws_api_gateway_integration" "sender_integration" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.sender.id
  http_method             = aws_api_gateway_method.sender_post.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.sender_lambda.invoke_arn
}

# Deployment
resource "aws_api_gateway_deployment" "my_api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  depends_on = [
    aws_api_gateway_integration.importer_integration,
    aws_api_gateway_integration.sender_integration
  ]
  stage_name = "prod"
}

# Permissions for Importer Lambda to be invoked by API Gateway
resource "aws_lambda_permission" "importer_permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.importer_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.my_api.execution_arn}/*/POST/importer"
}

# Permissions for Sender Lambda to be invoked by API Gateway
resource "aws_lambda_permission" "sender_permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sender_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.my_api.execution_arn}/*/POST/sender"
}

# Enable CORS for Importer
resource "aws_api_gateway_method" "importer_options" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.importer.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "importer_options_integration" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.importer.id
  http_method             = aws_api_gateway_method.importer_options.http_method
  type                    = "MOCK"
  request_templates = {
    "application/json" = <<EOF
{
  "statusCode": 200
}
EOF
  }
}

resource "aws_api_gateway_method_response" "importer_options_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.importer.id
  http_method = aws_api_gateway_method.importer_options.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }

  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "importer_options_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.importer.id
  http_method = aws_api_gateway_method.importer_options.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,Authorization'"
    "method.response.header.Access-Control-Allow-Methods" = "'OPTIONS,POST'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
}

# Enable CORS for Sender
resource "aws_api_gateway_method" "sender_options" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.sender.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "sender_options_integration" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.sender.id
  http_method             = aws_api_gateway_method.sender_options.http_method
  type                    = "MOCK"
  request_templates = {
    "application/json" = <<EOF
{
  "statusCode": 200
}
EOF
  }
}

resource "aws_api_gateway_method_response" "sender_options_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.sender.id
  http_method = aws_api_gateway_method.sender_options.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

resource "aws_api_gateway_integration_response" "sender_options_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.sender.id
  http_method = aws_api_gateway_method.sender_options.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,Authorization'"
    "method.response.header.Access-Control-Allow-Methods" = "'OPTIONS,POST'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
}

