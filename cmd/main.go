package main

import (
	myshoo "MyShoo"
	_ "MyShoo/docs"
	"MyShoo/internal/di"
	"MyShoo/internal/domain/config"
	infra "MyShoo/internal/infrastructure"
	"MyShoo/internal/initializers"
	"fmt"
)

//	@title MyShoo API
//	@version 1.0
//	@description Main entry point
func main() {
	fmt.Println("Handler ::: main()")

	//for dev purpose
	myshoo.Test()
	//don't remove this until the project is complete
	//(to avoid frequent changes in main.go file just for testing purpose)
	//which leads to frequent git tracking and commiting

	if err := initializers.LoadEnvVariables(); err != nil {
		fmt.Println("Couldn't load env variables")
		return
	}
	if err := infra.ConnectToDB(); err != nil {
		fmt.Println("Couldn't connect to DB")
		return
	}
	if err := infra.ConnectToCloud(); err != nil {
		fmt.Println("Couldn't connect to Cloud")
		return
	}
	if err := config.LoadConfig(); err != nil {
		fmt.Println("Couldn't load config. Error: ", err)
		return
	}

	di.InitializeAndStartAPI()
}
