package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	isValidCurrency = regexp.MustCompile(`^(USD|EUR|CAD)$`).MatchString
)

func ValidateString(val string, minLength int, maxLength int) error {
	n := len(val)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %i - %i characters", minLength, maxLength)
	}
	return nil
}

func ValidateCurrency(val string) error {
	if err := ValidateString(val, 3, 3); err != nil {
		return err
	}

	if !isValidCurrency(val) {
		return fmt.Errorf("currency can only be CAD, EUR, USD")
	}
	return nil
}

func ValidateBalance(val int64) error {
	if val <= 0 {
		return fmt.Errorf("balance must be a positive integer")
	}
	return nil
}

func ValidateUsername(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(val) {
		return fmt.Errorf("must contain only lovercase letters, digits or underscores")
	}
	return nil
}

func ValidatePassword(val string) error {
	return ValidateString(val, 0, 100)
}

func ValidateEmailAddress(val string) error {
	if err := ValidateString(val, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(val); err != nil {
		return fmt.Errorf("entered email address is not valid")
	}
	return nil
}

func ValidateUserFullName(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(val) {
		return fmt.Errorf("Must contain only letters or spaces")
	}
	return nil
}
