package cloudinary

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	cld "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)
	DeleteImage(imageURL string) error
}

type cloudinaryService struct {
	cld *cld.Cloudinary
}

func NewCloudinaryService(cloudinaryURL string) (CloudinaryService, error) {
	c, err := cld.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &cloudinaryService{cld: c}, nil
}

func (s *cloudinaryService) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ctx := context.Background()
	uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "zenith-pay/products",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func (s *cloudinaryService) DeleteImage(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	publicID, err := extractPublicID(imageURL)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

func extractPublicID(imageURL string) (string, error) {
	uploadIdx := strings.Index(imageURL, "/upload/")
	if uploadIdx == -1 {
		return "", errors.New("invalid cloudinary URL: missing /upload/")
	}

	afterUpload := imageURL[uploadIdx+8:]

	slashIdx := strings.Index(afterUpload, "/")
	if slashIdx == -1 {
		return "", errors.New("invalid cloudinary URL: missing version")
	}

	publicIDWithExt := afterUpload[slashIdx+1:]

	extIdx := strings.LastIndex(publicIDWithExt, ".")
	if extIdx != -1 {
		return publicIDWithExt[:extIdx], nil
	}

	return publicIDWithExt, nil
}
