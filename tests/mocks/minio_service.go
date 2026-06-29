package mocks

import "mime/multipart"

type MinioService struct {
	UploadImageFn func(fileHeader *multipart.FileHeader) (string, error)
	DeleteImageFn func(imageURL string) error
}

func (m *MinioService) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	return m.UploadImageFn(fileHeader)
}

func (m *MinioService) DeleteImage(imageURL string) error {
	return m.DeleteImageFn(imageURL)
}
