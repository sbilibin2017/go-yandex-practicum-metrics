package formatters

import (
	"strconv"
)

// FormatFloat64 formats a float64 to a string using strconv with specified precision.
func FormatFloat64(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}

// ParseFloat64 parses a string into a float64. Returns an error if the string cannot be parsed.
func ParseFloat64(value string) (float64, bool) {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return parsedValue, true
}
