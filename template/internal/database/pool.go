package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"{{ .ModulePath }}/internal/utils/e"
)

func GetPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, e.Wrap(err, "error creating connection pool")
	}
	return pool, nil
}
