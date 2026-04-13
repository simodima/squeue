package sqs

import (
	"context"

	sqsv2 "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func (d *Driver) Ping() error {
	_, err := d.sqsClient.GetQueueAttributes(context.Background(), &sqsv2.GetQueueAttributesInput{
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
		QueueUrl:       &d.url,
	})
	return err
}
