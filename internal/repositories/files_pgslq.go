package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type FilesRepo struct {
	conn *pgx.Conn
}

func (r *FilesRepo) CreateBulk(jobID int, names []string) ([]int, error) {
	query := `INSERT INTO FILES (job_id, name) VALUES (@jobID, @name) RETURNING id;`
	batch := &pgx.Batch{}
	ids := []int{}
	for _, name := range names {
		args := pgx.NamedArgs{
			"jobID": jobID,
			"name":  name,
		}
		batch.Queue(query, args)
	}
	result := r.conn.SendBatch(context.Background(), batch)
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

func NewFilesRepository(conn *pgx.Conn) *FilesRepo {
	return &FilesRepo{conn}
}
