package infra

import (
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go"
)

var CloudinaryClient *cloudinary.Cloudinary

func ConnectToCloud() error {
	var cloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
	var api_key = os.Getenv("CLOUDINARY_API_KEY")
	var api_secret = os.Getenv("CLOUDINARY_API_SECRET")
	var err error
	CloudinaryClient, err = cloudinary.NewFromParams(cloudName, api_key, api_secret)
	if err != nil {
		fmt.Println("Error creating Cloudinary client:", err)
		return err
	}
	return nil
}