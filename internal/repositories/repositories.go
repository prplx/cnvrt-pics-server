package repositories

import (
	"context"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/cnvrt/internal/models"
)

type Jobs interface {
	GetByID(ctx context.Context, id int64) (*models.Job, error)
	Create(ctx context.Context, session string) (int64, error)
	Delete(ctx context.Context, jobID int64) error
}

type Files interface {
	GetWithJobByID(ctx context.Context, id int64) (*models.File, error)
	CreateBulk(ctx context.Context, jobID int64, fileNames []string) ([]models.File, error)
	GetWithLatestOperationsByJobID(ctx context.Context, jobId int64) ([]*models.File, error)
	AddToJob(ctx context.Context, jobID int64, fileName string) (models.File, error)
	GetByJobID(ctx context.Context, jobID int64) ([]models.File, error)
	DeleteFromJob(ctx context.Context, jobID int64, fileID int64) error
}

type Operations interface {
	Create(ctx context.Context, o models.Operation) (int64, error)
	GetByParams(ctx context.Context, o models.Operation) (*models.Operation, error)
	GetLatestOperation(ctx context.Context, jobID, fileID string) (*models.Operation, error)
}

type PlannedFlushes interface {
	Create(ctx context.Context, jobID int64, flushAt time.Time) (int64, error)
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
