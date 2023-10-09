package models

type File struct {
	ID        int
	JobID     int
	Name      string
	CreatedAt string
}

type Job struct {
	ID        int
	CreatedAt string
}

type Operation struct {
	ID        int
	JobID     int
	FileID    int
	Format    string
	Quality   int
	FileName  string
	Width     int
	Height    int
	CreatedAt string
}
