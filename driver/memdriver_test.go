package driver_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/toretto460/squeue/driver"
)

type MemoryTestSuite struct {
	suite.Suite
}

// // this function executes before each test case
// func (suite *MemoryTestSuite) SetupTest() {
// }

// // this function executes after each test case
// func (suite *MemoryTestSuite) TearDownTest() {
// }

func (suite *MemoryTestSuite) TestNewDriver() {
	suite.NotNil(driver.NewMemoryDriver(time.Millisecond))
}

func (suite *MemoryTestSuite) TestAckDONotBreakAnything() {
	d := driver.NewMemoryDriver(time.Millisecond)
	suite.Nil(d.Ack("test", "1"))
}

func (suite *MemoryTestSuite) TestEnqueueDequeueSuccess() {
	d := driver.NewMemoryDriver(time.Millisecond)

	first := []byte("1")
	err := d.Enqueue("test", first)
	suite.Nil(err)
	second := []byte("2")
	err = d.Enqueue("test", second)
	suite.Nil(err)
	third := []byte("3")
	err = d.Enqueue("test", third)
	suite.Nil(err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messages, err := d.Consume(ctx, "test")
	suite.Nil(err)

	m1 := <-messages
	m2 := <-messages
	m3 := <-messages

	suite.Equal(first, m1.Body)
	suite.Equal(second, m2.Body)
	suite.Equal(third, m3.Body)

}

func TestSQSTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}
