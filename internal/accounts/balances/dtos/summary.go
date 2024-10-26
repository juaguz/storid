package dtos

import (
	"github.com/juaguz/storid/internal/platform/currencies"
	"github.com/juaguz/storid/internal/platform/months"
)

type Balance struct {
	TotalBalance     int  `json:"total_balance"`
	AvrDebitAmount   int  `json:"avr_debit_amount"`
	AvrCreditAmount  int  `json:"avr_credit_amount"`
	TransactionCount int  `json:"count"`
	AccountID        uint `json:"account_id"`
}

type MonthlyBalance struct {
	Balance
	Month int `json:"month"`
}

type SummaryBalance struct {
	Balance
	MonthlyBalance map[months.Month]*MonthlyBalance `json:"monthly_balance"`
}

type ToMapOption func(map[string]interface{})

func WithDecimalConversion() ToMapOption {
	return func(result map[string]interface{}) {
		if val, ok := result["TotalBalance"].(int); ok {
			result["TotalBalance"] = currencies.CentsToString(val)
		}
		if val, ok := result["AvrDebitAmount"].(int); ok {
			result["AvrDebitAmount"] = currencies.CentsToString(val)
		}
		if val, ok := result["AvrCreditAmount"].(int); ok {
			result["AvrCreditAmount"] = currencies.CentsToString(val)
		}
	}
}

type monthBalance struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

func (sb *SummaryBalance) ToMap(options ...ToMapOption) map[string]interface{} {
	result := map[string]interface{}{
		"TotalBalance":     sb.TotalBalance,
		"AvrDebitAmount":   sb.AvrDebitAmount,
		"AvrCreditAmount":  sb.AvrCreditAmount,
		"TransactionCount": sb.TransactionCount,
	}

	var monthlyBalances []monthBalance

	for _, m := range months.OrderMonths {
		month := months.Month(m)
		if _, ok := sb.MonthlyBalance[month]; !ok {
			continue
		}
		monthlyBalances = append(monthlyBalances, monthBalance{
			Month: month.String(),
			Count: sb.MonthlyBalance[month].TransactionCount,
		})
	}

	result["MonthlyBalance"] = monthlyBalances

	for _, opt := range options {
		opt(result)
	}

	return result
}
