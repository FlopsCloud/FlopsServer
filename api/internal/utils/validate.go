package utils

import "unicode"

// IsAlphanumeric checks if a string contains only letters and numbers
func IsAlphanumeric(s string) bool {
	for i, r := range s {

		if i == 0 && unicode.IsNumber(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
