package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type JobsRepo struct {
	conn *pgx.Conn
}

func (r *JobsRepo) Create() (string, error) {
	var jobID string

	err := r.conn.QueryRow(context.Background(), `
		INSERT INTO jobs (created_at) VALUES (NOW()) RETURNING id;
	`).Scan(&jobID)
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func NewJobsRepository(conn *pgx.Conn) *JobsRepo {
	return &JobsRepo{conn}
}
