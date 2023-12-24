// Code generated by MockGen. DO NOT EDIT.
// Source: internal/types/types.go
//
// Generated by this command:
//
//	mockgen -source=internal/types/types.go --destination=internal/mocks/types.go --package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockWebsocketConnection is a mock of WebsocketConnection interface.
type MockWebsocketConnection struct {
	ctrl     *gomock.Controller
	recorder *MockWebsocketConnectionMockRecorder
}

// MockWebsocketConnectionMockRecorder is the mock recorder for MockWebsocketConnection.
type MockWebsocketConnectionMockRecorder struct {
	mock *MockWebsocketConnection
}

// NewMockWebsocketConnection creates a new mock instance.
func NewMockWebsocketConnection(ctrl *gomock.Controller) *MockWebsocketConnection {
	mock := &MockWebsocketConnection{ctrl: ctrl}
	mock.recorder = &MockWebsocketConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebsocketConnection) EXPECT() *MockWebsocketConnectionMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockWebsocketConnection) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockWebsocketConnectionMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWebsocketConnection)(nil).Close))
}

// Params mocks base method.
func (m *MockWebsocketConnection) Params(key string, defaultValue ...string) string {
	m.ctrl.T.Helper()
	varargs := []any{key}
	for _, a := range defaultValue {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Params", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// Params indicates an expected call of Params.
func (mr *MockWebsocketConnectionMockRecorder) Params(key any, defaultValue ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{key}, defaultValue...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Params", reflect.TypeOf((*MockWebsocketConnection)(nil).Params), varargs...)
}

// ReadMessage mocks base method.
func (m *MockWebsocketConnection) ReadMessage() (int, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockWebsocketConnectionMockRecorder) ReadMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockWebsocketConnection)(nil).ReadMessage))
}

// WriteJSON mocks base method.
func (m *MockWebsocketConnection) WriteJSON(v any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteJSON", v)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteJSON indicates an expected call of WriteJSON.
func (mr *MockWebsocketConnectionMockRecorder) WriteJSON(v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteJSON", reflect.TypeOf((*MockWebsocketConnection)(nil).WriteJSON), v)
}
