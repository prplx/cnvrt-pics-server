package services

import (
	"github.com/prplx/lighter.pics/internal/communicator"
)

type Communicator interface {
	SendStartProcessing(jobID, fileName string) error
	SendErrorProcessing(jobID, fileName string) error
}

type Logger interface {
	PrintInfo(message string, properties map[string]string)
	PrintError(err error, properties map[string]string)
	PrintFatal(err error, properties map[string]string)
	Write(message []byte) (n int, err error)
}

type Services struct {
	Communicator Communicator
	Logger       Logger
}

type Deps struct {
	Logger Logger
}

func NewServices(deps Deps) *Services {
	return &Services{
		Logger:       deps.Logger,
		Communicator: communicator.NewCommunicator(),
	}
}
