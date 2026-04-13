package sqs

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqsv2 "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/simodima/squeue/driver"
)

func safeDoOnReceiveMessage(do func(*sqsv2.ReceiveMessageInput)) func(m any) {
	return func(m any) {
		if SQSMessage, ok := m.(*sqsv2.ReceiveMessageInput); ok {
			do(SQSMessage)
		}
	}
}

func WithConsumeWaitTimeSeconds(wait int32) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqsv2.ReceiveMessageInput) {
		message.WaitTimeSeconds = wait
	})
}

func WithConsumeVisibilityTimeout(timeout int32) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqsv2.ReceiveMessageInput) {
		message.VisibilityTimeout = timeout
	})
}

func WithConsumeRequestAttemptId(id string) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqsv2.ReceiveMessageInput) {
		message.ReceiveRequestAttemptId = aws.String(id)
	})
}

func WithConsumeMessageSystemAttributeNames(attributes []string) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqsv2.ReceiveMessageInput) {
		if len(attributes) == 0 {
			message.MessageSystemAttributeNames = []types.MessageSystemAttributeName{
				types.MessageSystemAttributeNameAll,
			}
		} else {
			names := make([]types.MessageSystemAttributeName, len(attributes))
			for i, a := range attributes {
				names[i] = types.MessageSystemAttributeName(a)
			}
			message.MessageSystemAttributeNames = names
		}
	})
}

func WithConsumeMessageAttributeNames(names []string) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqsv2.ReceiveMessageInput) {
		message.MessageAttributeNames = names
	})
}

func WithConsumeMaxNumberOfMessages(max int) func(m any) {
	return safeDoOnReceiveMessage(func(message *sqsv2.ReceiveMessageInput) {
		message.MaxNumberOfMessages = int32(max)
	})
}

func (d *Driver) Consume(queue string, opts ...func(message any)) (*driver.ConsumerController, error) {
	if d == nil {
		return nil, fmt.Errorf("invalid SQS client")
	}

	ctrl := driver.NewConsumerController()

	go func() {
		for {
			select {
			case <-ctrl.Done():
				return
			case <-time.After(time.Nanosecond):
			}

			messages, err := d.fetchMessages(queue, opts...)
			if err != nil {
				ctrl.Send(driver.Message{
					Error: err,
				})
				continue
			}

			for _, msg := range messages {
				ctrl.Send(driver.Message{
					Body: []byte(msg[0]),
					ID:   msg[1],
				})
			}
		}
	}()

	return ctrl, nil
}

func (d *Driver) fetchMessages(queue string, opts ...func(message any)) ([][2]string, error) {
	req := &sqsv2.ReceiveMessageInput{
		VisibilityTimeout:   90,
		MaxNumberOfMessages: 10,
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl: &queue,
	}

	for _, o := range opts {
		o(req)
	}

	msgResult, err := d.sqsClient.ReceiveMessage(context.Background(), req)
	if err != nil {
		return nil, err
	}

	messages := [][2]string{}
	for _, m := range msgResult.Messages {
		messages = append(messages, [2]string{*m.Body, *m.ReceiptHandle})
	}

	return messages, nil
}
