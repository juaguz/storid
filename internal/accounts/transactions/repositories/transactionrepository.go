package repositories

import (
	"github.com/juaguz/storid/internal/accounts/models"
	"github.com/juaguz/storid/internal/accounts/transactions/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionDBRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionDBRepository {
	return &TransactionDBRepository{
		DB: db,
	}
}

func (tr *TransactionDBRepository) Create(transactions []*dto.Transaction) error {
	var transactionsToCreate []models.Transaction
	for _, t := range transactions {

		transaction := models.Transaction{
			ExternalID: t.ExternalID,
			Date:       t.Date,
			Amount:     t.Amount,
			Type:       string(t.Type),
			AccountID:  t.AccountID,
		}

		transactionsToCreate = append(transactionsToCreate, transaction)
	}

	tr.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(transactions)

	return nil
}
