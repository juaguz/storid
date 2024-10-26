package summary

import (
	"context"
	"sync"

	"github.com/juaguz/storid/internal/accounts/balances/dtos"
	"go.uber.org/zap"
)

type Notifier interface {
	Send(ctx context.Context, accountID uint, summary *dtos.SummaryBalance) error
}

type SummaryGenerator interface {
	GetSummaryBalance(ctx context.Context) (map[uint]*dtos.SummaryBalance, error)
}

type Sender struct {
	SummaryGenerator SummaryGenerator
	Notifier         []Notifier
	logger           *zap.Logger
}

func NewSender(summaryGenerator SummaryGenerator, notifiers []Notifier, logger *zap.Logger) *Sender {
	return &Sender{
		SummaryGenerator: summaryGenerator,
		Notifier:         notifiers,
		logger:           logger,
	}
}

func (s *Sender) Send(ctx context.Context) error {
	summaries, err := s.SummaryGenerator.GetSummaryBalance(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for accountID, summary := range summaries {
		for _, notifier := range s.Notifier {
			// Increment the WaitGroup counter
			wg.Add(1)

			// Run the sending function as a goroutine
			go func(notifier Notifier, accountID uint, summary *dtos.SummaryBalance) {
				defer wg.Done() // Decrement the counter when done
				err := notifier.Send(ctx, accountID, summary)
				if err != nil {
					s.logger.Error("error sending email", zap.Error(err))
				}
			}(notifier, accountID, summary)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return nil
}
