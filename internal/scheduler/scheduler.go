package scheduler

import (
	"os"
	"sync"
	"time"

	"github.com/prplx/lighter.pics/internal/helpers"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
)

type Scheduler struct {
	mu           sync.Mutex
	timers       map[int]*time.Timer
	config       *types.Config
	logger       services.Logger
	communitator services.Communicator
}

func NewScheduler(config *types.Config, l services.Logger, c services.Communicator) *Scheduler {
	return &Scheduler{
		timers:       make(map[int]*time.Timer),
		config:       config,
		logger:       l,
		communitator: c,
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
	})

	return nil
}
