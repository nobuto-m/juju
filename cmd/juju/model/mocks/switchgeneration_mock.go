// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/model (interfaces: SwitchGenerationCommandAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSwitchGenerationCommandAPI is a mock of SwitchGenerationCommandAPI interface
type MockSwitchGenerationCommandAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSwitchGenerationCommandAPIMockRecorder
}

// MockSwitchGenerationCommandAPIMockRecorder is the mock recorder for MockSwitchGenerationCommandAPI
type MockSwitchGenerationCommandAPIMockRecorder struct {
	mock *MockSwitchGenerationCommandAPI
}

// NewMockSwitchGenerationCommandAPI creates a new mock instance
func NewMockSwitchGenerationCommandAPI(ctrl *gomock.Controller) *MockSwitchGenerationCommandAPI {
	mock := &MockSwitchGenerationCommandAPI{ctrl: ctrl}
	mock.recorder = &MockSwitchGenerationCommandAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSwitchGenerationCommandAPI) EXPECT() *MockSwitchGenerationCommandAPIMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockSwitchGenerationCommandAPI) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockSwitchGenerationCommandAPIMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSwitchGenerationCommandAPI)(nil).Close))
}

// HasNextGeneration mocks base method
func (m *MockSwitchGenerationCommandAPI) HasNextGeneration(arg0 string) (bool, error) {
	ret := m.ctrl.Call(m, "HasNextGeneration", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasNextGeneration indicates an expected call of HasNextGeneration
func (mr *MockSwitchGenerationCommandAPIMockRecorder) HasNextGeneration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNextGeneration", reflect.TypeOf((*MockSwitchGenerationCommandAPI)(nil).HasNextGeneration), arg0)
}
