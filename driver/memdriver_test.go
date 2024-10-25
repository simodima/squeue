package driver_test

import (
	"testing"
	"time"

	"github.com/simodima/squeue/driver"
	"github.com/stretchr/testify/suite"
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

func (suite *MemoryTestSuite) TestPingDONotBreakAnything() {
	d := driver.NewMemoryDriver(time.Millisecond)
	suite.Nil(d.Ping())
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

	ctrl, err := d.Consume("test")
	suite.Nil(err)

	m1 := <-ctrl.Data()
	m2 := <-ctrl.Data()
	m3 := <-ctrl.Data()

	ctrl.Stop()

	suite.Equal(first, m1.Body)
	suite.Equal(second, m2.Body)
	suite.Equal(third, m3.Body)
}

func TestSQSTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}
