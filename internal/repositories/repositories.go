package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/lighter.pics/internal/models"
)

type Jobs interface {
	Create(ctx context.Context) (int, error)
}

type Files interface {
	GetById(ctx context.Context, id int) (*models.File, error)
	CreateBulk(ctx context.Context, jobID int, fileNames []string) ([]int, error)
}

type Operations interface {
	Create(ctx context.Context, o models.Operation) (int, error)
	GetByParams(ctx context.Context, o models.Operation) (*models.Operation, error)
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
