package provider

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/xid"
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
	// TODO: Since we want this to be as secure as possible, error on http:// scheme
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

func (b *AmazonS3Backend) Bootstrap() string {
	return ""
}

func (b *AmazonS3Backend) Distribute(filename string) string {
	uid := b.generateUID()

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Cannot read %s: %s", filename, err.Error())
	}

	err = b.putObject(filename, uid, contents)
	if err != nil {
		log.Fatal("Cannot upload to bucket: %s", err.Error())
	}

	return "https://<endpoint>/" + uid
}

func (b AmazonS3Backend) generateUID() string {
	guid := xid.New()

	return guid.String()
}

func (b AmazonS3Backend) putObject(filename, uid string, content []byte) error {
	// This tag will be vital for us later on.
	urlencodedUid := url.QueryEscape("divvy-up-uid=" + uid)

	s3Input := &s3manager.UploadInput{
		Bucket:  aws.String(b.Bucket),
		Key:     aws.String(path.Join(b.Prefix, filename)),
		Body:    bytes.NewBuffer(content),
		Tagging: aws.String(urlencodedUid),
	}
	_, err := b.Uploader.Upload(s3Input)
	return err
}
