package sqs

import "github.com/aws/aws-sdk-go/aws"

type Option func(*Driver)

func WithClient(c sqsClient) Option {
	return func(d *Driver) {
		d.sqsClient = c
	}
}

func AutoTestConnection() Option {
	return func(d *Driver) {
		d.testConnectionOnStartup = true
	}
}

func WithUrl(name string) Option {
	return func(d *Driver) {
		d.url = name
	}
}

func WithRegion(region string) Option {
	return func(d *Driver) {
		d.region = region
	}
}

func WithVisibilityTimeout(val int64) Option {
	return func(d *Driver) {
		d.visibilityTimeout = aws.Int64(val)
	}
}
func WithMaxNumberOfMessages(val int64) Option {
	return func(d *Driver) {
		d.maxNumberOfMessages = aws.Int64(val)
	}
}
