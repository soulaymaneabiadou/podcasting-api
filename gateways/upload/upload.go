package upload

import (
	"mime/multipart"
)

type FileHandler interface {
	Upload(file *multipart.FileHeader) (string, error)
}
