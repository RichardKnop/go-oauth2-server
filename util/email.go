package util

import (
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

// ValidateEmail validates an email address based on a regular expression
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}
