package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/lighter.pics/internal/models"
)

type OperationsRepo struct {
	pool *pgxpool.Pool
}

func (r *OperationsRepo) Create(operation models.Operation) error {
	query := `INSERT INTO operations (job_id, file_id, format, quality, fileName, width, height) VALUES (@jobID, @fileID, @format, @quality, @fileName, @width, @height);`
	args := pgx.NamedArgs{
		"jobID":    operation.JobID,
		"fileID":   operation.FileID,
		"format":   operation.Format,
		"quality":  operation.Quality,
		"fileName": operation.FileName,
		"width":    operation.Width,
		"height":   operation.Height,
	}

	_, err := r.pool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	return nil
}

func NewOperationsRepository(pool *pgxpool.Pool) *OperationsRepo {
	return &OperationsRepo{pool}
}
