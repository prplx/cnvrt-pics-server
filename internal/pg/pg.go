package pg

import (
	"context"
	"log"
	"sync"
	"time"

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
	ctxwt, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return pg.Pool.Ping(ctxwt)
}

func (pg *postgres) Close() {
	pg.Pool.Close()
}
