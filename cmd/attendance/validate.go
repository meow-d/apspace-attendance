package main

import (
	"fmt"
	"strconv"
)

// not using their built-in validation thingy

func validateExists(s string) error {
	if s == "" {
		return fmt.Errorf("fields must not be empty")
	}
	return nil
}

func validateUsername(username string) error {
	if err := validateExists(username); err != nil {
		return err
	}
	if len(username) != 8 {
		return fmt.Errorf("username must be 8 characters")
	}
	if username[:2] != "TP" {
		return fmt.Errorf("username must start with TP")
	}
	return nil
}

func validateCode(code string) error {
	if len(code) != 3 {
		return fmt.Errorf("code must be 3 digits")
	}

	_, err := strconv.Atoi(code)
	if err != nil {
		return err
	}

	// shouldn't be needed but just in case
	for _, c := range code {
		if c < '0' || c > '9' {
			return fmt.Errorf("code must contain only digits")
		}
	}

	return nil
}

// run on inputView update
func filterNumbers(s string) string {
	var res string
	for _, r := range s {
		if r >= '0' && r <= '9' {
			res += string(r)
		}
	}
	return res
}
