package usecases

import (
	"errors"
	"unicode"
)

var (
	ErrUsernameInvalidCharacters = errors.New("username contains invalid character")
	ErrUsernameTooShort          = errors.New("username is too short")
	ErrUsernameTooLong           = errors.New("username is too long")
	ErrPasswordInvalidCharacters = errors.New("password contains invalid character")
	ErrPasswordTooShort          = errors.New("password is too short")
	ErrPasswordTooLong           = errors.New("password is too long")
)

const (
	minUsernameLength = 3
	maxUsernameLength = 20
	minPasswordLength = 6
	maxPasswordLength = 30
)

func checkInvalidCharacters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func validateUsername(username string) error {
	if flag := checkInvalidCharacters(username); flag {
		return ErrUsernameInvalidCharacters
	}
	if len(username) < minUsernameLength {
		return ErrUsernameTooShort
	}
	if len(username) > maxUsernameLength {
		return ErrUsernameTooLong
	}
	return nil
}

func validatePassword(password string) error {
	if flag := checkInvalidCharacters(password); flag {
		return ErrPasswordInvalidCharacters
	}
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}
	return nil
}
