package services

import (
	"github.com/prplx/lighter.pics/internal/communicator"
	"github.com/prplx/lighter.pics/internal/imageProcessor"
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/types"
)

type Communicator interface {
	SendStartProcessing(jobID, fileID int, fileName string) error
	SendErrorProcessing(jobID, fileID int, fileName string) error
	SendSuccessProcessing(jobID int, result types.SuccessResult) error
}

type Logger interface {
	PrintInfo(message string, properties ...types.AnyMap)
	PrintError(err error, properties ...types.AnyMap)
	PrintFatal(err error, properties ...types.AnyMap)
	Write(message []byte) (n int, err error)
}

type ImageProcessor interface {
	NewImage(buf []byte) *imageProcessor.Image
	Write(path string, buf []byte) error
}

type Services struct {
	Communicator   Communicator
	Logger         Logger
	Repositories   repositories.Repositories
	ImageProcessor ImageProcessor
}

type Deps struct {
	Logger Logger
	repositories.Repositories
	ImageProcessor ImageProcessor
}

func NewServices(deps Deps) *Services {
	return &Services{
		Logger:         deps.Logger,
		Communicator:   communicator.NewCommunicator(),
		Repositories:   deps.Repositories,
		ImageProcessor: deps.ImageProcessor,
	}
}
