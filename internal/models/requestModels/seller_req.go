package request

type SellerSignUpReq struct {
	FirstName string `json:"SignUp_firstName" validate:"required,gte=3"`
	LastName  string `json:"SignUp_lastName"`
	Email     string `json:"SignUp_email" validate:"required,email"`
	Phone     string `json:"SignUp_phone" validate:"required,e164"`
	Password  string `json:"SignUp_password" validate:"gte=3"` //1st password shall only be considered. Its frontend to check both typed passwords are same. No security implicatoions
}
type SellerSignInReq struct {
	Email    string `json:"Login_email" validate:"required,email"`
	Password string `json:"Login_password" validate:"required,gte=3"`
}
