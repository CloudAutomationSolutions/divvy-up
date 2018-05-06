package provider

import (
	"bytes"
	"log"
	"net/url"
	"path"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/xid"
)

// AmazonBackend is a backend for Amazon S3
type AmazonBackend struct {
	Bucket    string
	S3Client  *s3.S3
	CFNClient *cloudformation.CloudFormation
	Prefix    string
	Uploader  *s3manager.Uploader
}

type UserSpecifiedCFNParameters struct {
	Endpoint string
}

// NewAmazonBackend creates a new instance of AmazonBackend
func NewAmazonBackend(bucket string, prefix string, region string, endpoint string) *AmazonBackend {
	s3Client := s3.New(session.New(), &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint), //HAS TO BE HTTPS ALWAYS. Maybe add error handling here?
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(endpoint != ""),
	})

	cfnClient := cloudformation.New(session.New(), &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(endpoint != ""),
	})

	b := &AmazonBackend{
		Bucket:    bucket,
		S3Client:  s3Client,
		CFNClient: cfnClient,
		Prefix:    cleanPrefix(prefix),
		Uploader:  s3manager.NewUploaderWithClient(s3Client),
	}
	return b
}

func (b *AmazonBackend) Bootstrap(userSpecifiedParametersLocation string) {
	// TODO: Use this default "https://raw.githubusercontent.com/CloudAutomationSolutions/divvy-up/master/bootstrap-manifests/aws/cfn-template.yaml")
	// The template should be presented as a url.
	// This will be fetched from the github repository.
	// A flag for re-writting the url will be added. if empty, the default will be used
	userParameters := getUserSpecifiedBootstrapConfig(userSpecifiedParametersLocation)
	for i, element := range userParameters {
		cfnParameters := b.getUserSpecifiedCFNParameters(element.BoottrapParameters)

		params := &cloudformation.CreateStackInput{
			// TODO: Something about naming
			StackName: aws.String("divvy-up" + strconv.Itoa(i)), // Required
			Capabilities: []*string{
				aws.String("CAPABILITY_IAM"), // Required
			},
			NotificationARNs: []*string{
				aws.String("NotificationARN"), // Required
			},
			// Should maybe be passed from list?
			Parameters: cfnParameters,
			Tags: []*cloudformation.Tag{
				&cloudformation.Tag{
					Key:   aws.String("application"),
					Value: aws.String("divvy-up"),
				},
			},
			TemplateURL:      aws.String(element.TemplateFile),
			TimeoutInMinutes: aws.Int64(10),
		}

		resp, err := b.CFNClient.CreateStack(params)
		if err != nil {
			// TODO: Evaluate need to print code
			log.Fatalf("Error occurred: %s", err.Error)
		}
		log.Printf("The output from stack creation:\n%v", resp)
	}
}

func (b *AmazonBackend) Distribute(filePath string) string {
	uid := b.generateUID()

	contents := readFile(filePath)

	err := b.putObject(filePath, uid, contents)
	if err != nil {
		log.Fatalf("Cannot upload to bucket: %s", err.Error())
	}

	return "https://<endpoint>/" + uid
}

func (b AmazonBackend) generateUID() string {
	guid := xid.New()

	return guid.String()
}

func (b AmazonBackend) putObject(filePath, uid string, content []byte) error {
	// This tag will be vital for us later on.
	urlencodedUid := url.QueryEscape("divvy-up-uid=" + uid)

	s3Input := &s3manager.UploadInput{
		Bucket:  aws.String(b.Bucket),
		Key:     aws.String(path.Join(b.Prefix, filePath)),
		Body:    bytes.NewBuffer(content),
		Tagging: aws.String(urlencodedUid),
	}
	_, err := b.Uploader.Upload(s3Input)
	return err
}

func (b AmazonBackend) getUserSpecifiedCFNParameters(bootstrapParameters []BootstrapParameterElement) []*cloudformation.Parameter {
	var cfnParameters []*cloudformation.Parameter

	// TODO: User pointers for bootstrapParameters
	for _, element := range bootstrapParameters {
		cfnParameters = append(cfnParameters, &cloudformation.Parameter{
			ParameterKey:     aws.String(element.Key),
			ParameterValue:   aws.String(element.Value),
			UsePreviousValue: aws.Bool(true),
		})
	}
	return cfnParameters
}
