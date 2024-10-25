package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (d *Driver) Ping() error {
	_, err := d.sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("All")},
		QueueUrl:       aws.String(d.url),
	})
	return err
}
