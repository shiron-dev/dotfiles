// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/deps.go
//
// Generated by this command:
//
//	mockgen -source=internal/infrastructure/deps.go -destination=./gen/mock/infrastructure/deps.go
//

// Package mock_infrastructure is a generated GoMock package.
package mock_infrastructure

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDepsInfrastructure is a mock of DepsInfrastructure interface.
type MockDepsInfrastructure struct {
	ctrl     *gomock.Controller
	recorder *MockDepsInfrastructureMockRecorder
	isgomock struct{}
}

// MockDepsInfrastructureMockRecorder is the mock recorder for MockDepsInfrastructure.
type MockDepsInfrastructureMockRecorder struct {
	mock *MockDepsInfrastructure
}

// NewMockDepsInfrastructure creates a new mock instance.
func NewMockDepsInfrastructure(ctrl *gomock.Controller) *MockDepsInfrastructure {
	mock := &MockDepsInfrastructure{ctrl: ctrl}
	mock.recorder = &MockDepsInfrastructureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDepsInfrastructure) EXPECT() *MockDepsInfrastructureMockRecorder {
	return m.recorder
}

// CheckInstalled mocks base method.
func (m *MockDepsInfrastructure) CheckInstalled(name string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInstalled", name)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckInstalled indicates an expected call of CheckInstalled.
func (mr *MockDepsInfrastructureMockRecorder) CheckInstalled(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInstalled", reflect.TypeOf((*MockDepsInfrastructure)(nil).CheckInstalled), name)
}

// OpenWithCode mocks base method.
func (m *MockDepsInfrastructure) OpenWithCode(path ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range path {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OpenWithCode", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// OpenWithCode indicates an expected call of OpenWithCode.
func (mr *MockDepsInfrastructureMockRecorder) OpenWithCode(path ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenWithCode", reflect.TypeOf((*MockDepsInfrastructure)(nil).OpenWithCode), path...)
}
