package sqs

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/toretto460/squeue/driver"
)

func (d *Driver) Consume(ctx context.Context, queue string) (chan driver.Message, error) {
	if d == nil {
		return nil, fmt.Errorf("invalid SQS client")
	}

	results := make(chan driver.Message)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(results)
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
