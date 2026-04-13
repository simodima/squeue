package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqsv2 "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func safeDoOnSendMessage(do func(*sqsv2.SendMessageInput)) func(m any) {
	return func(m any) {
		if SQSMessage, ok := m.(*sqsv2.SendMessageInput); ok {
			do(SQSMessage)
		}
	}
}

func WithEnqueueDelaySeconds(delay int32) func(m any) {
	return safeDoOnSendMessage(func(message *sqsv2.SendMessageInput) {
		message.DelaySeconds = delay
	})
}

func WithEnqueueMessageAttributes(attrs map[string]types.MessageAttributeValue) func(m any) {
	return safeDoOnSendMessage(func(message *sqsv2.SendMessageInput) {
		message.MessageAttributes = attrs
	})
}

func WithEnqueueMessageSystemAttributes(attrs map[string]types.MessageSystemAttributeValue) func(m any) {
	return safeDoOnSendMessage(func(message *sqsv2.SendMessageInput) {
		message.MessageSystemAttributes = attrs
	})
}

func WithEnqueueMessageDeduplicationId(id string) func(m any) {
	return safeDoOnSendMessage(func(message *sqsv2.SendMessageInput) {
		message.MessageDeduplicationId = &id
	})
}

func WithEnqueueMessageGroupId(id string) func(m any) {
	return safeDoOnSendMessage(func(message *sqsv2.SendMessageInput) {
		message.MessageGroupId = &id
	})
}

func (d *Driver) Enqueue(queue string, data []byte, opts ...func(message any)) error {
	if d == nil {
		return fmt.Errorf("invalid SQS client")
	}

	req := &sqsv2.SendMessageInput{
		MessageBody: aws.String(string(data)),
		QueueUrl:    &queue,
	}

	for _, opt := range opts {
		opt(req)
	}

	_, err := d.sqsClient.SendMessage(context.Background(), req)

	return err
}
