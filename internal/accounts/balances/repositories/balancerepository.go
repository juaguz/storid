package repositories

import (
	"context"
	"fmt"

	"github.com/juaguz/storid/internal/accounts/balances/dtos"
	"github.com/juaguz/storid/internal/accounts/models"
	"github.com/juaguz/storid/internal/platform/months"
	"gorm.io/gorm"
)

type BalancesDBRepository struct {
	DB *gorm.DB
}

func NewBalancesRepository(db *gorm.DB) *BalancesDBRepository {
	return &BalancesDBRepository{
		DB: db,
	}
}

// GetBalanceByAccountID returns the balance for the given account ID.
// Normally this kind of repository tends to grow a lot,
// I would use a Criteria pattern to filter the results.
func (br *BalancesDBRepository) GetBalanceByAccountID(_ context.Context, accountID uint) (*models.Balance, error) {
	var balance models.Balance
	err := br.DB.Where("account_id = ?", accountID).Find(&balance).Error
	if err != nil {
		return nil, fmt.Errorf("error getting balance by account ID: %w", err)
	}

	return &balance, nil
}

// GetMonthlyBalancesByAccountID returns the monthly balances for the given account ID.
func (br *BalancesDBRepository) GetMonthlyBalancesByAccountID(_ context.Context, accountID uint) (map[months.Month]*dtos.Balance, error) {
	var balances []models.MonthlyBalance
	err := br.DB.Where("account_id = ?", accountID).Order("month").Find(&balances).Error
	if err != nil {
		return nil, err
	}

	monthlyBalances := make(map[months.Month]*dtos.Balance)

	for _, b := range balances {
		m := months.Month(b.Month)

		monthlyBalances[m] = &dtos.Balance{
			TotalBalance:     b.TotalBalance,
			AvrDebitAmount:   b.AvgDebitAmount,
			AvrCreditAmount:  b.AvgDebitAmount,
			TransactionCount: b.TransactionCount,
		}
	}

	return monthlyBalances, nil
}

func (br *BalancesDBRepository) GetSummaryBalance(_ context.Context) (map[uint]*dtos.SummaryBalance, error) {
	var results []struct {
		models.Balance
		models.MonthlyBalance
	}

	err := br.DB.Table(models.BalancesTable).
		Joins("INNER JOIN monthly_balances a ON a.account_id = balances.account_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	balances := make(map[uint]*dtos.SummaryBalance)

	for _, r := range results {
		b := r.Balance

		d, exists := balances[b.AccountID]
		if !exists {
			d = &dtos.SummaryBalance{
				Balance: dtos.Balance{
					TotalBalance:     b.TotalBalance,
					AvrDebitAmount:   b.AvgDebitAmount,
					AvrCreditAmount:  b.AvgCreditAmount,
					TransactionCount: b.TransactionCount,
					AccountID:        b.AccountID,
				},
				MonthlyBalance: make(map[months.Month]*dtos.MonthlyBalance),
			}
			balances[b.AccountID] = d
		}

		mb := &dtos.MonthlyBalance{
			Balance: dtos.Balance{
				TotalBalance:     r.MonthlyBalance.TotalBalance,
				AvrDebitAmount:   r.MonthlyBalance.AvgDebitAmount,
				AvrCreditAmount:  r.MonthlyBalance.AvgCreditAmount,
				TransactionCount: r.MonthlyBalance.TransactionCount,
				AccountID:        r.MonthlyBalance.AccountID,
			},
			Month: r.MonthlyBalance.Month,
		}

		d.MonthlyBalance[months.Month(mb.Month)] = mb
	}

	return balances, nil
}
