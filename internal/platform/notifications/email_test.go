package notifications

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTPService_Parse(t *testing.T) {
	s := NewSMTPService(&SMTPConfig{}, nil)
	template := s.templates["summary"]
	variables := map[string]interface{}{
		"TotalBalance":    1000,
		"AvrDebitAmount":  100,
		"AvrCreditAmount": -100,
		"MonthlyBalance": []map[string]interface{}{
			{"Month": "Jan", "Count": 1000},
			{"Month": "Feb", "Count": 100},
			{"Month": "Mar", "Count": -100},
		},
	}

	rendered, err := s.parseTemplate(template, variables)
	assert.NoError(t, err)
	assert.NotNil(t, rendered)

	assert.Contains(t, rendered, "1000")
	assert.Contains(t, rendered, "100")
	assert.Contains(t, rendered, "-100")
	assert.Contains(t, rendered, "Jan")
	assert.Contains(t, rendered, "Feb")
	assert.Contains(t, rendered, "Mar")

}
