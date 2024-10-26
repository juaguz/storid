package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name         string        `json:"name"`
	LastName     string        `json:"last_name"`
	Email        string        `json:"email"`
	Transactions []Transaction `json:"transactions"`
}
