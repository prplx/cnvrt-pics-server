package scheduler

import (
	"sync"
	"time"

	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
	"github.com/prplx/lighter.pics/internal/types"
)

type Scheduler struct {
	mu           sync.Mutex
	timers       map[int]*time.Timer
	repositories *repositories.Repositories
	config       *types.Config
	logger       services.Logger
}

func NewScheduler(c *types.Config, r *repositories.Repositories, l services.Logger) *Scheduler {
	return &Scheduler{
		timers:       make(map[int]*time.Timer),
		repositories: r,
		config:       c,
	}
}

func (s *Scheduler) Schedule(jobID int, timeout time.Duration) error {
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

		// DB stuf goes here

		delete(s.timers, jobID)
	})

	return nil
}
