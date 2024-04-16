package response

import (
	e "MyShoo/internal/domain/customErrors"
	"fmt"
)

type SM struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type NewSM struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SME struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
type SMT struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}
type SMED struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	// Error   string      `json:"error"`
	Data interface{} `json:"data"`
}

func FailedSME(message string, err error) SME {
	return SME{
		Status:  "failed",
		Message: message,
		Error:   err.Error(),
	}
}

func SuccessSM(message string) SM {
	return SM{
		Status:  "success",
		Message: message,
	}

}
func TokenReturn(token *string) SMT {
	return SMT{
		Status:  "success",
		Message: "token generated",
		Token:   *token,
	}

}

func FromError(err error) SME {
	switch err := err.(type) {
	case *e.Error:
		return SME{
			Status:  "failed",
			Message: err.Msg,
			Error:   err.Err.Error(),
		}

	default:
		return SME{
			Status:  "failed",
			Message: "",
			Error:   err.Error(),
		}
	}
}
func FromErrByText(text string) SME {
	return SME{
		Status:  "failed",
		Message: text,
	}
}
func MsgAndError(txt string, err error) SME {
	return SME{
		Status:  "failed",
		Message: txt,
		Error:   err.Error(),
	}
}
func ErrOnBindingReq(err error) SME {
	return SME{
		Status:  "failed",
		Message: "error in binding request",
		Error:   err.Error(),
	}
}
func ErrOnFormValidation(err *[]string) SME {
	return SME{
		Status:  "failed",
		Message: "error in form validation",
		Error:   fmt.Sprint(err),
	}
}
