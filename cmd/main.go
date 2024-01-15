package main

import (
	myshoo "MyShoo"
	_ "MyShoo/docs"
	"MyShoo/internal/di"
	infra "MyShoo/internal/infrastructure"
	"MyShoo/internal/initializers"
	"fmt"
)

// @title MyShoo AP0
// @version 1.0
// @description This is a sebcgbdgfrver for MyShoo API.
func main() {
	fmt.Println("Handler ::: main()")

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

	//test purposes
	myshoo.Test()

	di.InitializeAndStartAPI()
}
