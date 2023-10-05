package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type FilesRepo struct {
	conn *pgx.Conn
}

func (r *FilesRepo) Create(jobID, name string) error {
	query := `INSERT INTO files (job_id, name) VALUES ($1, $2);`
	_, err := r.conn.Exec(context.Background(), query, jobID, name)
	return err
}

func (r *FilesRepo) CreateBulk(jobID string, names []string) error {
	rows := make([][]interface{}, len(names))
	for idx, name := range names {
		rows[idx] = []interface{}{jobID, name}
	}
	_, err := r.conn.CopyFrom(context.Background(), pgx.Identifier{"files"}, []string{"job_id", "name"}, pgx.CopyFromRows(rows))
	return err
}

func (r *FilesRepo) GetByJobIDAndName(jobID, name string) (*File, error) {
	query := `SELECT id, job_id, name, created_at FROM files WHERE job_id = $1 AND name = $2 LIMIT 1;`
	row := r.conn.QueryRow(context.Background(), query, jobID, name)
	file := &File{}
	err := row.Scan(&file.ID, &file.JobID, &file.Name, &file.CreatedAt)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func NewFilesRepository(conn *pgx.Conn) *FilesRepo {
	return &FilesRepo{conn}
}
