package util

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[A-Za-z0-9_\.]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[A-Za-z\s]+`).MatchString
)

func ValidateString(s string, minLen int, maxLen int) error {
	if len(s) < minLen || len(s) > maxLen {
		return fmt.Errorf("string %s is not between %d and %d characters", s, minLen, maxLen)
	}

	return nil
}

func ValidateUsername(s string) error {
	if !isValidUsername(s) {
		return fmt.Errorf("username must only contain alphanumeric characters and underscorea and periods")
	}

	return ValidateString(s, 3, 50)
}

func ValidateFullName(s string) error {
	if !isValidFullName(s) {
		return fmt.Errorf("full name must only contain alphabetical characters and spaces")
	}

	return ValidateString(s, 3, 100)
}

func ValidateEmail(s string) error {
	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("email address invalid: %v", err)
	}

	return ValidateString(s, 3, 200)
}

func ValidatePassword(s string) error {
	if !isStrongPassword(s) {
		return fmt.Errorf("password does not meet character requirements")
	}

	return ValidateString(s, 6, 100)
}

func isStrongPassword(s string) bool {
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(s)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(s)
	hasDigit := regexp.MustCompile(`\d`).MatchString(s)
	hasSpecial := regexp.MustCompile(`[^A-Za-z\d]`).MatchString(s)

	return hasLower && hasUpper && hasDigit && hasSpecial
}
