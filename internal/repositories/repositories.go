package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/lighter.pics/internal/models"
)

type Jobs interface {
	Create() (int, error)
}

type Files interface {
	GetById(id int) (*models.File, error)
	CreateBulk(jobID int, fileNames []string) ([]int, error)
}

type Operations interface {
	Create(operation models.Operation) error
}

type Repositories struct {
	Jobs
	Files
	Operations
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Jobs:       NewJobsRepository(pool),
		Files:      NewFilesRepository(pool),
		Operations: NewOperationsRepository(pool),
	}
}
