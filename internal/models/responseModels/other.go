package response

type SME struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
type SMET struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
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

func SuccessSME(message string) SME {
	return SME{
		Status:  "success",
		Message: message,
	}

}
