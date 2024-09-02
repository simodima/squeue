package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func safeDoOnSendMessage(do func(*sqs.SendMessageInput)) func(m any) {
	return func(m any) {
		if SQSMessage, ok := m.(*sqs.SendMessageInput); ok {
			do(SQSMessage)
		}
	}
}

func WithEnqueueDelaySeconds(delay int64) func(m any) {
	return safeDoOnSendMessage(func(message *sqs.SendMessageInput) {
		message.DelaySeconds = &delay
	})
}

func WithEnqueueMessageAttributes(attrs map[string]*sqs.MessageAttributeValue) func(m any) {
	return safeDoOnSendMessage(func(message *sqs.SendMessageInput) {
		message.MessageAttributes = attrs
	})
}

func WithEnqueueMessageSystemAttributes(attrs map[string]*sqs.MessageSystemAttributeValue) func(m any) {
	return safeDoOnSendMessage(func(message *sqs.SendMessageInput) {
		message.MessageSystemAttributes = attrs
	})
}

func WithEnqueueMessageDeduplicationId(id string) func(m any) {
	return safeDoOnSendMessage(func(message *sqs.SendMessageInput) {
		message.MessageDeduplicationId = &id
	})
}

func WithEnqueueMessageGroupId(id string) func(m any) {
	return safeDoOnSendMessage(func(message *sqs.SendMessageInput) {
		message.MessageGroupId = &id
	})
}

func (d *Driver) Enqueue(queue string, data []byte, opts ...func(message any)) error {
	if d == nil {
		return fmt.Errorf("invalid SQS client")
	}

	req := &sqs.SendMessageInput{
		MessageBody: aws.String(string(data)),
		QueueUrl:    &queue,
	}

	for _, opt := range opts {
		opt(req)
	}

	_, err := d.sqsClient.SendMessage(req)

	return err
}
