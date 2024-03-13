package response

var UnauthorizedAccess = struct {
	//Status  string `json:"status"`
	Message string `json:"message"`
}{
	// "failed",
	"Unauthorized Access",
}

var InvalidToken = struct {
	//Status  string `json:"status"`
	Message string `json:"message"`
}{
	// "failed",
	"Invalid token. Access Denied",
}
