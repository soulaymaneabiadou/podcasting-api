package upload

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Handler struct {
}

func NewS3Handler() *S3Handler {
	return &S3Handler{}
}

func (s *S3Handler) Upload(fh *multipart.FileHeader) (string, error) {
	filename := randomizeFileName(fh.Filename)
	file, err := fh.Open()
	defer file.Close()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("banners/%s", filename)),
		Body:   file,
		ACL:    types.ObjectCannedACL(os.Getenv("S3_ACL")),
	})

	return result.Location, err
}
