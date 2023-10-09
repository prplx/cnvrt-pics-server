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

func (r *FilesRepo) GetById(id int) (*models.File, error) {
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

func (r *FilesRepo) CreateBulk(jobID int, names []string) ([]int, error) {
	query := `INSERT INTO files (job_id, name) VALUES (@jobID, @name) RETURNING id;`
	batch := &pgx.Batch{}
	ids := []int{}
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
		var id int
		row := result.QueryRow()
		err := row.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func NewFilesRepository(pool *pgxpool.Pool) *FilesRepo {
	return &FilesRepo{pool}
}
