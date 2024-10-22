package extensions

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidStrComb = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	// isValidPassword = regexp.MustCompile(`^([*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\\S]+$`).MatchString
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
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}

	return nil
}

func ValidateFullname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullname(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}

	// TODO: Write regex for password
	// if !isValidPassword(value) {
	// 	return fmt.Errorf("must contain only letters, digits, or underscore")
	// }

	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return err
	}

	return nil
}
