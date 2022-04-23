package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IClientS3 interface {
	// Upload a new object to a bucket and returns its URL to view/download.
	UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error)
	// Downloads an existing object from a bucket.
	DownloadObject(ctx context.Context, bucket, fileName string, body io.WriterAt) error
}

type S3 struct {
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3(config Config) (IClientS3, error) {
	session, err := New(config)
	if err != nil {
		return nil, err
	}

	return &S3{
		timeout:    time.Second * 5,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}, nil
}

func (s S3) UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}
	return res.Location, nil
}

func (s S3) DownloadObject(ctx context.Context, bucket, fileName string, body io.WriterAt) error {
	if _, err := s.downloader.DownloadWithContext(ctx, body, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	return nil
}
