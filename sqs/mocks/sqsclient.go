// Code generated by MockGen. DO NOT EDIT.
// Source: sqs.go

// Package mock_sqs is a generated GoMock package.
package mock_sqs

import (
	reflect "reflect"

	sqs "github.com/aws/aws-sdk-go/service/sqs"
	gomock "github.com/golang/mock/gomock"
)

// MocksqsClient is a mock of sqsClient interface.
type MocksqsClient struct {
	ctrl     *gomock.Controller
	recorder *MocksqsClientMockRecorder
}

// MocksqsClientMockRecorder is the mock recorder for MocksqsClient.
type MocksqsClientMockRecorder struct {
	mock *MocksqsClient
}

// NewMocksqsClient creates a new mock instance.
func NewMocksqsClient(ctrl *gomock.Controller) *MocksqsClient {
	mock := &MocksqsClient{ctrl: ctrl}
	mock.recorder = &MocksqsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksqsClient) EXPECT() *MocksqsClientMockRecorder {
	return m.recorder
}

// DeleteMessage mocks base method.
func (m *MocksqsClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", input)
	ret0, _ := ret[0].(*sqs.DeleteMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MocksqsClientMockRecorder) DeleteMessage(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MocksqsClient)(nil).DeleteMessage), input)
}

// GetQueueAttributes mocks base method.
func (m *MocksqsClient) GetQueueAttributes(input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueueAttributes", input)
	ret0, _ := ret[0].(*sqs.GetQueueAttributesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueueAttributes indicates an expected call of GetQueueAttributes.
func (mr *MocksqsClientMockRecorder) GetQueueAttributes(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueueAttributes", reflect.TypeOf((*MocksqsClient)(nil).GetQueueAttributes), input)
}

// ReceiveMessage mocks base method.
func (m *MocksqsClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveMessage", input)
	ret0, _ := ret[0].(*sqs.ReceiveMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReceiveMessage indicates an expected call of ReceiveMessage.
func (mr *MocksqsClientMockRecorder) ReceiveMessage(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveMessage", reflect.TypeOf((*MocksqsClient)(nil).ReceiveMessage), input)
}

// SendMessage mocks base method.
func (m *MocksqsClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", input)
	ret0, _ := ret[0].(*sqs.SendMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MocksqsClientMockRecorder) SendMessage(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MocksqsClient)(nil).SendMessage), input)
}
