package services

import (
	"MyShoo/internal/models/requestModels"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type FileUploadService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewFileUploadService(cloudinary *cloudinary.Cloudinary) *FileUploadService {
	return &FileUploadService{
		cloudinary: cloudinary,
	}
}

func (service *FileUploadService) UploadImage(req *requestModels.ImageFileReq) (string, error) {
	result, err := service.cloudinary.Upload.Upload(
		context.Background(),
		req.File,
		req.UploadParams,
	)
	if err != nil {
		log.Println("error happened while uploading image to cloudinary. err: ", err)
		return "", err
	}

	if result.Error.Message != "" {
		log.Println("error happened while uploading image to cloudinary. result.Error: ", result.Error)
		return "", errors.New(fmt.Sprint(result.Error))
	}

	fmt.Println("Url:", result.SecureURL)
	return result.SecureURL, nil
}

// func to upload excel file to cloudinary
func (service *FileUploadService) UploadExcelFile(req *requestModels.ExcelFileReq) (string, error) {
	result, err := service.cloudinary.Upload.Upload(
		context.Background(),
		req.File,
		req.UploadParams,
	)
	if err != nil {
		log.Println("error happened while uploading image to cloudinary. err: ", err)
		return "", err
	}

	if result.Error.Message != "" {
		log.Println("error happened while uploading image to cloudinary. result.Error: ", result.Error)
		return "", errors.New(fmt.Sprint(result.Error))
	}

	fmt.Println("url:", result.SecureURL)

	return result.SecureURL, nil
}

// func to upload invoice to cloudinary
func (service *FileUploadService) UploadInvoice(filePath string, nameToSaveAs string) (string, error) {
	uploadParams := uploader.UploadParams{
		PublicID:  nameToSaveAs,
		Folder:    "invoices",
		Overwrite: true,
	}
	result, err := service.cloudinary.Upload.Upload(
		context.Background(),
		filePath,
		uploadParams,
	)
	if err != nil {
		log.Println("error happened while uploading pdf file to cloudinary. err: ", err)
		return "", err
	}

	if result.Error.Message != "" {
		log.Println("error happened while uploading pdf file to cloudinary. result.Error: ", result.Error)
		return "", errors.New(fmt.Sprint(result.Error))
	}

	url := fmt.Sprintf("https://res.cloudinary.com/ddslfzycp/raw/upload/%s", result.PublicID)
	fmt.Println("url:", url)

	return url, nil
}


// func to upload invoice to cloudinary
func (service *FileUploadService) UploadFile(file interface{}, uploadParams uploader.UploadParams) (string, error) {

	result, err := service.cloudinary.Upload.Upload(
		context.Background(),
		file,
		uploadParams,
	)
	if err != nil {
		log.Println("error happened while uploading image to cloudinary. err: ", err)
		return "", err
	}

	if result.Error.Message != "" {
		log.Println("error happened while uploading image to cloudinary. result.Error: ", result.Error)
		return "", errors.New(fmt.Sprint(result.Error))
	}

	fmt.Println("url:", result.SecureURL)

	return result.SecureURL, nil
}
