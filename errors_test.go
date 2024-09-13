package squeue

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorTestSuite struct {
	suite.Suite
}

func (suite *ErrorTestSuite) TestWrapError() {
	err := errors.New("test error")

	tests := []struct {
		code            int
		expectedMessage string
	}{
		{
			code:            ErrUnmarshal,
			expectedMessage: "unmarshalling error",
		},
		{
			code:            ErrMarshal,
			expectedMessage: "marshalling error",
		},
		{
			code:            ErrDriver,
			expectedMessage: "driver error",
		},
		{
			code:            100,
			expectedMessage: "unknown error",
		},
	}

	errData := []byte("test")
	for _, t := range tests {
		wrapped := wrapErr(err, t.code, errData)
		suite.Contains(wrapped.Error(), t.expectedMessage)
		suite.Equal(t.code, wrapped.Code())
		suite.Equal(errData, wrapped.Data())
		suite.Equal(err, wrapped.Unwrap())
	}
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}
