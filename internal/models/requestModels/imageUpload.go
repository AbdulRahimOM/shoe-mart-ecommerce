package requestModels

import (
	"os"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type ImageFileReq struct {
	// Ctx          context.Context
	File         *os.File
	UploadParams uploader.UploadParams
	// Directory    string
	// FileName     string
}

type ExcelFileReq struct {
	// Ctx          context.Context
	File         string //tempFilePath
	UploadParams uploader.UploadParams
}
