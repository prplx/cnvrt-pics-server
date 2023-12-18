package scheduler

import (
	"math/rand"
	"testing"
	"time"

	"github.com/prplx/lighter.pics/internal/config"
	"github.com/prplx/lighter.pics/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestScheduler_ScheduleFlush__should_reset_timer_when_already_added(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMockLogger(ctrl)
	communicator := mocks.NewMockCommunicator(ctrl)
	jobsRepo := mocks.NewMockJobs(ctrl)
	scheduler := NewScheduler(config.TestConfig(), logger, communicator, jobsRepo)

	callFirst := logger.EXPECT().PrintInfo("Start flushing job", gomock.Any()).Times(1)
	logger.EXPECT().PrintInfo("Sucessfully flushed the job", gomock.Any()).Times(1).After(callFirst)
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(1)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1)

	jobID := jobID()
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
	scheduler := NewScheduler(config.TestConfig(), logger, communicator, jobsRepo)

	logger.EXPECT().PrintInfo("Start flushing job", gomock.Any()).Times(2)
	logger.EXPECT().PrintInfo("Sucessfully flushed the job", gomock.Any()).Times(2)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(2)
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(2)

	jobID1 := jobID()
	jobID2 := jobID()
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
	scheduler := NewScheduler(config.TestConfig(), logger, communicator, jobsRepo)

	logger.EXPECT().PrintInfo("Start flushing job", gomock.Any()).Times(1)
	logger.EXPECT().PrintInfo("Sucessfully flushed the job", gomock.Any()).Times(1)
	jobsRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1)
	communicator.EXPECT().SendSuccessFlushing(gomock.Any()).Times(1)

	jobID := jobID()
	err := scheduler.ScheduleFlush(jobID, 100*time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(200 * time.Millisecond)
	assert.Nil(t, scheduler.timers[jobID])
}

func jobID() int {
	return rand.Intn(2e8)
}
