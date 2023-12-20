package services

import (
	"context"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/prplx/cnvrt/internal/repositories"
	"github.com/prplx/cnvrt/internal/types"
)

type Communicator interface {
	AddClient(jobID int, connection *websocket.Conn)
	RemoveClient(jobID int)
	SendStartProcessing(jobID, fileID int, fileName string) error
	SendErrorProcessing(jobID, fileID int, fileName string) error
	SendSuccessProcessing(jobID int, result types.SuccessResult) error
	SendStartArchiving(jobID int) error
	SendErrorArchiving(jobID int) error
	SendSuccessArchiving(jobID int, path string) error
	SendSuccessFlushing(jobID int) error
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
	Archive(jobID int) error
}

type Scheduler interface {
	ScheduleFlush(jobID int, timeout time.Duration) error
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
