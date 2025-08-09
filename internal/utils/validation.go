package utils

import "strings"

// IsNonEmpty returns true if s contains non-whitespace chars
func IsNonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}
