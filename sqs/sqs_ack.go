package sqs

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func (c *Driver) Ack(queue, messageID string) error {
	if c == nil {
		return errors.New("invalid sqs driver")
	}

	_, err := c.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queue,
		ReceiptHandle: &messageID,
	})

	return err
}
