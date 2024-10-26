package models

type balanceMixin struct {
	TotalBalance     int     `json:"total_balance"`
	TransactionCount int     `json:"transaction_count"`
	AvgCreditAmount  int     `json:"avg_credit_amount"`
	AvgDebitAmount   int     `json:"avg_debit_amount"`
	AccountID        uint    `json:"account_id"`
	Account          Account `json:"account" gorm:"foreignKey:account_id"`
	Count            uint    `json:"count"`
}

func (MonthlyBalance) TableName() string {
	return MonthlyBalanceTable
}

type Balance struct {
	TotalBalance     int     `json:"total_balance"`
	TransactionCount int     `json:"transaction_count"`
	AvgCreditAmount  int     `json:"avg_credit_amount"`
	AvgDebitAmount   int     `json:"avg_debit_amount"`
	AccountID        uint    `json:"account_id"`
	Account          Account `json:"account" gorm:"foreignKey:account_id"`
	Count            uint    `json:"count"`
}

const BalancesTable = "balances"

func (Balance) TableName() string {
	return BalancesTable
}

type MonthlyBalance struct {
	TotalBalance     int     `json:"total_balance"`
	TransactionCount int     `json:"transaction_count"`
	AvgCreditAmount  int     `json:"avg_credit_amount"`
	AvgDebitAmount   int     `json:"avg_debit_amount"`
	AccountID        uint    `json:"account_id"`
	Account          Account `json:"account" gorm:"foreignKey:account_id"`
	Count            uint    `json:"count"`
	Month            int     `json:"month"`
}

const MonthlyBalanceTable = "monthly_balances"
