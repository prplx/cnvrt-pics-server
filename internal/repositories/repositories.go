package repositories

import (
	"github.com/jackc/pgx/v5"
)

type Jobs interface {
	Create() (int, error)
}

type File struct {
	ID        int
	JobID     int
	Name      string
	CreatedAt string
}

type Files interface {
	CreateBulk(jobID int, fileNames []string) ([]int, error)
}

type Operations interface {
	Create(jobID, fileID int, format string, quality int, fileName string, width, height int) (string, error)
}

type Repositories struct {
	Jobs
	Files
	Operations
}

func NewRepositories(conn *pgx.Conn) *Repositories {
	return &Repositories{
		Jobs:       NewJobsRepository(conn),
		Files:      NewFilesRepository(conn),
		Operations: NewOperationsRepository(conn),
	}
}
