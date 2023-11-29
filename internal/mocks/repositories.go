// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/repositories.go
//
// Generated by this command:
//
//	mockgen -source=internal/repositories/repositories.go --destination=internal/mocks/repositories.go --package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/prplx/lighter.pics/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockJobs is a mock of Jobs interface.
type MockJobs struct {
	ctrl     *gomock.Controller
	recorder *MockJobsMockRecorder
}

// MockJobsMockRecorder is the mock recorder for MockJobs.
type MockJobsMockRecorder struct {
	mock *MockJobs
}

// NewMockJobs creates a new mock instance.
func NewMockJobs(ctrl *gomock.Controller) *MockJobs {
	mock := &MockJobs{ctrl: ctrl}
	mock.recorder = &MockJobsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobs) EXPECT() *MockJobsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockJobs) Create(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockJobsMockRecorder) Create(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockJobs)(nil).Create), ctx)
}

// MockFiles is a mock of Files interface.
type MockFiles struct {
	ctrl     *gomock.Controller
	recorder *MockFilesMockRecorder
}

// MockFilesMockRecorder is the mock recorder for MockFiles.
type MockFilesMockRecorder struct {
	mock *MockFiles
}

// NewMockFiles creates a new mock instance.
func NewMockFiles(ctrl *gomock.Controller) *MockFiles {
	mock := &MockFiles{ctrl: ctrl}
	mock.recorder = &MockFilesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFiles) EXPECT() *MockFilesMockRecorder {
	return m.recorder
}

// CreateBulk mocks base method.
func (m *MockFiles) CreateBulk(ctx context.Context, jobID int, fileNames []string) ([]models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBulk", ctx, jobID, fileNames)
	ret0, _ := ret[0].([]models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBulk indicates an expected call of CreateBulk.
func (mr *MockFilesMockRecorder) CreateBulk(ctx, jobID, fileNames any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBulk", reflect.TypeOf((*MockFiles)(nil).CreateBulk), ctx, jobID, fileNames)
}

// GetById mocks base method.
func (m *MockFiles) GetById(ctx context.Context, id int) (*models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockFilesMockRecorder) GetById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockFiles)(nil).GetById), ctx, id)
}

// GetWithLatestOperationsByJobID mocks base method.
func (m *MockFiles) GetWithLatestOperationsByJobID(jobID int) ([]*models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithLatestOperationsByJobID", jobID)
	ret0, _ := ret[0].([]*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithLatestOperationsByJobID indicates an expected call of GetWithLatestOperationsByJobID.
func (mr *MockFilesMockRecorder) GetWithLatestOperationsByJobID(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithLatestOperationsByJobID", reflect.TypeOf((*MockFiles)(nil).GetWithLatestOperationsByJobID), jobID)
}

// MockOperations is a mock of Operations interface.
type MockOperations struct {
	ctrl     *gomock.Controller
	recorder *MockOperationsMockRecorder
}

// MockOperationsMockRecorder is the mock recorder for MockOperations.
type MockOperationsMockRecorder struct {
	mock *MockOperations
}

// NewMockOperations creates a new mock instance.
func NewMockOperations(ctrl *gomock.Controller) *MockOperations {
	mock := &MockOperations{ctrl: ctrl}
	mock.recorder = &MockOperationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOperations) EXPECT() *MockOperationsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOperations) Create(ctx context.Context, o models.Operation) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, o)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockOperationsMockRecorder) Create(ctx, o any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOperations)(nil).Create), ctx, o)
}

// GetByParams mocks base method.
func (m *MockOperations) GetByParams(ctx context.Context, o models.Operation) (*models.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByParams", ctx, o)
	ret0, _ := ret[0].(*models.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByParams indicates an expected call of GetByParams.
func (mr *MockOperationsMockRecorder) GetByParams(ctx, o any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByParams", reflect.TypeOf((*MockOperations)(nil).GetByParams), ctx, o)
}

// GetLatestOperation mocks base method.
func (m *MockOperations) GetLatestOperation(ctx context.Context, jobID, fileID string) (*models.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestOperation", ctx, jobID, fileID)
	ret0, _ := ret[0].(*models.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestOperation indicates an expected call of GetLatestOperation.
func (mr *MockOperationsMockRecorder) GetLatestOperation(ctx, jobID, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestOperation", reflect.TypeOf((*MockOperations)(nil).GetLatestOperation), ctx, jobID, fileID)
}
