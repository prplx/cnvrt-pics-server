// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/services.go
//
// Generated by this command:
//
//	mockgen -source=internal/services/services.go --destination=internal/mocks/services.go --package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	websocket "github.com/gofiber/contrib/websocket"
	types "github.com/prplx/cnvrt/internal/types"
	gomock "go.uber.org/mock/gomock"
)

// MockCommunicator is a mock of Communicator interface.
type MockCommunicator struct {
	ctrl     *gomock.Controller
	recorder *MockCommunicatorMockRecorder
}

// MockCommunicatorMockRecorder is the mock recorder for MockCommunicator.
type MockCommunicatorMockRecorder struct {
	mock *MockCommunicator
}

// NewMockCommunicator creates a new mock instance.
func NewMockCommunicator(ctrl *gomock.Controller) *MockCommunicator {
	mock := &MockCommunicator{ctrl: ctrl}
	mock.recorder = &MockCommunicatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunicator) EXPECT() *MockCommunicatorMockRecorder {
	return m.recorder
}

// AddClient mocks base method.
func (m *MockCommunicator) AddClient(jobID int, connection *websocket.Conn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddClient", jobID, connection)
}

// AddClient indicates an expected call of AddClient.
func (mr *MockCommunicatorMockRecorder) AddClient(jobID, connection any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClient", reflect.TypeOf((*MockCommunicator)(nil).AddClient), jobID, connection)
}

// RemoveClient mocks base method.
func (m *MockCommunicator) RemoveClient(jobID int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveClient", jobID)
}

// RemoveClient indicates an expected call of RemoveClient.
func (mr *MockCommunicatorMockRecorder) RemoveClient(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveClient", reflect.TypeOf((*MockCommunicator)(nil).RemoveClient), jobID)
}

// SendErrorArchiving mocks base method.
func (m *MockCommunicator) SendErrorArchiving(jobID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendErrorArchiving", jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendErrorArchiving indicates an expected call of SendErrorArchiving.
func (mr *MockCommunicatorMockRecorder) SendErrorArchiving(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendErrorArchiving", reflect.TypeOf((*MockCommunicator)(nil).SendErrorArchiving), jobID)
}

// SendErrorProcessing mocks base method.
func (m *MockCommunicator) SendErrorProcessing(jobID, fileID int, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendErrorProcessing", jobID, fileID, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendErrorProcessing indicates an expected call of SendErrorProcessing.
func (mr *MockCommunicatorMockRecorder) SendErrorProcessing(jobID, fileID, fileName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendErrorProcessing", reflect.TypeOf((*MockCommunicator)(nil).SendErrorProcessing), jobID, fileID, fileName)
}

// SendStartArchiving mocks base method.
func (m *MockCommunicator) SendStartArchiving(jobID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendStartArchiving", jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendStartArchiving indicates an expected call of SendStartArchiving.
func (mr *MockCommunicatorMockRecorder) SendStartArchiving(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendStartArchiving", reflect.TypeOf((*MockCommunicator)(nil).SendStartArchiving), jobID)
}

// SendStartProcessing mocks base method.
func (m *MockCommunicator) SendStartProcessing(jobID, fileID int, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendStartProcessing", jobID, fileID, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendStartProcessing indicates an expected call of SendStartProcessing.
func (mr *MockCommunicatorMockRecorder) SendStartProcessing(jobID, fileID, fileName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendStartProcessing", reflect.TypeOf((*MockCommunicator)(nil).SendStartProcessing), jobID, fileID, fileName)
}

// SendSuccessArchiving mocks base method.
func (m *MockCommunicator) SendSuccessArchiving(jobID int, path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSuccessArchiving", jobID, path)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSuccessArchiving indicates an expected call of SendSuccessArchiving.
func (mr *MockCommunicatorMockRecorder) SendSuccessArchiving(jobID, path any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSuccessArchiving", reflect.TypeOf((*MockCommunicator)(nil).SendSuccessArchiving), jobID, path)
}

// SendSuccessFlushing mocks base method.
func (m *MockCommunicator) SendSuccessFlushing(jobID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSuccessFlushing", jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSuccessFlushing indicates an expected call of SendSuccessFlushing.
func (mr *MockCommunicatorMockRecorder) SendSuccessFlushing(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSuccessFlushing", reflect.TypeOf((*MockCommunicator)(nil).SendSuccessFlushing), jobID)
}

// SendSuccessProcessing mocks base method.
func (m *MockCommunicator) SendSuccessProcessing(jobID int, result types.SuccessResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSuccessProcessing", jobID, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSuccessProcessing indicates an expected call of SendSuccessProcessing.
func (mr *MockCommunicatorMockRecorder) SendSuccessProcessing(jobID, result any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSuccessProcessing", reflect.TypeOf((*MockCommunicator)(nil).SendSuccessProcessing), jobID, result)
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// PrintError mocks base method.
func (m *MockLogger) PrintError(err error, properties ...types.AnyMap) {
	m.ctrl.T.Helper()
	varargs := []any{err}
	for _, a := range properties {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "PrintError", varargs...)
}

// PrintError indicates an expected call of PrintError.
func (mr *MockLoggerMockRecorder) PrintError(err any, properties ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{err}, properties...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintError", reflect.TypeOf((*MockLogger)(nil).PrintError), varargs...)
}

// PrintFatal mocks base method.
func (m *MockLogger) PrintFatal(err error, properties ...types.AnyMap) {
	m.ctrl.T.Helper()
	varargs := []any{err}
	for _, a := range properties {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "PrintFatal", varargs...)
}

// PrintFatal indicates an expected call of PrintFatal.
func (mr *MockLoggerMockRecorder) PrintFatal(err any, properties ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{err}, properties...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintFatal", reflect.TypeOf((*MockLogger)(nil).PrintFatal), varargs...)
}

// PrintInfo mocks base method.
func (m *MockLogger) PrintInfo(message string, properties ...types.AnyMap) {
	m.ctrl.T.Helper()
	varargs := []any{message}
	for _, a := range properties {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "PrintInfo", varargs...)
}

// PrintInfo indicates an expected call of PrintInfo.
func (mr *MockLoggerMockRecorder) PrintInfo(message any, properties ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{message}, properties...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintInfo", reflect.TypeOf((*MockLogger)(nil).PrintInfo), varargs...)
}

// Write mocks base method.
func (m *MockLogger) Write(message []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", message)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockLoggerMockRecorder) Write(message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockLogger)(nil).Write), message)
}

// MockProcessor is a mock of Processor interface.
type MockProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockProcessorMockRecorder
}

// MockProcessorMockRecorder is the mock recorder for MockProcessor.
type MockProcessorMockRecorder struct {
	mock *MockProcessor
}

// NewMockProcessor creates a new mock instance.
func NewMockProcessor(ctrl *gomock.Controller) *MockProcessor {
	mock := &MockProcessor{ctrl: ctrl}
	mock.recorder = &MockProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessor) EXPECT() *MockProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockProcessor) Process(ctx context.Context, input types.ImageProcessInput) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Process", ctx, input)
}

// Process indicates an expected call of Process.
func (mr *MockProcessorMockRecorder) Process(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockProcessor)(nil).Process), ctx, input)
}

// Shutdown mocks base method.
func (m *MockProcessor) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockProcessorMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockProcessor)(nil).Shutdown))
}

// Startup mocks base method.
func (m *MockProcessor) Startup() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Startup")
}

// Startup indicates an expected call of Startup.
func (mr *MockProcessorMockRecorder) Startup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Startup", reflect.TypeOf((*MockProcessor)(nil).Startup))
}

// MockArchiver is a mock of Archiver interface.
type MockArchiver struct {
	ctrl     *gomock.Controller
	recorder *MockArchiverMockRecorder
}

// MockArchiverMockRecorder is the mock recorder for MockArchiver.
type MockArchiverMockRecorder struct {
	mock *MockArchiver
}

// NewMockArchiver creates a new mock instance.
func NewMockArchiver(ctrl *gomock.Controller) *MockArchiver {
	mock := &MockArchiver{ctrl: ctrl}
	mock.recorder = &MockArchiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArchiver) EXPECT() *MockArchiverMockRecorder {
	return m.recorder
}

// Archive mocks base method.
func (m *MockArchiver) Archive(jobID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Archive", jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Archive indicates an expected call of Archive.
func (mr *MockArchiverMockRecorder) Archive(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Archive", reflect.TypeOf((*MockArchiver)(nil).Archive), jobID)
}

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// ScheduleFlush mocks base method.
func (m *MockScheduler) ScheduleFlush(jobID int, timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleFlush", jobID, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleFlush indicates an expected call of ScheduleFlush.
func (mr *MockSchedulerMockRecorder) ScheduleFlush(jobID, timeout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleFlush", reflect.TypeOf((*MockScheduler)(nil).ScheduleFlush), jobID, timeout)
}
