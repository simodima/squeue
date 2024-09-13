package sqs

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/toretto460/squeue/driver"
)

func safeDoOnReceiveMessage(do func(*sqs.ReceiveMessageInput)) func(m any) {
	return func(m any) {
		if SQSMessage, ok := m.(*sqs.ReceiveMessageInput); ok {
			do(SQSMessage)
		}
	}
}

func WithConsumeWaitTimeSeconds(wait int64) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqs.ReceiveMessageInput) {
		message.WaitTimeSeconds = aws.Int64(wait)
	})
}

func WithConsumeVisibilityTimeout(timeout int64) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqs.ReceiveMessageInput) {
		message.VisibilityTimeout = aws.Int64(timeout)
	})
}

func WithConsumeRequestAttemptId(id string) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqs.ReceiveMessageInput) {
		message.ReceiveRequestAttemptId = aws.String(id)
	})
}

func WithConsumeMessageSystemAttributeNames(attributes []string) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqs.ReceiveMessageInput) {
		if len(attributes) == 0 {
			message.MessageSystemAttributeNames = aws.StringSlice(
				[]string{sqs.MessageSystemAttributeNameAll},
			)
		} else {
			message.MessageSystemAttributeNames = aws.StringSlice(attributes)
		}

	})
}

func WithConsumeMessageAttributeNames(names []string) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqs.ReceiveMessageInput) {
		message.MessageAttributeNames = aws.StringSlice(names)
	})
}

func WithConsumeMaxNumberOfMessages(max int) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqs.ReceiveMessageInput) {
		message.MaxNumberOfMessages = aws.Int64(int64(max))
	})
}

func (d *Driver) Consume(ctx context.Context, queue string, opts ...func(message any)) (chan driver.Message, error) {
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

			messages, err := d.fetchMessages(queue, opts...)
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

func (d *Driver) fetchMessages(queue string, opts ...func(message any)) ([][2]string, error) {
	req := &sqs.ReceiveMessageInput{
		VisibilityTimeout:   aws.Int64(90),
		MaxNumberOfMessages: aws.Int64(10),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl: &queue,
	}

	for _, o := range opts {
		o(req)
	}

	msgResult, err := d.sqsClient.ReceiveMessage(req)

	if err != nil {
		return nil, err
	}

	messages := [][2]string{}

	for _, m := range msgResult.Messages {
		messages = append(messages, [2]string{*m.Body, *m.ReceiptHandle})
	}

	return messages, nil
}
