package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/lighter.pics/internal/models"
)

type FilesRepo struct {
	pool *pgxpool.Pool
}

func (r *FilesRepo) GetById(ctx context.Context, id int) (*models.File, error) {
	query := `SELECT id, name FROM files WHERE id = @id;`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.pool.QueryRow(context.Background(), query, args)
	file := &models.File{}
	err := row.Scan(&file.ID, &file.Name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *FilesRepo) CreateBulk(ctx context.Context, jobID int, names []string) ([]models.File, error) {
	query := `INSERT INTO files (job_id, name) VALUES (@jobID, @name) RETURNING id, name;`
	batch := &pgx.Batch{}
	files := []models.File{}
	for _, name := range names {
		args := pgx.NamedArgs{
			"jobID": jobID,
			"name":  name,
		}
		batch.Queue(query, args)
	}
	result := r.pool.SendBatch(context.Background(), batch)
	defer result.Close()

	for range names {
		file := models.File{}
		row := result.QueryRow()
		err := row.Scan(&file.ID, &file.Name)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *FilesRepo) GetWithLatestOperationsByJobID(jobID int) ([]*models.File, error) {
	query := `
		SELECT f.id, f.name, o.filename, o.format
		FROM files f
		INNER JOIN operations o ON o.file_id = f.id
		WHERE f.job_id = @jobID
		AND o.created_at = (
			SELECT MAX(created_at)
			FROM operations
			WHERE file_id = f.id
		)
		ORDER BY f.id ASC;
	`
	args := pgx.NamedArgs{
		"jobID": jobID,
	}
	rows, err := r.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []*models.File{}
	for rows.Next() {
		file := &models.File{LatestOperation: &models.Operation{}}
		err := rows.Scan(&file.ID, &file.Name, &file.LatestOperation.FileName, &file.LatestOperation.Format)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func NewFilesRepository(pool *pgxpool.Pool) *FilesRepo {
	return &FilesRepo{pool}
}
