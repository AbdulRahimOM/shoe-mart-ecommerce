package main

import (
	"fmt"
	"os"

	_ "github.com/AbdulRahimOM/shoe-mart-ecommerce/docs"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/config"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/di"
	infra "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/infrastructure"
)

// @title MyShoo API
// @version 1.0
// @description E-commerce application (Product: Shoes)
// @securityDefinitions.apikey	BearerTokenAuth
// @in							header
// @name						Authorization
func main() {
	fmt.Println("Program initializing............")

	defer os.Exit(1)

	if err := config.LoadEnvVariables(); err != nil {
		fmt.Println("Couldn't load env variables. Err: ", err)
	}

	if err := infra.ConnectToDB(); err != nil {
		fmt.Println("Couldn't connect to DB. Error: ", err)
	}
	if err := infra.ConnectToCloud(); err != nil {
		fmt.Println("Couldn't connect to Cloud. Error: ", err)
	}
	if err := config.LoadDeliveryConfig(); err != nil {
		fmt.Println("Couldn't load config. Error: ", err)
	}

	di.InitializeAndStartAPI()
}
