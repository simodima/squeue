package sqs

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/toretto460/squeue/driver"
)

type sqsClient interface {
	DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	ListQueues(input *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error)
}

type Driver struct {
	region                  string
	url                     string
	sqsClient               sqsClient
	testConnectionOnStartup bool

	visibilityTimeout   *int64
	maxNumberOfMessages *int64
}

func New(options ...Option) (*Driver, error) {
	driver := &Driver{
		visibilityTimeout:       aws.Int64(90),
		maxNumberOfMessages:     aws.Int64(10),
		testConnectionOnStartup: false,
	}

	for _, o := range options {
		o(driver)
	}

	clientCredentials, err := getCredentials()
	if err != nil {
		return nil, err
	}

	client, err := createClient(driver.url, driver.region, clientCredentials)
	if err != nil {
		return nil, err
	}

	driver.sqsClient = client

	if driver.testConnectionOnStartup {
		if err := driver.testConnection(); err != nil {
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
	parsedUrl, err := url.Parse(queueUrl)
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

func (d *Driver) testConnection() error {
	_, err := d.sqsClient.ListQueues(&sqs.ListQueuesInput{})
	return err
}

func (d *Driver) Enqueue(queue string, data []byte) error {
	if d == nil {
		return fmt.Errorf("invalid SQS client")
	}

	_, err := d.sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(data)),
		QueueUrl:    &queue,
	})

	return err
}

func (d *Driver) Consume(ctx context.Context, queue string) (chan driver.Message, error) {
	if d == nil {
		return nil, fmt.Errorf("invalid SQS client")
	}

	results := make(chan driver.Message)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Nanosecond):
			}

			messages, err := d.fetchMessages(queue)
			if err != nil {
				results <- driver.Message{
					Error: err,
				}
				continue
			}

			for _, msg := range messages {
				results <- driver.Message{
					Body: []byte(msg[0]),
					ID:   msg[1],
				}
			}
		}
	}()

	return results, nil
}

func (d *Driver) fetchMessages(queue string) ([][2]string, error) {
	msgResult, err := d.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		VisibilityTimeout:   aws.Int64(90), // transform in options
		MaxNumberOfMessages: aws.Int64(10),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameApproximateReceiveCount),
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		QueueUrl: &queue,
	})

	if err != nil {
		return nil, err
	}

	messages := [][2]string{}

	for _, m := range msgResult.Messages {
		messages = append(messages, [2]string{*m.Body, *m.ReceiptHandle})
	}

	return messages, nil
}
