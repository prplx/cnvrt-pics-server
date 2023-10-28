package pg

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	Pool *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) *postgres {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatal("unable to create connection pool: %w", err)
		}

		pgInstance = &postgres{db}
	})

	return pgInstance
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.Pool.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.Pool.Close()
}
