package e

import "errors"

type Error struct {
	Err        error
	StatusCode int
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

var (
	ErrEmailNodtRegistered_401 = Error{Err: errors.New("this email is not registered"), StatusCode: 401}

	ErrInvalidPassword_401 = Error{Err: errors.New("password mismatch"), StatusCode: 401}


// ErrOnBindingReq  = errors.New("error binding request")
// ErrOnValidation   = errors.New("error validating the request")

// 	ErrPhoneNumberAlreadyUsed = errors.New("conflict: phone number already used")

// 	// errors.New("orderID doesn't exist in records")
// 	ErrOrderIDDoesNotExist = errors.New("orderID doesn't exist in records")
// 	// errors.New("order does not belong to user")
// 	ErrOrderNotOfUser = errors.New("order does not belong to user")

// // errors.New("order amount exceeds maximum amount for COD")
)


func DBQueryError(err *error) *Error {
	return &Error{Err: errors.New("db querry err:" + (*err).Error()), StatusCode: 500}
}

func TextError(text string, statusCode int) *Error {
	return &Error{Err: errors.New(text), StatusCode: statusCode}
}
func TextCumError(text string, err error, statusCode int) *Error {
	return &Error{Err: errors.New(text+err.Error()), StatusCode: statusCode}
}
func GetError(error Error) *Error {
	err:=error
	return &err
}
