package extensions

import (
	"fmt"
	"regexp"
)

var (
	isValidStrComb = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("input must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidStrComb(value) {
		return fmt.Errorf("must contain only letters, digits, or underscore")
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}

	if !isValidStrComb(value) {
		return fmt.Errorf("must contain only letters, digits, or underscore")
	}

	return nil
}
