package utils

import "mime/multipart"

type ImageUploader interface {
	Upload(file *multipart.FileHeader) (string, error)
}
