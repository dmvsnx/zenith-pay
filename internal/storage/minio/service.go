package minio

import (
	"context"
	"mime/multipart"
	"path"

	"github.com/minio/minio-go/v7"
)

type Service interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)
	DeleteImage(imageURL string) error
}

func (m *Client) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	return m.Upload(file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
}

func (m *Client) DeleteImage(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	objectName := path.Base(imageURL)
	if objectName == "." || objectName == "/" {
		return nil
	}

	return m.Client.RemoveObject(context.Background(), m.Bucket, objectName, minio.RemoveObjectOptions{})
}
