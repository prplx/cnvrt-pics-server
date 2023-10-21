package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobsRepo struct {
	pool *pgxpool.Pool
}

func (r *JobsRepo) Create(ctx context.Context) (int, error) {
	var jobID int

	err := r.pool.QueryRow(ctx, `
		INSERT INTO jobs (created_at) VALUES (NOW()) RETURNING id;
	`).Scan(&jobID)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func NewJobsRepository(pool *pgxpool.Pool) *JobsRepo {
	return &JobsRepo{pool}
}
