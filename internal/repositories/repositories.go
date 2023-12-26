package repositories

import (
	"context"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/cnvrt/internal/models"
)

type Jobs interface {
	Create(ctx context.Context) (int, error)
	Delete(ctx context.Context, jobID int) error
}

type Files interface {
	GetByID(ctx context.Context, id int) (*models.File, error)
	CreateBulk(ctx context.Context, jobID int, fileNames []string) ([]models.File, error)
	GetWithLatestOperationsByJobID(jobID int) ([]*models.File, error)
}

type Operations interface {
	Create(ctx context.Context, o models.Operation) (int, error)
	GetByParams(ctx context.Context, o models.Operation) (*models.Operation, error)
	GetLatestOperation(ctx context.Context, jobID, fileID string) (*models.Operation, error)
}

type PlannedFlushes interface {
	Create(ctx context.Context, jobID int, flushAt time.Time) (int, error)
	GetAll(ctx context.Context) ([]*models.PlannedFlush, error)
}

type Repositories struct {
	Jobs
	Files
	Operations
	PlannedFlushes
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Jobs:           NewJobsRepository(pool),
		Files:          NewFilesRepository(pool),
		Operations:     NewOperationsRepository(pool),
		PlannedFlushes: NewPlannedFlushesRepository(pool),
	}
}
