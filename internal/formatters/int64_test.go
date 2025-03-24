package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatInt64(t *testing.T) {
	// Test cases for FormatInt64 function
	tests := []struct {
		input    int64
		expected string
	}{
		{123, "123"},
		{-456, "-456"},
		{0, "0"},
		{9876543210, "9876543210"},
		{-9876543210, "-9876543210"},
	}

	for _, tt := range tests {
		t.Run("FormatInt64", func(t *testing.T) {
			result := FormatInt64(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseInt64(t *testing.T) {
	// Test cases for ParseInt64 function
	tests := []struct {
		input    string
		expected int64
		valid    bool
	}{
		{"123", 123, true},
		{"-456", -456, true},
		{"0", 0, true},
		{"9876543210", 9876543210, true},
		{"-9876543210", -9876543210, true},
		{"abc", 0, false},    // invalid case
		{"123.45", 0, false}, // invalid case, not an integer
	}

	for _, tt := range tests {
		t.Run("ParseInt64", func(t *testing.T) {
			result, valid := ParseInt64(tt.input)
			assert.Equal(t, tt.valid, valid)
			if tt.valid {
				assert.Equal(t, tt.expected, result)
			} else {
				assert.Equal(t, int64(0), result)
			}
		})
	}
}
