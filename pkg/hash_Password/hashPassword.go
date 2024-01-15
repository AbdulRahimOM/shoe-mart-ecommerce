package hashpassword

import "golang.org/x/crypto/bcrypt"

func Hashpassword(password string) (string, error) {
	hashedPWBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPWBytes), err
}
func CompareHashedPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
