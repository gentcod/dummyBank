// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gentcod/DummyBank/worker (interfaces: TaskDistributor)

// Package mockwk is a generated GoMock package.
package mockwk

import (
	context "context"
	reflect "reflect"

	worker "github.com/gentcod/DummyBank/worker"
	gomock "github.com/golang/mock/gomock"
	asynq "github.com/hibiken/asynq"
)

// MockTaskDistributor is a mock of TaskDistributor interface.
type MockTaskDistributor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskDistributorMockRecorder
}

// MockTaskDistributorMockRecorder is the mock recorder for MockTaskDistributor.
type MockTaskDistributorMockRecorder struct {
	mock *MockTaskDistributor
}

// NewMockTaskDistributor creates a new mock instance.
func NewMockTaskDistributor(ctrl *gomock.Controller) *MockTaskDistributor {
	mock := &MockTaskDistributor{ctrl: ctrl}
	mock.recorder = &MockTaskDistributorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskDistributor) EXPECT() *MockTaskDistributorMockRecorder {
	return m.recorder
}

// DistributeTaskSendVerifyEmail mocks base method.
func (m *MockTaskDistributor) DistributeTaskSendVerifyEmail(arg0 context.Context, arg1 *worker.PayloadSendVerifyEmail, arg2 ...asynq.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DistributeTaskSendVerifyEmail", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DistributeTaskSendVerifyEmail indicates an expected call of DistributeTaskSendVerifyEmail.
func (mr *MockTaskDistributorMockRecorder) DistributeTaskSendVerifyEmail(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DistributeTaskSendVerifyEmail", reflect.TypeOf((*MockTaskDistributor)(nil).DistributeTaskSendVerifyEmail), varargs...)
}
