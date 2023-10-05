package repositories

import (
	"github.com/jackc/pgx/v5"
)

type Jobs interface {
	Create() (string, error)
}

type File struct {
	ID        int
	JobID     int
	Name      string
	CreatedAt string
}

type Files interface {
	Create(jobID, fileName string) error
	CreateBulk(jobID string, fileNames []string) error
	GetByJobIDAndName(jobID, fileName string) (*File, error)
}

type Operations interface {
	Create(jobID string, fileID int, format string, quality int, fileName string, width, height int) (string, error)
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
