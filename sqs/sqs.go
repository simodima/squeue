package sqs

import (
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//go:generate mockgen -source=sqs.go -destination=mocks/sqsclient.go
type sqsClient interface {
	DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	GetQueueAttributes(input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error)
}

type Driver struct {
	region                  string
	url                     string
	sqsClient               sqsClient
	testConnectionOnStartup bool
}

func New(options ...Option) (*Driver, error) {
	driver := &Driver{
		testConnectionOnStartup: false,
	}

	for _, o := range options {
		o(driver)
	}

	if driver.sqsClient == nil {
		client, err := createClient(driver.url, driver.region, getCredentials())
		if err != nil {
			return nil, err
		}

		driver.sqsClient = client
	}

	if driver.testConnectionOnStartup {
		if err := driver.Ping(); err != nil {
			return nil, err
		}
	}

	return driver, nil
}

// getCredentials returns explicit credentials when the legacy env vars are set,
// or nil to fall through to the AWS SDK default credential chain (which supports
// ECS/Pod Identity, instance profiles, env vars, and shared credentials files).
func getCredentials() *credentials.Credentials {
	if os.Getenv("AWS_SHARED_CREDENTIALS_FILE") != "" {
		return credentials.NewSharedCredentials("", "")
	}
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		return credentials.NewEnvCredentials()
	}
	// nil tells the SDK to use its built-in default chain, including Pod Identity.
	return nil
}

func createClient(queueUrl string, region string, clientCredentials *credentials.Credentials) (*sqs.SQS, error) {
	parsedUrl, err := url.ParseRequestURI(queueUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating sqs client: %w", err)
	}

	options := session.Options{
		Config: aws.Config{
			Endpoint:    aws.String(fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)),
			Region:      aws.String(region),
			Credentials: clientCredentials,
		},
	}

	return sqs.New(session.Must(session.NewSessionWithOptions(options))), nil
}
