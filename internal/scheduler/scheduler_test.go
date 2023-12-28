package scheduler_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/prplx/cnvrt/internal/config"
	"github.com/prplx/cnvrt/internal/mocks"
	"github.com/prplx/cnvrt/internal/models"
	sch "github.com/prplx/cnvrt/internal/scheduler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestScheduler_ScheduleFlush__should_reset_timer_when_already_added(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	jobsRepo := mocks.NewMockJobs(ctrl)
	plannedFlushesRepo := mocks.NewMockPlannedFlushes(ctrl)
	plannedFlushesRepo.EXPECT().GetAll(gomock.Any()).Times(1).Return(nil, nil)
	scheduler := sch.NewScheduler(config.TestConfig(), logger, communicator, jobsRepo, plannedFlushesRepo)
	jobID := jobID()

	callFirst := logger.EXPECT().PrintInfo("Start flushing job", gomock.Any()).Times(1)
	logger.EXPECT().PrintInfo("Sucessfully flushed the job", gomock.Any()).Times(1).After(callFirst)
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(1)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1)
	plannedFlushesRepo.EXPECT().Create(gomock.Any(), jobID, gomock.Any()).Times(1)

	err := scheduler.ScheduleFlush(jobID, 100*time.Millisecond)
	assert.NoError(t, err)
	err = scheduler.ScheduleFlush(jobID, 10*time.Millisecond)
	assert.NoError(t, err)
	scheduler.ScheduleFlush(jobID, 10*time.Millisecond)
	scheduler.ScheduleFlush(jobID, 10*time.Millisecond)
	scheduler.ScheduleFlush(jobID, 10*time.Millisecond)
	scheduler.ScheduleFlush(jobID, 10*time.Millisecond)
	scheduler.ScheduleFlush(jobID, 10*time.Millisecond)
	time.Sleep(1 * time.Second)
}

func TestScheduler_ScheduleFlush__should_report_twice_when_called_with_different_job_ids(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	jobsRepo := mocks.NewMockJobs(ctrl)
	plannedFlushesRepo := mocks.NewMockPlannedFlushes(ctrl)
	plannedFlushesRepo.EXPECT().GetAll(gomock.Any()).Times(1).Return(nil, nil)
	scheduler := sch.NewScheduler(config.TestConfig(), logger, communicator, jobsRepo, plannedFlushesRepo)
	jobID1 := jobID()
	jobID2 := jobID()

	logger.EXPECT().PrintInfo("Start flushing job", gomock.Any()).Times(2)
	logger.EXPECT().PrintInfo("Sucessfully flushed the job", gomock.Any()).Times(2)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(2)
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(2)
	plannedFlushesRepo.EXPECT().Create(gomock.Any(), jobID1, gomock.Any()).Times(1)
	plannedFlushesRepo.EXPECT().Create(gomock.Any(), jobID2, gomock.Any()).Times(1)

	err := scheduler.ScheduleFlush(jobID1, 100*time.Millisecond)
	assert.NoError(t, err)
	err = scheduler.ScheduleFlush(jobID2, 100*time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(300 * time.Millisecond)
}

func TestScheduler_ScheduleFlush__should_remove_the_job_from_the_map_after_executing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	jobsRepo := mocks.NewMockJobs(ctrl)
	plannedFlushesRepo := mocks.NewMockPlannedFlushes(ctrl)
	plannedFlushesRepo.EXPECT().GetAll(gomock.Any()).Times(1).Return(nil, nil)
	scheduler := sch.NewScheduler(config.TestConfig(), logger, communicator, jobsRepo, plannedFlushesRepo)
	jobID := jobID()

	logger.EXPECT().PrintInfo("Start flushing job", gomock.Any()).Times(1)
	logger.EXPECT().PrintInfo("Sucessfully flushed the job", gomock.Any()).Times(1)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1)
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(1)
	plannedFlushesRepo.EXPECT().Create(gomock.Any(), jobID, gomock.Any()).Times(1)

	err := scheduler.ScheduleFlush(jobID, 100*time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(200 * time.Millisecond)
}

func TestScheduler_ScheduleFlush__should_execute_flushes_if_there_any_in_db(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	jobsRepo := mocks.NewMockJobs(ctrl)
	plannedFlushesRepo := mocks.NewMockPlannedFlushes(ctrl)
	plannedFlushesRepo.EXPECT().GetAll(gomock.Any()).Times(1).Return([]*models.PlannedFlush{
		{
			ID:      1,
			JobID:   1,
			FlushAt: time.Now().UTC(),
		},
		{
			ID:      2,
			JobID:   2,
			FlushAt: time.Now().UTC(),
		},
	}, nil)
	logger.EXPECT().PrintInfo(gomock.Any(), gomock.Any()).AnyTimes()
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(2)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(2)
	sch.NewScheduler(config.TestConfig(), logger, communicator, jobsRepo, plannedFlushesRepo)
	time.Sleep(200 * time.Millisecond)
}

func jobID() int {
	return rand.Intn(2e8)
}
