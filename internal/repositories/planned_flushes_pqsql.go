package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/cnvrt/internal/models"
)

type PlannedFlushesRepo struct {
	pool *pgxpool.Pool
}

func (r *PlannedFlushesRepo) Create(ctx context.Context, jobID int64, flushAt time.Time) (int64, error) {
	var flushID int64
	query := `INSERT INTO planned_flushes (job_id, flush_at) VALUES (@jobID, @flushAt) RETURNING id;`
	args := pgx.NamedArgs{
		"jobID":   jobID,
		"flushAt": flushAt,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&flushID)
	if err != nil {
		return 0, err
	}

	return flushID, nil
}

func (r *PlannedFlushesRepo) GetAll(ctx context.Context) ([]*models.PlannedFlush, error) {
	query := `SELECT id, job_id, flush_at, created_at FROM planned_flushes;`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flushes []*models.PlannedFlush

	for rows.Next() {
		var flush models.PlannedFlush

		err = rows.Scan(&flush.ID, &flush.JobID, &flush.FlushAt, &flush.CreatedAt)
		if err != nil {
			return nil, err
		}

		flushes = append(flushes, &flush)
	}

	return flushes, nil
}

func NewPlannedFlushesRepository(pool *pgxpool.Pool) *PlannedFlushesRepo {
	return &PlannedFlushesRepo{pool}
}
