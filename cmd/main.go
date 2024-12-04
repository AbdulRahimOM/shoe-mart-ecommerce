package main

import (
	"fmt"
	"os"

	_ "github.com/AbdulRahimOM/shoe-mart-ecommerce/docs"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/config"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/di"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	infra "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/infrastructure"
	hashpassword "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/hashPassword"
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

	if err := seedSuperAdmin(); err != nil {
		fmt.Println("Couldn't seed super admin. Error: ", err)
	}

	di.InitializeAndStartAPI()
}

func seedSuperAdmin() error {
	// Check if super admin already exists
	var count int64
	if err := infra.DB.Model(&entities.Admin{}).Count(&count).Error; err != nil {
		return fmt.Errorf("couldn't count admins. DB Error: %v", err)
	}

	if count > 0 {
		fmt.Println("Super admin already exists")
		return nil
	}

	hashpassword, err := hashpassword.Hashpassword(config.InitialSuperAdminPassword)
	if err != nil {
		return err
	}
	admin := &entities.Admin{
		Email:     config.InitialSuperAdminEmail,
		Password:  hashpassword,
		FirstName: config.InitialSuperAdminFirstName,
		LastName:  config.InitialSuperAdminLastName,
		Phone:     config.InitialSuperAdminPhone,
	}

	if err := infra.DB.Create(admin).Error; err != nil {
		return err
	}

	return nil
}
