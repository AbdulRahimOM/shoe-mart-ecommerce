package e

import "errors"

var (
	ErrEmailAlreadyUsed   = errors.New("conflict: email already registered")
	ErrEmailNotRegistered = errors.New("this email is not registered")
	ErrInvalidPassword    = errors.New("password mismatch")

	ErrPhoneNumberAlreadyUsed = errors.New("conflict: phone number already used")


	// errors.New("orderID doesn't exist in records")
	ErrOrderIDDoesNotExist = errors.New("orderID doesn't exist in records")
	// errors.New("order does not belong to user")
	ErrOrderNotOfUser = errors.New("order does not belong to user")

	// errors.New("order amount exceeds maximum amount for COD")
	ErrOrderExceedsMaxAmountForCOD = errors.New("order amount exceeds maximum amount for COD")

	ErrCODNotAvailable = errors.New("COD not available")
	ErrEmptyCart       = errors.New("cart is empty")
	ErrOnBindingReq  = errors.New("error binding request")
	ErrOnValidation   = errors.New("error validating the request")
)

