package archiver

import (
	"github.com/prplx/lighter.pics/internal/repositories"
	"github.com/prplx/lighter.pics/internal/services"
)

type Archiver struct {
	repositories *repositories.Repositories
	logger       services.Logger
	communicator services.Communicator
}

func NewArchiver(r *repositories.Repositories, l services.Logger, c services.Communicator) *Archiver {
	return &Archiver{
		repositories: r,
		logger:       l,
		communicator: c,
	}
}

func (a *Archiver) Archive(jobID int) error {
	return nil
}
