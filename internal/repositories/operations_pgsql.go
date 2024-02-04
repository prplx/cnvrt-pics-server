package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/cnvrt/internal/models"
)

type OperationsRepo struct {
	pool *pgxpool.Pool
}

func (r *OperationsRepo) Create(ctx context.Context, o models.Operation) (int64, error) {
	var operationID int64
	query := `INSERT INTO operations (job_id, file_id, format, quality, fileName, width, height) VALUES (@jobID, @fileID, @format, @quality, @fileName, @width, @height) RETURNING id;`
	args := pgx.NamedArgs{
		"jobID":    o.JobID,
		"fileID":   o.FileID,
		"format":   o.Format,
		"quality":  o.Quality,
		"fileName": o.FileName,
		"width":    o.Width,
		"height":   o.Height,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&operationID)
	if err != nil {
		return 0, err
	}

	return operationID, nil
}

func (r *OperationsRepo) GetByParams(ctx context.Context, o models.Operation) (*models.Operation, error) {
	query := `SELECT id, filename FROM operations WHERE job_id = $1 AND file_id = $2 AND format = $3 AND quality = $4 AND width = $5 AND height = $6;`

	err := r.pool.QueryRow(ctx, query, o.JobID, o.FileID, o.Format, o.Quality, o.Width, o.Height).Scan(&o.ID, &o.FileName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &o, nil
}

func (r *OperationsRepo) GetLatestOperation(ctx context.Context, jobID, fileID string) (*models.Operation, error) {
	query := `SELECT id, filename FROM operations WHERE job_id = $1 AND file_id = $2 AND created_at = (SELECT MAX(created_at) FROM operations WHERE job_id = $1 AND file_id = $2);`

	var o models.Operation

	err := r.pool.QueryRow(ctx, query, jobID, fileID).Scan(&o.ID, &o.FileName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &o, nil
}

func NewOperationsRepository(pool *pgxpool.Pool) *OperationsRepo {
	return &OperationsRepo{pool}
}
