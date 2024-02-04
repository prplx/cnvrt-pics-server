package models

import "time"

type File struct {
	ID              int64
	JobID           int64
	Name            string
	CreatedAt       time.Time
	LatestOperation *Operation
}

type Job struct {
	ID        int64
	CreatedAt time.Time
}

type Operation struct {
	ID        int64
	JobID     int64
	FileID    int64
	Format    string
	Quality   int
	FileName  string
	Width     int
	Height    int
	CreatedAt time.Time
}

type PlannedFlush struct {
	ID        int64
	JobID     int64
	FlushAt   time.Time
	CreatedAt time.Time
}
