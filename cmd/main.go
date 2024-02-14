package main

import (
	myshoo "MyShoo"
	_ "MyShoo/docs"
	"MyShoo/internal/di"
	"MyShoo/internal/domain/config"
	infra "MyShoo/internal/infrastructure"
	"MyShoo/internal/initializers"
	"fmt"
	"os"
	"path/filepath"
)

// @title MyShoo API
// @version 1.0
// @description E-commerce application (Product: Shoes)
// @securityDefinitions.apikey	BearerTokenAuth
// @in							header
// @name						Authorization
func main() {
	fmt.Println("Program initializing..........")

	//for dev purpose
	myshoo.Test()
	//don't remove this until the project is complete
	//(to avoid frequent changes in main.go file just for testing purpose)
	//which leads to frequent git tracking and commiting

	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	envPath := filepath.Join(executableDir, ".env") //env file is presumed to be alongside the executable
	fmt.Println("envPath: ", envPath)

	if err := initializers.LoadEnvVariables(envPath); err != nil {
		fmt.Println("Couldn't load env variables. Err: ", err)
		fmt.Println("Trying to load .env file for 'go run command mode'....")
		if err := initializers.LoadEnvVariables(".env"); err != nil {//for 'go run cmd/main.go' command mode
			fmt.Println("Couldn't load env variables. Err: ", err)
			fmt.Println("Failed to load .env file even in 'go run command mode'....\nExiting....")
			return
		}else{
			fmt.Println("Successfully loaded .env file for 'go run command mode'")
		
		}
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
