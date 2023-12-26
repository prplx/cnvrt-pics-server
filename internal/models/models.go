package models

import "time"

type File struct {
	ID              int
	JobID           int
	Name            string
	CreatedAt       time.Time
	LatestOperation *Operation
}

type Job struct {
	ID        int
	CreatedAt time.Time
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
	CreatedAt time.Time
}

type PlannedFlush struct {
	ID        int
	JobID     int
	FlushAt   time.Time
	CreatedAt time.Time
}
