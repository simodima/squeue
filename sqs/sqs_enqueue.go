package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

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
