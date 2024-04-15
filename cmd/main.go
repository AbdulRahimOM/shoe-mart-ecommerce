package main

import (
	myshoo "MyShoo"
	_ "MyShoo/docs"
	"MyShoo/internal/config"
	"MyShoo/internal/di"
	infra "MyShoo/internal/infrastructure"
	"fmt"
	"os"
)

// @title MyShoo API
// @version 1.0
// @description E-commerce application (Product: Shoes)
// @securityDefinitions.apikey	BearerTokenAuth
// @in							header
// @name						Authorization
func main() {
	fmt.Println("Program initializing............")

	//for dev purpose
	myshoo.Test()
	//don't remove this until the project is complete
	//(to avoid frequent changes in main.go file just for testing purpose)
	//which leads to frequent git tracking and commiting

	if err := config.LoadEnvVariables(); err != nil {
		fmt.Println("Couldn't load env variables. Err: ", err)
		os.Exit(1)
	}

	if err := infra.ConnectToDB(); err != nil {
		fmt.Println("Couldn't connect to DB. Error: ", err)
		os.Exit(1)
	}
	if err := infra.ConnectToCloud(); err != nil {
		fmt.Println("Couldn't connect to Cloud. Error: ", err)
		os.Exit(1)
	}
	if err := config.LoadDeliveryConfig(); err != nil {
		fmt.Println("Couldn't load config. Error: ", err)
		os.Exit(1)
	}

	di.InitializeAndStartAPI()
}
