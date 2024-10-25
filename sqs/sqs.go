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

	sharedCredentials *struct {
		filename string
		profile  string
	}
}

func (d *Driver) SetSharedCredentials(filename, profile string) {
	d.sharedCredentials = &struct {
		filename string
		profile  string
	}{filename: filename, profile: profile}
}

func New(options ...Option) (*Driver, error) {
	driver := &Driver{
		testConnectionOnStartup: false,
	}

	for _, o := range options {
		o(driver)
	}

	if driver.sqsClient == nil {
		clientCredentials, err := driver.getCredentials()
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

func (d *Driver) getCredentials() (*credentials.Credentials, error) {
	if d.sharedCredentials != nil {
		return credentials.NewSharedCredentials(
			d.sharedCredentials.filename,
			d.sharedCredentials.profile,
		), nil
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		return credentials.NewEnvCredentials(), nil
	}

	return nil, errors.New(
		"missing shared credentials and AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY env vars",
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

	sqsSession, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, fmt.Errorf("error creating sqs session: %w", err)
	}

	return sqs.New(sqsSession), nil
}
