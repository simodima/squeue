package squeue_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/simodima/squeue"
)

type wrongMessage struct{}

func (e *wrongMessage) UnmarshalJSON(data []byte) error {
	return errors.New("UnmarshalJSON error")
}

func (e *wrongMessage) MarshalJSON() ([]byte, error) {
	return nil, errors.New("MarshalJSON error")
}

type QueueTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	driver *MockDriver
}

func (suite *QueueTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.driver = NewMockDriver(suite.ctrl)
}

func (suite *QueueTestSuite) TearDownTest() {
	suite.ctrl.Finish()
	suite.ctrl = nil
	suite.driver = nil
}

func (suite *QueueTestSuite) TestNewProducer() {
	squeue.NewProducer(suite.driver, "test-queue")
}

func (suite *QueueTestSuite) TestEnqueueMessage_DriverError() {
	queue := "test-queue"
	producer := squeue.NewProducer(suite.driver, queue)

	suite.driver.
		EXPECT().
		Enqueue(queue, []byte(`{"name":"test message"}`)).
		Return(errors.New("producer error"))

	err := producer.Enqueue(&TestMessage{Name: "test message"})
	suite.Error(err)
}

func (suite *QueueTestSuite) TestEnqueueMessage_MarshalingError() {
	queue := "test-queue"
	producer := squeue.NewProducer(suite.driver, queue)

	suite.driver.
		EXPECT().
		Enqueue(queue, nil).
		Times(0)

	err := producer.Enqueue(&wrongMessage{})
	suite.Error(err)
}

func (suite *QueueTestSuite) TestEnqueueMessage() {
	queue := "test-queue"
	producer := squeue.NewProducer(suite.driver, queue)
	message := &TestMessage{Name: "test message"}
	suite.driver.
		EXPECT().
		Enqueue(queue, []byte(`{"name":"test message"}`)).
		Return(nil)

	err := producer.Enqueue(message)
	suite.Nil(err)
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
