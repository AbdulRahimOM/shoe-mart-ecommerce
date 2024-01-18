package requestModels

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type ImageFileReq struct {
	Ctx          context.Context
	File         *os.File
	UploadParams uploader.UploadParams
	Directory    string
	FileName     string
}
