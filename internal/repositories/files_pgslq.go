package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prplx/cnvrt/internal/models"
)

type FilesRepo struct {
	pool *pgxpool.Pool
}

func NewFilesRepository(pool *pgxpool.Pool) *FilesRepo {
	return &FilesRepo{pool}
}

func (r *FilesRepo) GetWithJobByID(ctx context.Context, id int64) (*models.File, error) {
	query := `SELECT f.id, f.name, j.id, j.session FROM files f INNER JOIN jobs j ON f.job_id = j.id WHERE f.id = @id;`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.pool.QueryRow(context.Background(), query, args)
	file := &models.File{Job: &models.Job{}}
	err := row.Scan(&file.ID, &file.Name, &file.Job.ID, &file.Job.Session)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *FilesRepo) CreateBulk(ctx context.Context, jobID int64, names []string) ([]models.File, error) {
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

func (r *FilesRepo) GetWithLatestOperationsByJobID(ctx context.Context, jobID int64) ([]*models.File, error) {
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
	rows, err := r.pool.Query(ctx, query, args)
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

func (r *FilesRepo) AddToJob(ctx context.Context, jobID int64, fileName string) (models.File, error) {
	query := `INSERT INTO files (job_id, name) VALUES (@jobID, @fileName) RETURNING id, name;`
	args := pgx.NamedArgs{
		"jobID":    jobID,
		"fileName": fileName,
	}
	row := r.pool.QueryRow(context.Background(), query, args)
	file := models.File{}
	err := row.Scan(&file.ID, &file.Name)
	if err != nil {
		return models.File{}, err
	}
	return file, nil
}

func (r *FilesRepo) GetByJobID(ctx context.Context, jobID int64) ([]models.File, error) {
	query := `SELECT id, name FROM files WHERE job_id = @jobID;`
	args := pgx.NamedArgs{
		"jobID": jobID,
	}
	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []models.File{}
	for rows.Next() {
		file := models.File{}
		err := rows.Scan(&file.ID, &file.Name)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *FilesRepo) DeleteFromJob(ctx context.Context, jobID int64, fileID int64) error {
	query := `DELETE FROM files WHERE job_id = @jobID AND id = @fileID;`
	args := pgx.NamedArgs{
		"jobID":  jobID,
		"fileID": fileID,
	}
	_, err := r.pool.Exec(ctx, query, args)
	return err
}
