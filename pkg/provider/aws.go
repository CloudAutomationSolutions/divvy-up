package provider

import (
	//"bytes"
	//"io/ioutil"
	//pathutil "path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AmazonS3Backend is a backend for Amazon S3
type AmazonS3Backend struct {
	Bucket   string
	Client   *s3.S3
	Prefix   string
	Uploader *s3manager.Uploader
}

// NewAmazonS3Backend creates a new instance of AmazonS3Backend
func NewAmazonS3Backend(bucket string, prefix string, region string, endpoint string) *AmazonS3Backend {
	service := s3.New(session.New(), &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(strings.HasPrefix(endpoint, "http://")),
		S3ForcePathStyle: aws.Bool(endpoint != ""),
	})
	b := &AmazonS3Backend{
		Bucket:   bucket,
		Client:   service,
		Prefix:   cleanPrefix(prefix),
		Uploader: s3manager.NewUploaderWithClient(service),
	}
	return b
}

func (c *AmazonS3Backend) Bootstrap() string {
	return ""
}

func (c *AmazonS3Backend) Distribute() string {
	return ""
}
