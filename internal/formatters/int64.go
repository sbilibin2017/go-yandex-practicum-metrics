package formatters

import (
	"strconv"
)

// FormatInt64 formats an int64 to a string using strconv.
func FormatInt64(value int64) string {
	return strconv.FormatInt(value, 10) // 10 for decimal base
}

// ParseInt64 parses a string into an int64. Returns an error if the string cannot be parsed.
func ParseInt64(value string) (int64, bool) {
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsedValue, true
}
