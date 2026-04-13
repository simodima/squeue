package sqs_test

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqsv2 "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")

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
	queueUrl := "aws-sqs-queue-url"
	suite.sqsMock.
		EXPECT().
		GetQueueAttributes(gomock.Any(), &sqsv2.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:       &queueUrl,
		}).
		Return(&sqsv2.GetQueueAttributesOutput{}, nil)

	sqsDriver, err := sqs.New(
		sqs.WithClient(suite.sqsMock),
		sqs.AutoTestConnection(),
		sqs.WithUrl(queueUrl),
	)

	suite.Nil(err)
	suite.NotNil(sqsDriver)
}

func (suite *SQSTestSuite) TestNewAutoTestConnectionFail() {
	queueUrl := "aws-sqs-queue-url"
	suite.sqsMock.
		EXPECT().
		GetQueueAttributes(gomock.Any(), &sqsv2.GetQueueAttributesInput{
			AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
			QueueUrl:       &queueUrl,
		}).
		Return(nil, errors.New("error calling aws"))

	sqsDriver, err := sqs.New(
		sqs.WithClient(suite.sqsMock),
		sqs.AutoTestConnection(),
		sqs.WithUrl(queueUrl),
	)

	suite.NotNil(err)
	suite.Nil(sqsDriver)
}

func (suite *SQSTestSuite) TestEnqueueSuccess() {
	testQueue := "test-queue"
	suite.sqsMock.EXPECT().
		SendMessage(gomock.Any(), &sqsv2.SendMessageInput{
			MessageBody:            aws.String("test message"),
			QueueUrl:               &testQueue,
			DelaySeconds:           1,
			MessageDeduplicationId: aws.String("dedup-id-1"),
			MessageGroupId:         aws.String("group-id-1"),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"tenant": {
					DataType:    aws.String("String"),
					StringValue: aws.String("tenant-1"),
				},
			},
			MessageSystemAttributes: map[string]types.MessageSystemAttributeValue{
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
		testQueue,
		[]byte("test message"),
		sqs.WithEnqueueDelaySeconds(1),
		sqs.WithEnqueueMessageGroupId("group-id-1"),
		sqs.WithEnqueueMessageDeduplicationId("dedup-id-1"),
		sqs.WithEnqueueMessageAttributes(map[string]types.MessageAttributeValue{
			"tenant": {
				DataType:    aws.String("String"),
				StringValue: aws.String("tenant-1"),
			},
		}),
		sqs.WithEnqueueMessageSystemAttributes(map[string]types.MessageSystemAttributeValue{
			"request-id": {
				DataType:    aws.String("String"),
				StringValue: aws.String("12345"),
			},
		}),
	)

	suite.Nil(err)
}

func (suite *SQSTestSuite) TestConsumeSuccess() {
	testQueue := "test-queue"
	suite.sqsMock.EXPECT().
		ReceiveMessage(gomock.Any(), &sqsv2.ReceiveMessageInput{
			MaxNumberOfMessages:         9,
			MessageAttributeNames:       []string{"All"},
			MessageSystemAttributeNames: []types.MessageSystemAttributeName{types.MessageSystemAttributeNameAll},
			QueueUrl:                    &testQueue,
			ReceiveRequestAttemptId:     aws.String("attempt-1"),
			VisibilityTimeout:           2,
			WaitTimeSeconds:             1,
		}).
		Return(&sqsv2.ReceiveMessageOutput{
			Messages: []types.Message{
				{Body: aws.String(`{"id": 1}`), ReceiptHandle: aws.String("1")},
				{Body: aws.String(`{"id": 2}`), ReceiptHandle: aws.String("1")},
				{Body: aws.String(`{"id": 3}`), ReceiptHandle: aws.String("1")},
			},
		}, nil).AnyTimes()

	sqsDriver := must(sqs.New(
		sqs.WithClient(suite.sqsMock),
	))

	ctrl, err := sqsDriver.Consume(
		testQueue,
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
		SendMessage(gomock.Any(), &sqsv2.SendMessageInput{
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
