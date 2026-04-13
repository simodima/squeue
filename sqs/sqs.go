package sqs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	sqsv2 "github.com/aws/aws-sdk-go-v2/service/sqs"
)

//go:generate mockgen -source=sqs.go -destination=mocks/sqsclient.go
type sqsClient interface {
	DeleteMessage(ctx context.Context, params *sqsv2.DeleteMessageInput, optFns ...func(*sqsv2.Options)) (*sqsv2.DeleteMessageOutput, error)
	SendMessage(ctx context.Context, params *sqsv2.SendMessageInput, optFns ...func(*sqsv2.Options)) (*sqsv2.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *sqsv2.ReceiveMessageInput, optFns ...func(*sqsv2.Options)) (*sqsv2.ReceiveMessageOutput, error)
	GetQueueAttributes(ctx context.Context, params *sqsv2.GetQueueAttributesInput, optFns ...func(*sqsv2.Options)) (*sqsv2.GetQueueAttributesOutput, error)
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
		client, err := createClient(driver.url, driver.region)
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

func createClient(queueUrl string, region string) (*sqsv2.Client, error) {
	parsedUrl, err := url.ParseRequestURI(queueUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating sqs client: %w", err)
	}

	endpoint := fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithBaseEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	return sqsv2.NewFromConfig(cfg), nil
}
