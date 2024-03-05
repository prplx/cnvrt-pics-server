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
	time "time"

	models "github.com/prplx/cnvrt/internal/models"
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
func (m *MockJobs) Create(ctx context.Context, session string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, session)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockJobsMockRecorder) Create(ctx, session any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockJobs)(nil).Create), ctx, session)
}

// Delete mocks base method.
func (m *MockJobs) Delete(ctx context.Context, jobID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockJobsMockRecorder) Delete(ctx, jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockJobs)(nil).Delete), ctx, jobID)
}

// GetByID mocks base method.
func (m *MockJobs) GetByID(ctx context.Context, id int64) (*models.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockJobsMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockJobs)(nil).GetByID), ctx, id)
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

// AddToJob mocks base method.
func (m *MockFiles) AddToJob(ctx context.Context, jobID int64, fileName string) (models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToJob", ctx, jobID, fileName)
	ret0, _ := ret[0].(models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddToJob indicates an expected call of AddToJob.
func (mr *MockFilesMockRecorder) AddToJob(ctx, jobID, fileName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToJob", reflect.TypeOf((*MockFiles)(nil).AddToJob), ctx, jobID, fileName)
}

// CreateBulk mocks base method.
func (m *MockFiles) CreateBulk(ctx context.Context, jobID int64, fileNames []string) ([]models.File, error) {
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

// DeleteFromJob mocks base method.
func (m *MockFiles) DeleteFromJob(ctx context.Context, jobID, fileID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFromJob", ctx, jobID, fileID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFromJob indicates an expected call of DeleteFromJob.
func (mr *MockFilesMockRecorder) DeleteFromJob(ctx, jobID, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFromJob", reflect.TypeOf((*MockFiles)(nil).DeleteFromJob), ctx, jobID, fileID)
}

// GetByJobID mocks base method.
func (m *MockFiles) GetByJobID(ctx context.Context, jobID int64) ([]models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByJobID", ctx, jobID)
	ret0, _ := ret[0].([]models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByJobID indicates an expected call of GetByJobID.
func (mr *MockFilesMockRecorder) GetByJobID(ctx, jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByJobID", reflect.TypeOf((*MockFiles)(nil).GetByJobID), ctx, jobID)
}

// GetWithJobByID mocks base method.
func (m *MockFiles) GetWithJobByID(ctx context.Context, id int64) (*models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithJobByID", ctx, id)
	ret0, _ := ret[0].(*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithJobByID indicates an expected call of GetWithJobByID.
func (mr *MockFilesMockRecorder) GetWithJobByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithJobByID", reflect.TypeOf((*MockFiles)(nil).GetWithJobByID), ctx, id)
}

// GetWithLatestOperationsByJobID mocks base method.
func (m *MockFiles) GetWithLatestOperationsByJobID(ctx context.Context, jobId int64) ([]*models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithLatestOperationsByJobID", ctx, jobId)
	ret0, _ := ret[0].([]*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithLatestOperationsByJobID indicates an expected call of GetWithLatestOperationsByJobID.
func (mr *MockFilesMockRecorder) GetWithLatestOperationsByJobID(ctx, jobId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithLatestOperationsByJobID", reflect.TypeOf((*MockFiles)(nil).GetWithLatestOperationsByJobID), ctx, jobId)
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
func (m *MockOperations) Create(ctx context.Context, o models.Operation) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, o)
	ret0, _ := ret[0].(int64)
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

// MockPlannedFlushes is a mock of PlannedFlushes interface.
type MockPlannedFlushes struct {
	ctrl     *gomock.Controller
	recorder *MockPlannedFlushesMockRecorder
}

// MockPlannedFlushesMockRecorder is the mock recorder for MockPlannedFlushes.
type MockPlannedFlushesMockRecorder struct {
	mock *MockPlannedFlushes
}

// NewMockPlannedFlushes creates a new mock instance.
func NewMockPlannedFlushes(ctrl *gomock.Controller) *MockPlannedFlushes {
	mock := &MockPlannedFlushes{ctrl: ctrl}
	mock.recorder = &MockPlannedFlushesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlannedFlushes) EXPECT() *MockPlannedFlushesMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPlannedFlushes) Create(ctx context.Context, jobID int64, flushAt time.Time) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, jobID, flushAt)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPlannedFlushesMockRecorder) Create(ctx, jobID, flushAt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPlannedFlushes)(nil).Create), ctx, jobID, flushAt)
}

// GetAll mocks base method.
func (m *MockPlannedFlushes) GetAll(ctx context.Context) ([]*models.PlannedFlush, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*models.PlannedFlush)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPlannedFlushesMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPlannedFlushes)(nil).GetAll), ctx)
}
