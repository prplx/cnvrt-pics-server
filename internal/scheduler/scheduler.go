package scheduler

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/types"
)

type Scheduler struct {
	mu             sync.Mutex
	timers         map[int]*time.Timer
	config         *types.Config
	logger         services.Logger
	communitator   services.Communicator
	jobsRepository repositories.Jobs
}

func NewScheduler(config *types.Config, l services.Logger, c services.Communicator, jr repositories.Jobs) *Scheduler {
	return &Scheduler{
		timers:         make(map[int]*time.Timer),
		config:         config,
		logger:         l,
		communitator:   c,
		jobsRepository: jr,
	}
}

func (s *Scheduler) ScheduleFlush(jobID int, timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.timers[jobID]; exists {
		s.timers[jobID].Reset(timeout)
		return nil
	}

	s.timers[jobID] = time.AfterFunc(timeout, func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.logger.PrintInfo("Start flushing job", types.AnyMap{
			"job_id": jobID,
		})

		dir := helpers.BuildPath(s.config.Process.UploadDir, jobID)
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			err := os.RemoveAll(dir)
			if err != nil {
				s.logger.PrintInfo("Error while flushing the job", types.AnyMap{
					"job_id": jobID,
				})
			}
		}

		delete(s.timers, jobID)

		s.communitator.SendSuccessFlushing(jobID)
		s.logger.PrintInfo("Sucessfully flushed the job", types.AnyMap{
			"job_id": jobID,
		})

		err := s.jobsRepository.Delete(context.Background(), jobID)
		if err != nil {
			s.logger.PrintError(errors.New("error while deleting the job"), types.AnyMap{
				"job_id": jobID,
			})
		}
	})

	return nil
}
