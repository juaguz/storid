package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ExternalID string    `json:"external_id"`
	Date       time.Time `json:"date"`
	Amount     int       `json:"amount"`
	Type       string    `json:"type"`

	//To simplify the example I will asume that the accountid in the file is the same as the account id in the database
	AccountID uint `json:"account_id"`
}
