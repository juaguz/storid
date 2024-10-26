package dto

import (
	"time"
)

type TransactionType string

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

type Transaction struct {
	ID         uint            `json:"id" gorm:"primary_key"`
	ExternalID string          `json:"external_id"`
	Date       time.Time       `json:"date"`
	Amount     int             `json:"amount"`
	AccountID  uint            `json:"account_id"`
	Type       TransactionType `json:"type"`
}
