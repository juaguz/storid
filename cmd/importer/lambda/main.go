package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/juaguz/storid/cmd/importer/internal"
	"github.com/juaguz/storid/internal/accounts/transactions/importer"
	"github.com/juaguz/storid/internal/platform/config"
	"github.com/juaguz/storid/internal/platform/dispatcher"
	"github.com/juaguz/storid/internal/platform/filereaders"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LambdaEvent struct {
	EventName string `json:"event_name"`
	FilePath  string `json:"file_path"` // New field for file path
}

func StartLambdaHandler(i *importer.FileImporter, logger *zap.Logger) {
	logger.Info("Starting Lambda handler")
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		logger.Info("Received HTTP request", zap.String("path", request.Path))

		var event LambdaEvent
		if err := json.Unmarshal([]byte(request.Body), &event); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Invalid input",
			}, err
		}

		// Check if file path is provided
		if event.FilePath == "" {
			logger.Error("File path not provided")
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "File path not provided",
			}, nil
		}

		logger.Info("Starting file import process", zap.String("file_path", event.FilePath))
		if err := i.Import(ctx, event.FilePath); err != nil {
			logger.Error("Failed to import file", zap.Error(err))
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Failed to process request",
			}, err
		}

		logger.Info("File imported successfully")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "File imported successfully",
		}, nil
	})
}

func main() {
	app := fx.New(
		internal.NewApp(),
		fx.Provide(
			func(cfg *config.Config) *s3.Client {
				return cfg.S3Config.Client
			},
		),
		fx.Provide(
			func(client *s3.Client) importer.FileReader {
				bucket := os.Getenv("S3_BUCKET_NAME")
				return filereaders.NewS3FileReader(client, bucket)
			},
		),
		fx.Provide(
			func() bool {
				return false
			},
			fx.Annotate(
				dispatcher.NewSimpleEventDispatcher,
				fx.As(new(importer.EventDispatcher)),
				fx.As(new(dispatcher.EventDispatcher)),
			),
		),
		fx.Invoke(StartLambdaHandler),
	)

	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

	defer app.Stop(startCtx)
}
