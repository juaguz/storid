package main

import (
	"context"
	"fmt"
	"log"

	"github.com/juaguz/storid/cmd/sender/internal"
	"github.com/juaguz/storid/internal/accounts/balances/summary"
	"go.uber.org/fx"
)

type NotifierHandler struct {
	summary *summary.Sender
}

func NewNotifierHandler(summary *summary.Sender) *NotifierHandler {
	return &NotifierHandler{summary: summary}
}

func (h *NotifierHandler) Run(ctx context.Context) {
	if err := h.summary.Send(ctx); err != nil {
		log.Fatalf("Error running sender: %v", err)
	}
	fmt.Println("Done")
}

func main() {
	app := fx.New(
		internal.NewApp(),
		fx.Provide(NewNotifierHandler),
		fx.Invoke(func(handler *NotifierHandler) {
			handler.Run(context.Background())
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

	defer app.Stop(context.Background())
}
