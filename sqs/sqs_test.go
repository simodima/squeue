package sqs_test

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/simodima/squeue/sqs"
	mock_sqs "github.com/simodima/squeue/sqs/mocks"
)

type SQSTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	sqsMock *mock_sqs.MocksqsClient
}

// this function executes before each test case
func (suite *SQSTestSuite) SetupTest() {
	// cleanup environment
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "")
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_ACCESS_SECRET_KEY", "")

	suite.ctrl = gomock.NewController(suite.T())
	suite.sqsMock = mock_sqs.NewMocksqsClient(suite.ctrl)
}

// this function executes after each test case
func (suite *SQSTestSuite) TearDownTest() {
	suite.ctrl.Finish()
	suite.ctrl = nil
	suite.sqsMock = nil
}

func (suite *SQSTestSuite) TestNewWIthUrlAndRegionOption() {
	_, err := sqs.New(
		sqs.WithUrl("https://sqs.eu-central-1.amazonaws.com"),
		sqs.WithRegion("us-east-1"),
		sqs.WithSharedCredentials("test", "default"),
	)

	suite.Nil(err)
}

func (suite *SQSTestSuite) TestNewWithDefaultOptions() {
	_, err := sqs.New()

	suite.Error(err)
	suite.Contains(err.Error(), "missing")
}

func (suite *SQSTestSuite) TestNew_InvalidQueueURL() {
	_, err := sqs.New(
		sqs.WithSharedCredentials("test", "default"),
		sqs.WithUrl("-"),
	)

	suite.Error(err)
	suite.Contains(err.Error(), "invalid URI")
}

func (suite *SQSTestSuite) TestNewWithAClient() {
	sqsDriver, err := sqs.New(sqs.WithClient(suite.sqsMock))

	suite.Nil(err)
	suite.NotNil(sqsDriver)
}

func (suite *SQSTestSuite) TestNewAutoTestConnectionSuccess() {
	suite.sqsMock.
		EXPECT().
		GetQueueAttributes(&awssqs.GetQueueAttributesInput{
			AttributeNames: []*string{aws.String("All")},
			QueueUrl:       aws.String("aws-sqs-queue-url"),
		}).
		Return(&awssqs.GetQueueAttributesOutput{}, nil)

	sqsDriver, err := sqs.New(
		sqs.WithClient(suite.sqsMock),
		sqs.AutoTestConnection(),
		sqs.WithUrl("aws-sqs-queue-url"),
	)

	suite.Nil(err)
	suite.NotNil(sqsDriver)
}

func (suite *SQSTestSuite) TestNewAutoTestConnectionFail() {
	suite.sqsMock.
		EXPECT().
		GetQueueAttributes(&awssqs.GetQueueAttributesInput{
			AttributeNames: []*string{aws.String("All")},
			QueueUrl:       aws.String("aws-sqs-queue-url"),
		}).
		Return(nil, errors.New("error calling aws"))

	sqsDriver, err := sqs.New(
		sqs.WithClient(suite.sqsMock),
		sqs.AutoTestConnection(),
		sqs.WithUrl("aws-sqs-queue-url"),
	)

	suite.NotNil(err)
	suite.Nil(sqsDriver)
}

func (suite *SQSTestSuite) TestEnqueueSuccess() {
	suite.sqsMock.EXPECT().
		SendMessage(&awssqs.SendMessageInput{
			MessageBody:            aws.String("test message"),
			QueueUrl:               aws.String("test-queue"),
			DelaySeconds:           aws.Int64(1),
			MessageDeduplicationId: aws.String("dedup-id-1"),
			MessageGroupId:         aws.String("group-id-1"),
			MessageAttributes: map[string]*awssqs.MessageAttributeValue{
				"tenant": {
					DataType:    aws.String("String"),
					StringValue: aws.String("tenant-1"),
				},
			},
			MessageSystemAttributes: map[string]*awssqs.MessageSystemAttributeValue{
				"request-id": {
					DataType:    aws.String("String"),
					StringValue: aws.String("12345"),
				},
			},
		}).
		Return(nil, nil)

	sqsDriver := must(sqs.New(
		sqs.WithClient(suite.sqsMock),
	))

	err := sqsDriver.Enqueue(
		"test-queue",
		[]byte("test message"),
		sqs.WithEnqueueDelaySeconds(1),
		sqs.WithEnqueueMessageGroupId("group-id-1"),
		sqs.WithEnqueueMessageDeduplicationId("dedup-id-1"),
		sqs.WithEnqueueMessageAttributes(map[string]*awssqs.MessageAttributeValue{
			"tenant": {
				DataType:    aws.String("String"),
				StringValue: aws.String("tenant-1"),
			},
		}),
		sqs.WithEnqueueMessageSystemAttributes(map[string]*awssqs.MessageSystemAttributeValue{
			"request-id": {
				DataType:    aws.String("String"),
				StringValue: aws.String("12345"),
			},
		}),
	)

	suite.Nil(err)
}

func (suite *SQSTestSuite) TestConsumeSuccess() {
	suite.sqsMock.EXPECT().
		ReceiveMessage(&awssqs.ReceiveMessageInput{
			MaxNumberOfMessages:         aws.Int64(9),
			MessageAttributeNames:       []*string{aws.String("All")},
			MessageSystemAttributeNames: []*string{aws.String("All")},
			QueueUrl:                    aws.String("test-queue"),
			ReceiveRequestAttemptId:     aws.String("attempt-1"),
			VisibilityTimeout:           aws.Int64(2),
			WaitTimeSeconds:             aws.Int64(1),
		}).
		Return(&awssqs.ReceiveMessageOutput{
			Messages: []*awssqs.Message{
				{Body: aws.String(`{"id": 1}`), ReceiptHandle: aws.String("1")},
				{Body: aws.String(`{"id": 2}`), ReceiptHandle: aws.String("1")},
				{Body: aws.String(`{"id": 3}`), ReceiptHandle: aws.String("1")},
			},
		}, nil).AnyTimes()

	sqsDriver := must(sqs.New(
		sqs.WithClient(suite.sqsMock),
	))

	ctrl, err := sqsDriver.Consume(
		"test-queue",
		sqs.WithConsumeWaitTimeSeconds(1),
		sqs.WithConsumeVisibilityTimeout(2),
		sqs.WithConsumeRequestAttemptId("attempt-1"),
		sqs.WithConsumeMessageSystemAttributeNames([]string{"All"}),
		sqs.WithConsumeMessageAttributeNames([]string{"All"}),
		sqs.WithConsumeMaxNumberOfMessages(9),
	)

	suite.Nil(err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	stop := make(chan struct{})
	go func() {
		// Assert only the first 3 messages and
		// continue consuming
		i := 1
		for msg := range ctrl.Data() {
			if i < 3 {
				suite.Equal("1", msg.ID)
			}
			if i == 3 {
				stop <- struct{}{}
			}
			i++
		}

		wg.Done()
	}()

	go func() {
		<-stop
		ctrl.Stop()
	}()

	wg.Wait()
}

func (suite *SQSTestSuite) TestEnqueueFail() {
	testQueue := "test-queue"
	suite.sqsMock.EXPECT().
		SendMessage(&awssqs.SendMessageInput{
			MessageBody: aws.String("test message"),
			QueueUrl:    &testQueue,
		}).
		Return(nil, errors.New("error calling aws"))

	sqsDriver, _ := sqs.New(
		sqs.WithClient(suite.sqsMock),
	)

	err := sqsDriver.Enqueue(testQueue, []byte("test message"))
	suite.Error(err)
}

func TestSQSTestSuite(t *testing.T) {
	suite.Run(t, new(SQSTestSuite))
}
