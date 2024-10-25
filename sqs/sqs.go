package sqs

import (
	"errors"
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
		clientCredentials, err := getCredentials()
		if err != nil {
			return nil, err
		}

		client, err := createClient(driver.url, driver.region, clientCredentials)
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

func getCredentials() (*credentials.Credentials, error) {
	if os.Getenv("AWS_SHARED_CREDENTIALS_FILE") != "" {
		return credentials.NewSharedCredentials("", ""), nil
	} else if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		return credentials.NewEnvCredentials(), nil
	}

	return nil, errors.New(
		"missing AWS_SHARED_CREDENTIALS_FILE and AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY env vars",
	)
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
