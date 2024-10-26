package repositories

import (
	"github.com/juaguz/storid/internal/accounts/dtos"
	"github.com/juaguz/storid/internal/accounts/models"
	"gorm.io/gorm"
)

type AccountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		DB: db,
	}
}

// GetAccountByID returns the account by the given ID.
func (ar *AccountRepository) GetAccountByID(accountID uint) (*dtos.Account, error) {
	var account models.Account
	err := ar.DB.Where("id = ?", accountID).Find(&account).Error
	if err != nil {
		return nil, err
	}

	return &dtos.Account{
		Email: account.Email,
	}, nil
}
