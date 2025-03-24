package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFloat64(t *testing.T) {
	// Test cases for FormatFloat64 function
	tests := []struct {
		input     float64
		precision int
		expected  string
	}{
		{123.456789, 2, "123.46"},
		{-123.456789, 3, "-123.457"},
		{0.0, 4, "0.0000"},
		{9876543.210987, 5, "9876543.21099"},
		{-9876543.210987, 0, "-9876543"},
	}

	for _, tt := range tests {
		t.Run("FormatFloat64", func(t *testing.T) {
			result := FormatFloat64(tt.input, tt.precision)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseFloat64(t *testing.T) {
	// Test cases for ParseFloat64 function
	tests := []struct {
		input    string
		expected float64
		valid    bool
	}{
		{"123.456", 123.456, true},
		{"-123.456", -123.456, true},
		{"0.0", 0.0, true},
		{"9876543.210987", 9876543.210987, true},
		{"-9876543.210987", -9876543.210987, true},
		{"abc", 0, false},       // invalid case
		{"123.45.67", 0, false}, // invalid case, more than one decimal point
	}

	for _, tt := range tests {
		t.Run("ParseFloat64", func(t *testing.T) {
			result, valid := ParseFloat64(tt.input)
			assert.Equal(t, tt.valid, valid)
			if tt.valid {
				assert.Equal(t, tt.expected, result)
			} else {
				assert.Equal(t, float64(0), result)
			}
		})
	}
}
