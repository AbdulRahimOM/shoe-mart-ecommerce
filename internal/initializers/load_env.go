package initializers

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(envPath string) error {
	// fmt.Println("etered load env variable function in inijtializers")
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Couldn't load env variables")
		return err
	}
	return nil
}
