package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type OperationsRepo struct {
	conn *pgx.Conn
}

func (r *OperationsRepo) Create(jobID, fileID int, format string, quality int, fileName string, width, height int) (string, error) {
	var fileId string
	query := `INSERT INTO operations (job_id, file_id, format, quality, fileName, width, height) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	err := r.conn.QueryRow(context.Background(), query, jobID, fileID, format, quality, fileName, width, height).Scan(&fileId)
	if err != nil {
		return "", err
	}

	return fileId, nil
}

func NewOperationsRepository(conn *pgx.Conn) *OperationsRepo {
	return &OperationsRepo{conn}
}
