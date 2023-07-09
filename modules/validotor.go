package modules

import (
	"net/mail"
	"net/url"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidURL(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	return true
}

func isValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	var (
		upper bool
		lower bool
		digit bool
		punct bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsDigit(c):
			digit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			punct = true
		}
	}
	if !upper || !lower || !digit || !punct {
		return false
	}

	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}

	return true
}

func isValidUsername(username string) bool {
	if len(username) < 4 {
		return false
	}
	if !unicode.IsLetter(rune(username[0])) && !unicode.IsDigit(rune(username[0])) {
		return false
	}
	for _, c := range username {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '-' {
			return false
		}
	}
	return true
}
