package services

import (
	"MyShoo/internal/models/requestModels"
	"errors"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go"
)

type ImageUploadService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewImageUploadService(cloudinary *cloudinary.Cloudinary) *ImageUploadService {
	return &ImageUploadService{
		cloudinary: cloudinary,
	}
}

func (service *ImageUploadService) UploadImage(file *requestModels.ImageFileReq) (string, error) {
	fmt.Println("HUI file.File: ", file.File)
	fmt.Println("HUI file.UploadParams: ", file.UploadParams)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("HUI file", file)
	result, err := service.cloudinary.Upload.Upload(
		file.Ctx,
		file.File,
		file.UploadParams,
	)
	if err != nil {
		log.Println("error happened while uploading image to cloudinary. err: ", err)
		return "", err
	}

	if result.Error.Message != "" {
		log.Println("error happened while uploading image to cloudinary. result.Error: ", result.Error)
		return "", errors.New(fmt.Sprint(result.Error))
	}

	// var result2 *uploader.UploadResult
	fmt.Println("HUI result:", result)

	return result.SecureURL, nil
}
