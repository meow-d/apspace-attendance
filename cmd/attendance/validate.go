package main

import (
	"fmt"
	"strconv"
)

func validateCode(code string) error {
	if len(code) != 3 {
		return fmt.Errorf("code must be 3 digits")
	}
	_, err := strconv.Atoi(code)
	return err
}

func filterNumbers(s string) string {
	var res string
	for _, r := range s {
		if r >= '0' && r <= '9' {
			res += string(r)
		}
	}
	return res
}
