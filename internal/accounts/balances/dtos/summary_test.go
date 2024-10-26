package dtos

import (
	"testing"

	"github.com/juaguz/storid/internal/platform/months"
	"github.com/stretchr/testify/assert"
)

func TestSummaryBalance_ToMap(t *testing.T) {
	summaryBalance := &SummaryBalance{
		Balance: Balance{
			TotalBalance:     1000,
			AvrDebitAmount:   200,
			AvrCreditAmount:  300,
			TransactionCount: 10,
			AccountID:        1,
		},
		MonthlyBalance: map[months.Month]*MonthlyBalance{
			months.January: {
				Balance: Balance{TransactionCount: 3},
				Month:   int(months.January),
			},
			months.February: {
				Balance: Balance{TransactionCount: 5},
				Month:   int(months.February),
			},
			months.March: {
				Balance: Balance{TransactionCount: 2},
				Month:   int(months.March),
			},
		},
	}

	result := summaryBalance.ToMap()

	assert.Equal(t, 1000, result["TotalBalance"])
	assert.Equal(t, 200, result["AvrDebitAmount"])
	assert.Equal(t, 300, result["AvrCreditAmount"])
	assert.Equal(t, 10, result["TransactionCount"])

	expectedMonthlyBalance := []monthBalance{
		{Month: "January", Count: 3},
		{Month: "February", Count: 5},
		{Month: "March", Count: 2},
	}

	monthlyBalance, ok := result["MonthlyBalance"].([]monthBalance)
	assert.True(t, ok)
	assert.Equal(t, expectedMonthlyBalance, monthlyBalance)

	result = summaryBalance.ToMap(WithDecimalConversion())

	assert.Equal(t, "10.00", result["TotalBalance"])
	assert.Equal(t, "2.00", result["AvrDebitAmount"])
	assert.Equal(t, "3.00", result["AvrCreditAmount"])
	assert.Equal(t, 10, result["TransactionCount"])

	monthlyBalance, ok = result["MonthlyBalance"].([]monthBalance)
	assert.True(t, ok)
	assert.Equal(t, expectedMonthlyBalance, monthlyBalance)
}
