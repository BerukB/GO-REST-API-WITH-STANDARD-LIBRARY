package validation

import (
	"errors"
	"regexp"

	usermodel "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/models"
)

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if !match {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePhone(phone usermodel.PhoneNumber) error {
	phoneRegex := `^((\+251|251|0)?(9|7)\d{8})$`
	match, _ := regexp.MatchString(phoneRegex, string(phone))
	if !match {
		return errors.New("invalid phone number format")
	}
	return nil
}

func ValidateAddress(address string) error {
	if address == "" {
		return nil
	}
	numberRegex := `\d`
	match, _ := regexp.MatchString(numberRegex, address)
	if match {
		return errors.New("address should not contain numbers")
	}
	return nil
}
