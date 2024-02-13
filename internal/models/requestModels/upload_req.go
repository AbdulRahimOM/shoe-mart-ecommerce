package request

import (
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type UploadFileReq struct {
	// Ctx          context.Context
	File         interface{}
	UploadParams uploader.UploadParams
}
