package sqs_test

import (
	"errors"
	"os"
	"testing"

	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/toretto460/squeue/sqs"
	mock_sqs "github.com/toretto460/squeue/sqs/mocks"
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

func (suite *SQSTestSuite) TestNewWithDefaultOptions() {
	_, err := sqs.New()

	suite.Error(err)
	suite.Contains(err.Error(), "missing")
}

func (suite *SQSTestSuite) TestNewWithAClient() {
	sqsDriver, err := sqs.New(sqs.WithClient(suite.sqsMock))

	suite.Nil(err)
	suite.NotNil(sqsDriver)
}

func (suite *SQSTestSuite) TestNewAutoTestConnectionSuccess() {
	suite.sqsMock.
		EXPECT().
		ListQueues(&awssqs.ListQueuesInput{}).
		Return(nil, nil)

	sqsDriver, err := sqs.New(
		sqs.WithClient(suite.sqsMock),
		sqs.AutoTestConnection(),
	)

	suite.Nil(err)
	suite.NotNil(sqsDriver)
}

func (suite *SQSTestSuite) TestNewAutoTestConnectionFail() {
	suite.sqsMock.
		EXPECT().
		ListQueues(&awssqs.ListQueuesInput{}).
		Return(nil, errors.New("error calling aws"))

	sqsDriver, err := sqs.New(
		sqs.WithClient(suite.sqsMock),
		sqs.AutoTestConnection(),
	)

	suite.NotNil(err)
	suite.Nil(sqsDriver)
}

func TestSQSTestSuite(t *testing.T) {
	suite.Run(t, new(SQSTestSuite))
}
