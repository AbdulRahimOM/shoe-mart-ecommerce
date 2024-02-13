package response

type SM struct {
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
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
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
