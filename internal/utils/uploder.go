package utils

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/saurav11sarkar/001practic/internal/config"
)

type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService(cfg config.Config) (*CloudinaryService, error) {
	client, err := cloudinary.NewFromParams(
		cfg.Cloudinary.CloudName,
		cfg.Cloudinary.APIKey,
		cfg.Cloudinary.APISecret,
	)
	if err != nil {
		return nil, err
	}
	return &CloudinaryService{client: client}, nil
}

func (c *CloudinaryService) UploadFile(ctx context.Context, file interface{}, folderName string) (string, error) {
	resp, err := c.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folderName,
		ResourceType: "auto",
	})

	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
