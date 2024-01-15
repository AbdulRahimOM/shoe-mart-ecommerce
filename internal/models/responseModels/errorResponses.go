package response

var UnauthorizedAccess = struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}{
	"failed",
	"Unauthorized Access",
}
