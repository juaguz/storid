package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/juaguz/storid/cmd/sender/internal"
	"github.com/juaguz/storid/internal/accounts/balances/summary"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LambdaEvent struct {
	EventName string `json:"event_name"`
	Payload   string `json:"payload"`
}

func StartLambdaHandler(s *summary.Sender, logger *zap.Logger) {
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

		logger.Info("Starting sender process")
		if err := s.Send(ctx); err != nil {
			logger.Error("Failed to send summary", zap.Error(err))
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Failed to process request",
			}, err
		}

		logger.Info("Summary sent successfully")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "Sender executed successfully",
		}, nil
	})
}

func main() {
	app := fx.New(
		internal.NewApp(),
		fx.Invoke(StartLambdaHandler),
	)

	// Ejecutar la aplicación
	app.Run()
}
