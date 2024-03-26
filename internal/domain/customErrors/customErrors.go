package e

import "errors"

type Error struct {
	Msg        string
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
	errInvalidReq=errors.New("invalid req")
	errInvalidCredentials=errors.New("invalid credentials")

	ErrEmailAlreadyUsed_401        = &Error{Msg: "email already used", Err: errors.New("invalid req"), StatusCode: 401}
	ErrEmailNotRegistered_401 = &Error{Msg: "this email is not registered",Err: errInvalidReq, StatusCode: 401}
	ErrInvalidPassword_401 = &Error{Msg:"password mismatch",Err: errInvalidCredentials, StatusCode: 401}

// ErrOnBindingReq  = errors.New("error binding request")
// ErrOnValidation   = errors.New("error validating the request")

// 	ErrPhoneNumberAlreadyUsed = errors.New("conflict: phone number already used")

// 	// errors.New("orderID doesn't exist in records")
// 	ErrOrderIDDoesNotExist = errors.New("orderID doesn't exist in records")
// 	// errors.New("order does not belong to user")
// 	ErrOrderNotOfUser = errors.New("order does not belong to user")

// // errors.New("order amount exceeds maximum amount for COD")
)

func DBQueryError_500(err *error) *Error {
	return &Error{Msg: "db query err",Err: *err, StatusCode: 500}
}

// func TextError(text string, statusCode int) *Error {
// 	return &Error{Err: errors.New(text), StatusCode: statusCode}
// }


func SetError(msg string, err error, statusCode int) *Error {
	return &Error{Msg:msg,Err: err, StatusCode: statusCode}
}
func GetError(error Error) *Error {
	err := error
	return &err
}
