package squeue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/simodima/squeue"
	"github.com/simodima/squeue/driver"
)

type ConsumerTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	driver *MockDriver
}

func (suite *ConsumerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.driver = NewMockDriver(suite.ctrl)
}

func (suite *ConsumerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
	suite.ctrl = nil
	suite.driver = nil
}

func (suite *ConsumerTestSuite) TestNewConsumer() {
	squeue.NewConsumer[*TestMessage](suite.driver, "test-queue")
}

func (suite *ConsumerTestSuite) TestConsumeMessages_DriverError() {
	queue := "test-queue"
	consumer := squeue.NewConsumer[*TestMessage](suite.driver, queue)
	ctx := context.Background()

	suite.driver.
		EXPECT().
		Consume(ctx, queue).
		Return(nil, errors.New("consume error"))

	messages, err := consumer.Consume(ctx)
	suite.Nil(messages)
	suite.Error(err)
}

func (suite *ConsumerTestSuite) TestConsumeMessages_OneMessageWithError() {
	queue := "test-queue"
	consumer := squeue.NewConsumer[*TestMessage](suite.driver, queue)
	ctx := context.Background()

	dMessages := make(chan driver.Message)
	go func() {
		dMessages <- driver.Message{Error: errors.New("error in message")}
		close(dMessages)
	}()

	suite.driver.
		EXPECT().
		Consume(ctx, queue).
		Return(dMessages, nil)

	messages, err := consumer.Consume(ctx)

	suite.NotNil(messages)
	suite.Nil(err)

	messageCount := 0
	for m := range messages {
		messageCount++
		suite.Error(m.Error)
	}

	suite.Equal(1, messageCount)
}

func (suite *ConsumerTestSuite) TestConsumeMessages_OneMessageUnmarshallError() {
	queue := "test-queue"
	consumer := squeue.NewConsumer[*TestMessage](suite.driver, queue)
	ctx := context.Background()

	dMessages := make(chan driver.Message)
	go func() {
		dMessages <- driver.Message{
			Body:  []byte("invalid json"),
			ID:    "1111",
			Error: nil,
		}
		close(dMessages)
	}()

	suite.driver.
		EXPECT().
		Consume(ctx, queue).
		Return(dMessages, nil)

	messages, err := consumer.Consume(ctx)

	suite.NotNil(messages)
	suite.Nil(err)

	m := <-messages

	suite.Error(m.Error)

	var driverErr *squeue.Err
	if errors.As(m.Error, &driverErr) {
		suite.Equal(squeue.ErrUnmarshal, driverErr.Code())
	} else {
		suite.Fail("unmarshal error was expected")
	}
}

func (suite *ConsumerTestSuite) TestConsumeMessages_RealWorldScenarioWithErrors() {
	queue := "test-queue"
	consumer := squeue.NewConsumer[*TestMessage](suite.driver, queue)
	ctx := context.Background()

	dMessages := make(chan driver.Message)
	go func() {
		dMessages <- driver.Message{
			Body:  []byte(`{"name":"test message"}`),
			ID:    "1111",
			Error: nil,
		}

		dMessages <- driver.Message{
			Error: errors.New("wire error"),
		}

		dMessages <- driver.Message{
			Body:  []byte(`{"name":"test another message"}`),
			ID:    "1111",
			Error: nil,
		}

		close(dMessages)
	}()

	suite.driver.
		EXPECT().
		Consume(ctx, queue).
		Return(dMessages, nil)

	messages, err := consumer.Consume(ctx)

	suite.NotNil(messages)
	suite.Nil(err)

	m := <-messages
	suite.Nil(m.Error)
	suite.NotNil(m.Content)
	suite.Equal(m.Content.Name, "test message")

	m = <-messages
	suite.Error(m.Error)
	suite.Contains(m.Error.Error(), "driver error")
	var driverErr *squeue.Err
	if errors.As(m.Error, &driverErr) {
		suite.Equal(squeue.ErrDriver, driverErr.Code())
	} else {
		suite.Fail("driver error was expected")
	}

	m = <-messages
	suite.Nil(m.Error)
	suite.NotNil(m.Content)
	suite.Contains(m.Content.Name, "test another message")
}

func (suite *ConsumerTestSuite) TestAckMessage_WithError() {
	queue := "test-queue"
	consumer := squeue.NewConsumer[*TestMessage](suite.driver, queue)

	msg := squeue.Message[*TestMessage]{
		Content: &TestMessage{},
		ID:      "123",
	}

	suite.driver.
		EXPECT().
		Ack(queue, "123").
		Return(errors.New("ack error"))

	err := consumer.Ack(msg)

	suite.Error(err)
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}
