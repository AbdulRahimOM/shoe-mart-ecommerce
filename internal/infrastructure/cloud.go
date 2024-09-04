package infra

import (
	"fmt"

	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/config"

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
