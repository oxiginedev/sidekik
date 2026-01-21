package sidekik

import "strings"

// IsStringEmpty checks if a string is empty
func IsStringEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
