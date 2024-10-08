package response

var UnauthorizedAccess = struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}{
	"failed",
	"Unauthorized Access",
}

var InvalidToken = struct {
	Status  string `json:"status"`
	Message string `json:"message"`

}{
	"unauthorized",
	"Invalid token. Access Denied",
}

var PasswordNotSet = struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}{
	"failed",
	"Password not set. Please set password to continue",
}