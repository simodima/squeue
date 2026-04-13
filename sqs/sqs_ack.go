package sqs

import (
	"context"
	"errors"

	sqsv2 "github.com/aws/aws-sdk-go-v2/service/sqs"
)

func (c *Driver) Ack(queue, messageID string) error {
	if c == nil {
		return errors.New("invalid sqs driver")
	}

	_, err := c.sqsClient.DeleteMessage(context.Background(), &sqsv2.DeleteMessageInput{
		QueueUrl:      &queue,
		ReceiptHandle: &messageID,
	})

	return err
}
