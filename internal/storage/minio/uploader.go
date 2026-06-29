package minio

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (m *Client) Upload(file multipart.File, size int64, contentType string) (string, error) {
	objectName := uuid.New().String()
	_, err := m.Client.PutObject(context.Background(), m.Bucket, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"http://localhost:9000/%s/%s",
		m.Bucket,
		objectName,
	), nil
}
