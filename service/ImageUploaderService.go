package service

import (
	"context"

	ImageKit "github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type ImageUploaderService interface {
	UploadImage(file string) (*uploader.UploadResponse, error)
}

type imageUploaderService struct {
}

func NewImageUploaderService() ImageUploaderService {
	return &imageUploaderService{}
}

func (s *imageUploaderService) UploadImage(file string) (*uploader.UploadResponse, error) {
	ik, err := ImageKit.New()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resp, err := ik.Uploader.Upload(ctx, file, uploader.UploadParam{
		FileName: "randomname",
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
