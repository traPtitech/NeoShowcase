// Code generated by MockGen. DO NOT EDIT.
// Source: bus.go

// Package mock_eventbus is a generated GoMock package.
package mock_eventbus

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	eventbus "github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
)

// MockBus is a mock of Bus interface.
type MockBus struct {
	ctrl     *gomock.Controller
	recorder *MockBusMockRecorder
}

// MockBusMockRecorder is the mock recorder for MockBus.
type MockBusMockRecorder struct {
	mock *MockBus
}

// NewMockBus creates a new mock instance.
func NewMockBus(ctrl *gomock.Controller) *MockBus {
	mock := &MockBus{ctrl: ctrl}
	mock.recorder = &MockBusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBus) EXPECT() *MockBusMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockBus) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBusMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBus)(nil).Close), ctx)
}

// Publish mocks base method.
func (m *MockBus) Publish(eventType string, body eventbus.Fields) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Publish", eventType, body)
}

// Publish indicates an expected call of Publish.
func (mr *MockBusMockRecorder) Publish(eventType, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockBus)(nil).Publish), eventType, body)
}

// Subscribe mocks base method.
func (m *MockBus) Subscribe(events ...string) eventbus.Subscription {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range events {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(eventbus.Subscription)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockBusMockRecorder) Subscribe(events ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockBus)(nil).Subscribe), events...)
}

// MockSubscription is a mock of Subscription interface.
type MockSubscription struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionMockRecorder
}

// MockSubscriptionMockRecorder is the mock recorder for MockSubscription.
type MockSubscriptionMockRecorder struct {
	mock *MockSubscription
}

// NewMockSubscription creates a new mock instance.
func NewMockSubscription(ctrl *gomock.Controller) *MockSubscription {
	mock := &MockSubscription{ctrl: ctrl}
	mock.recorder = &MockSubscriptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscription) EXPECT() *MockSubscriptionMockRecorder {
	return m.recorder
}

// Chan mocks base method.
func (m *MockSubscription) Chan() <-chan *eventbus.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chan")
	ret0, _ := ret[0].(<-chan *eventbus.Event)
	return ret0
}

// Chan indicates an expected call of Chan.
func (mr *MockSubscriptionMockRecorder) Chan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chan", reflect.TypeOf((*MockSubscription)(nil).Chan))
}

// Unsubscribe mocks base method.
func (m *MockSubscription) Unsubscribe() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unsubscribe")
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockSubscriptionMockRecorder) Unsubscribe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockSubscription)(nil).Unsubscribe))
}
