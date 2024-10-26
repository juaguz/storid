package summary

import (
	"context"
	_ "embed"
	"fmt"

	accountdtos "github.com/juaguz/storid/internal/accounts/dtos"

	"github.com/juaguz/storid/internal/accounts/balances/dtos"
)

type AccountRepository interface {
	GetAccountByID(accountID uint) (*accountdtos.Account, error)
}

type EmailService interface {
	Send(to, subject, template string, variables map[string]interface{}) error
}

type EmailSender struct {
	AccountRepository AccountRepository
	EmailService      EmailService
}

func NewEmailSender(repository AccountRepository, service EmailService) *EmailSender {
	return &EmailSender{
		AccountRepository: repository,
		EmailService:      service,
	}
}

func (se *EmailSender) Send(_ context.Context, accountID uint, summary *dtos.SummaryBalance) error {
	act, err := se.AccountRepository.GetAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("error getting account by ID: %w", err)
	}

	err = se.EmailService.Send(act.Email, "Summary Balance", "summary", summary.ToMap(dtos.WithDecimalConversion()))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
