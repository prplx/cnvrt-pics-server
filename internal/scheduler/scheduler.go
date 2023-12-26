package scheduler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/prplx/cnvrt/internal/helpers"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/services"
	"github.com/prplx/cnvrt/internal/types"
)

type Scheduler struct {
	mu                       sync.Mutex
	timers                   map[int]*time.Timer
	config                   *types.Config
	logger                   services.Logger
	communitator             services.Communicator
	jobsRepository           repositories.Jobs
	plannedFlushesRepository repositories.PlannedFlushes
}

func NewScheduler(config *types.Config, l services.Logger, c services.Communicator, jr repositories.Jobs, fr repositories.PlannedFlushes) *Scheduler {
	scheduler := &Scheduler{
		timers:                   make(map[int]*time.Timer),
		config:                   config,
		logger:                   l,
		communitator:             c,
		jobsRepository:           jr,
		plannedFlushesRepository: fr,
	}
	err := scheduler.schedulePlannedFlushes()
	if err != nil {
		l.PrintError(errors.New("error while scheduling planned flushes"), types.AnyMap{
			"error": err,
		})
	}

	return scheduler
}

func (s *Scheduler) ScheduleFlush(jobID int, timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.timers[jobID]; exists {
		s.timers[jobID].Reset(timeout)
		return nil
	}

	timeNowUTC := time.Now()
	flushAt := timeNowUTC.Add(timeout).UTC()

	_, err := s.plannedFlushesRepository.Create(context.Background(), jobID, flushAt)
	if err != nil {
		return err
	}

	return s.scheduleFlush(jobID, timeout)
}

func (s *Scheduler) scheduleFlush(jobID int, timeout time.Duration) error {
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

func (s *Scheduler) schedulePlannedFlushes() error {
	plannedFlushes, err := s.plannedFlushesRepository.GetAll(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(plannedFlushes) == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for _, flush := range plannedFlushes {
		diff := flush.FlushAt.Sub(now)
		s.scheduleFlush(flush.JobID, diff)
	}

	return nil
}
