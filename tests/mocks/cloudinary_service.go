package mocks

import "mime/multipart"

type CloudinaryService struct {
	UploadImageFn func(fileHeader *multipart.FileHeader) (string, error)
	DeleteImageFn func(imageURL string) error
}

func (m *CloudinaryService) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	return m.UploadImageFn(fileHeader)
}

func (m *CloudinaryService) DeleteImage(imageURL string) error {
	return m.DeleteImageFn(imageURL)
}
