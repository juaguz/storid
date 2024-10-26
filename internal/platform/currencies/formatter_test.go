package currencies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this is the perfect case to apply PBT (Property-Based Testing)
// to test a wide range of values and edge cases
func TestFloatToCents(t *testing.T) {
	tests := []struct {
		input    float64
		expected int
	}{
		{1.23, 123},
		{0.99, 99},
		{0.0, 0},
		{100.00, 10000},
		{-1.23, -123},
	}

	for _, test := range tests {
		result := FloatToCents(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestCentsToNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected float64
	}{
		{123, 1.23},
		{99, 0.99},
		{0, 0.0},
		{10000, 100.00},
		{-123, -1.23},
	}

	for _, test := range tests {
		result := CentsToNumber[float64](test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestStringToCents(t *testing.T) {
	tests := []struct {
		input       string
		expected    int
		expectError bool
	}{
		{"1.23", 123, false},
		{"0.99", 99, false},
		{"0", 0, false},
		{"-1.23", -123, false},
		{"abc", 0, true},
	}

	for _, test := range tests {
		result, err := StringToCents(test.input)
		if err != nil {
			assert.Equal(t, test.expectError, err)
		}

		assert.Equal(t, test.expected, result)
	}
}

func TestCentsToString(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{123, "1.23"},
		{99, "0.99"},
		{0, "0.00"},
		{10000, "100.00"},
		{-123, "-1.23"},
	}

	for _, test := range tests {
		result := CentsToString(test.input)
		assert.Equal(t, test.expected, result)
	}
}
