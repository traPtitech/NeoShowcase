// Code generated by MockGen. DO NOT EDIT.
// Source: environment.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/traPtitech/neoshowcase/pkg/domain"
)

// MockEnvironmentRepository is a mock of EnvironmentRepository interface.
type MockEnvironmentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEnvironmentRepositoryMockRecorder
}

// MockEnvironmentRepositoryMockRecorder is the mock recorder for MockEnvironmentRepository.
type MockEnvironmentRepositoryMockRecorder struct {
	mock *MockEnvironmentRepository
}

// NewMockEnvironmentRepository creates a new mock instance.
func NewMockEnvironmentRepository(ctrl *gomock.Controller) *MockEnvironmentRepository {
	mock := &MockEnvironmentRepository{ctrl: ctrl}
	mock.recorder = &MockEnvironmentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnvironmentRepository) EXPECT() *MockEnvironmentRepositoryMockRecorder {
	return m.recorder
}

// GetEnv mocks base method.
func (m *MockEnvironmentRepository) GetEnv(ctx context.Context, applicationID string) ([]*domain.Environment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnv", ctx, applicationID)
	ret0, _ := ret[0].([]*domain.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnv indicates an expected call of GetEnv.
func (mr *MockEnvironmentRepositoryMockRecorder) GetEnv(ctx, applicationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnv", reflect.TypeOf((*MockEnvironmentRepository)(nil).GetEnv), ctx, applicationID)
}

// SetEnv mocks base method.
func (m *MockEnvironmentRepository) SetEnv(ctx context.Context, applicationID, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetEnv", ctx, applicationID, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetEnv indicates an expected call of SetEnv.
func (mr *MockEnvironmentRepositoryMockRecorder) SetEnv(ctx, applicationID, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEnv", reflect.TypeOf((*MockEnvironmentRepository)(nil).SetEnv), ctx, applicationID, key, value)
}
