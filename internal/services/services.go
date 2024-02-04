package services

import (
	"context"
	"time"

	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/types"
)

type Communicator interface {
	AddClient(jobID int64, connection types.WebsocketConnection)
	RemoveClient(jobID int64)
	SendStartProcessing(jobID, fileID int64, fileName string) error
	SendErrorProcessing(jobID, fileID int64, fileName string) error
	SendSuccessProcessing(jobID int64, result types.SuccessResult) error
	SendStartArchiving(jobID int64) error
	SendErrorArchiving(jobID int64) error
	SendSuccessArchiving(jobID int64, path string) error
	SendSuccessFlushing(jobID int64) error
}

type Logger interface {
	PrintInfo(message string, properties ...types.AnyMap)
	PrintError(err error, properties ...types.AnyMap)
	PrintFatal(err error, properties ...types.AnyMap)
	Write(message []byte) (n int, err error)
}

type Processor interface {
	Startup()
	Process(ctx context.Context, input types.ImageProcessInput)
	Shutdown()
}

type Archiver interface {
	Archive(jobID int64) error
}

type Scheduler interface {
	ScheduleFlush(jobID int64, timeout time.Duration) error
}

type Services struct {
	Communicator Communicator
	Logger       Logger
	Repositories *repositories.Repositories
	Processor    Processor
	Archiver     Archiver
	Config       *types.Config
	Scheduler    Scheduler
}

type Deps struct {
	Logger       Logger
	Repositories *repositories.Repositories
	Processor    Processor
	Communicator Communicator
	Archiver     Archiver
	Config       *types.Config
	Scheduler    Scheduler
}

func NewServices(deps Deps) *Services {
	return &Services{
		Logger:       deps.Logger,
		Communicator: deps.Communicator,
		Repositories: deps.Repositories,
		Processor:    deps.Processor,
		Archiver:     deps.Archiver,
		Config:       deps.Config,
		Scheduler:    deps.Scheduler,
	}
}
