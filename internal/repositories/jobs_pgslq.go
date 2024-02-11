package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/cnvrt/internal/models"
)

type JobsRepo struct {
	pool *pgxpool.Pool
}

func (r *JobsRepo) Create(ctx context.Context, session string) (int64, error) {
	var jobID int64
	args := pgx.NamedArgs{
		"session": session,
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO jobs (created_at, session) VALUES (NOW(), @session) RETURNING id;
	`, args).Scan(&jobID)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func (r *JobsRepo) GetByID(ctx context.Context, id int64) (*models.Job, error) {
	var job models.Job
	err := r.pool.QueryRow(ctx, `
		SELECT id, created_at, session FROM jobs WHERE id = $1;
	`, id).Scan(&job.ID, &job.CreatedAt, &job.Session)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobsRepo) Delete(ctx context.Context, jobID int64) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM jobs WHERE id = $1;
	`, jobID)
	if err != nil {
		return err
	}

	return nil
}

func NewJobsRepository(pool *pgxpool.Pool) *JobsRepo {
	return &JobsRepo{pool}
}
