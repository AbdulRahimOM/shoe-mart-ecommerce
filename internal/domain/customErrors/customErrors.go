package e

import "errors"

var (
	ErrEmailAlreadyUsed   = errors.New("conflict: email already registered")
	ErrEmailNotRegistered = errors.New("this email is not registered")
	ErrInvalidPassword    = errors.New("password mismatch")

	ErrPhoneNumberAlreadyUsed = errors.New("conflict: phone number already used")
)
