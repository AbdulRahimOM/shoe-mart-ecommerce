package request

type UserSignUpReq struct {
	FirstName string `json:"SignUp_firstName" validate:"required,gte=3"`
	LastName  string `json:"SignUp_lastName"`
	Email     string `json:"SignUp_email" validate:"required,email"`
	Phone     string `json:"SignUp_phone" validate:"required,e164"`
	Password  string `json:"SignUp_password" validate:"gte=3"` //1st password shall only be considered. Its frontend to check both typed passwords are same. No security implicatoions
}
type UserSignInReq struct {
	Email    string `json:"Login_email" validate:"required,email"`
	Password string `json:"Login_password" validate:"required,gte=3"`
}
type AddUserAddress struct {
	UserID      uint   `json:"userId" validate:"required,number"`
	AddressName string `json:"addressName" validate:"required,gte=3"`
	FirstName   string `json:"firstName" validate:"required,gte=3"`
	LastName    string `json:"lastName" validate:"required,gte=3"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,e164"`
	Street      string `json:"street" validate:"required,gte=3"`
	LandMark    string `json:"landmark" validate:"gte=3"`
	City        string `json:"city" validate:"required,gte=3"`
	State       string `json:"state" validate:"required,gte=3"`
	Pincode     string `json:"pincode" validate:"required,pincode"`
}
type EditUserAddress struct {
	ID          uint   `json:"id" validate:"required,number"`
	UserID      uint   `json:"userId" validate:"required,number"`
	AddressName string `json:"addressName" validate:"required,gte=3"`
	FirstName   string `json:"firstName" validate:"required,gte=3"`
	LastName    string `json:"lastName" validate:"required,gte=3"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,e164"`
	Street      string `json:"street" validate:"required,gte=3"`
	LandMark    string `json:"landmark" validate:"gte=0"`
	City        string `json:"city" validate:"required,gte=3"`
	State       string `json:"state" validate:"required,gte=3"`
	Pincode     string `json:"pincode" validate:"required,pincode"`
}
type DeleteUserAddress struct {
	ID     uint `json:"id" validate:"required,number"`
	UserID uint `json:"userId" validate:"required,number"`
}

type GetUserAddresses struct {
	UserID uint `json:"userId" validate:"required,number"`
}

type EditProfileReq struct {
	FirstName string `json:"firstName" validate:"required,gte=3" gorm:"column:firstName;notNull"`
	LastName  string `json:"lastName" validate:"required,gte=3" gorm:"column:lastName"`
	Email     string `json:"email" validate:"required,email" gorm:"column:email;unique;notNull"`
	Phone     string `json:"phone" validate:"required,e164" gorm:"column:phone;notNull"`
}
type ApplyForPasswordResetReq struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
	NewPassword string `json:"newPassword" validate:"required,gte=3"`
}

type VerifyOTPReq struct {
	OTP string `json:"otp" validate:"required,number"`
}
