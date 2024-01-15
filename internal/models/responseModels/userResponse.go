package response

import "MyShoo/internal/domain/entities"

type UserLoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Error   string `json:"error"`
}

type GetUserAddressesResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Addresses []entities.UserAddress `json:"addresses"`
	Error     string                 `json:"error"`
}

type GetProfileResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Profile struct {
		UserDetails entities.UserDetails `json:"userDetails"`
		Addresses   []entities.UserAddress  `json:"addresses"`
	} `json:"profile"`
	Error string `json:"error"`
}
