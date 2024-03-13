package response

import (
	e "MyShoo/internal/domain/customErrors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SM struct {
	//Status  string `json:"status"`
	//Message string `json:"message"`
}
type NewSM struct {
	//Status  string `json:"status"`
	Message string `json:"message"`
}

//	type SME struct {
//		//Status  string `json:"status"`
//		//Message string `json:"message"`
//		//Error   string `json:"error"`
//	}
type SMT struct {
	//Status  string `json:"status"`
	//Message string `json:"message"`
	Token string `json:"token"`
}
type SMED struct {
	//Status  string      `json:"status"`
	//Message string      `json:"message"`
	//Error   string      `json:"error"`
	Data interface{} `json:"data"`
}

// func FailedSME(message string, err error) SME {
// 	return SME{
// 		//Status:  "failed",
// 		//Message: message,
// 		//Error:   err.Error(),
// 	}
// }

func SuccessSM(message string) SM {
	return SM{
		//Status:  "success",
		//Message: message,
	}

}
func FromError(err error) gin.H {
	switch err := err.(type) {
	case *e.Error:
		return gin.H{"error": err.Err.Error()}
	default:
		return gin.H{"error": err.Error()}
	}
}
func FromErrByText(text string) gin.H {
	return gin.H{"error": text}
}
func FromErrByTextCumError(txt string, err error) gin.H {
	return gin.H{"error": txt + err.Error()}
}
func ErrOnBindingReq(err error) gin.H {
	return gin.H{"error": "err in binding req" + err.Error()}
}
func ErrOnFormValidation(err *[]string) gin.H {
	return gin.H{"error": "err in form validation" + fmt.Sprint(err)}

}
