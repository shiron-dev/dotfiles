// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/git.go
//
// Generated by this command:
//
//	mockgen -source=internal/infrastructure/git.go -destination=./gen/mock/infrastructure/git.go
//

// Package mock_infrastructure is a generated GoMock package.
package mock_infrastructure

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockGitInfrastructure is a mock of GitInfrastructure interface.
type MockGitInfrastructure struct {
	ctrl     *gomock.Controller
	recorder *MockGitInfrastructureMockRecorder
	isgomock struct{}
}

// MockGitInfrastructureMockRecorder is the mock recorder for MockGitInfrastructure.
type MockGitInfrastructureMockRecorder struct {
	mock *MockGitInfrastructure
}

// NewMockGitInfrastructure creates a new mock instance.
func NewMockGitInfrastructure(ctrl *gomock.Controller) *MockGitInfrastructure {
	mock := &MockGitInfrastructure{ctrl: ctrl}
	mock.recorder = &MockGitInfrastructureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitInfrastructure) EXPECT() *MockGitInfrastructureMockRecorder {
	return m.recorder
}

// CheckoutFile mocks base method.
func (m *MockGitInfrastructure) CheckoutFile(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckoutFile", path)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckoutFile indicates an expected call of CheckoutFile.
func (mr *MockGitInfrastructureMockRecorder) CheckoutFile(path any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutFile", reflect.TypeOf((*MockGitInfrastructure)(nil).CheckoutFile), path)
}

// GitDifftool mocks base method.
func (m *MockGitInfrastructure) GitDifftool(ctx context.Context, sout, serror io.Writer, path ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sout, serror}
	for _, a := range path {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GitDifftool", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GitDifftool indicates an expected call of GitDifftool.
func (mr *MockGitInfrastructureMockRecorder) GitDifftool(ctx, sout, serror any, path ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sout, serror}, path...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GitDifftool", reflect.TypeOf((*MockGitInfrastructure)(nil).GitDifftool), varargs...)
}

// IsGitDiff mocks base method.
func (m *MockGitInfrastructure) IsGitDiff(path ...string) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range path {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsGitDiff", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsGitDiff indicates an expected call of IsGitDiff.
func (mr *MockGitInfrastructureMockRecorder) IsGitDiff(path ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsGitDiff", reflect.TypeOf((*MockGitInfrastructure)(nil).IsGitDiff), path...)
}

// SetGitDir mocks base method.
func (m *MockGitInfrastructure) SetGitDir(path string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGitDir", path)
}

// SetGitDir indicates an expected call of SetGitDir.
func (mr *MockGitInfrastructureMockRecorder) SetGitDir(path any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGitDir", reflect.TypeOf((*MockGitInfrastructure)(nil).SetGitDir), path)
}
