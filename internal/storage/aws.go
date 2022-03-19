package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Endpoint        string
	Region          string
	DisableSSL      bool
}

func New(config Config) (*session.Session, error) {
	return session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
				Endpoint:    &config.Endpoint,
				Region:      &config.Region,
				DisableSSL:  &config.DisableSSL,
			},
		},
	)
}
