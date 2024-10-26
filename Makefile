# Vars
S3_BUCKET=my-lambda-deployment-bucket
IMPORTER_ZIP=lambda_importer.zip
SENDER_ZIP=lambda_sender.zip
IMPORTER_FUNCTION=importer_lambda
SENDER_FUNCTION=sender_lambda
GO_BUILD=GOOS=linux GOARCH=amd64

# compiles importer
build-lambda-importer:
	@echo "Building lambda importer"
	$(GO_BUILD) go build -tags production -o lambda_importer -o bootstrap cmd/importer/lambda/main.go
	zip $(IMPORTER_ZIP) bootstrap
	mv $(IMPORTER_ZIP) infra/terraform

# compiles sender
build-lambda-sender:
	@echo "Building lambda sender"
	$(GO_BUILD) go build -tags production -o lambda_sender -o bootstrap cmd/sender/lambda/main.go
	zip $(SENDER_ZIP) bootstrap
	mv $(SENDER_ZIP) infra/terraform

build-docker-cli-importer:
	@echo "Building docker importer"
	docker build --target importer -t importer-app:latest .

run-importer:
	@echo "Running importer"
	docker run --env-file .env.docker --network storid_network -v $(volume) --rm importer-app:latest /app/importer --file=$(file) --mode=$(mode)

build-docker-cli-sender:
	@echo "Building docker sender"
	docker build --target sender -t sender-app:latest .

run-sender:
	@echo "Running sender"
	docker run --env-file .env.docker --rm --network storid_network sender-app:latest

clean:
	@echo "Cleaning up"
	rm -f lambda_importer lambda_sender bootstrap infra/terraform/$(IMPORTER_ZIP) infra/terraform/$(SENDER_ZIP)

# deploys terraform from infra/terraform
deploy-terraform:
	@echo "Deploying terraform"
	cd infra/terraform && terraform plan -var-file="terraform.tfvars" -out="tfplan"  &&  terraform apply "tfplan"

# builds, deploys, and cleans up
all: build-lambda-importer build-lambda-sender deploy-terraform clean

docker: build-docker-cli-importer build-docker-cli-sender

