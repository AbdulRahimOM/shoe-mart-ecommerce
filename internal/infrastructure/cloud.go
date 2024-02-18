package infra

import (
	"MyShoo/internal/config"
	"fmt"

	"github.com/cloudinary/cloudinary-go"
)

var CloudinaryClient *cloudinary.Cloudinary

func ConnectToCloud() error {
	var err error
	CloudinaryClient, err = cloudinary.NewFromParams(config.CloudinaryCloudName, config.CloudinaryApiKey, config.CloudinaryApiSecret)
	if err != nil {
		fmt.Println("Error creating Cloudinary client:", err)
		return err
	}
	return nil
}
