package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korroziea/taxi/user-service/internal/config"
)

const connectionTimeout = 3 * time.Second

func Connect(cfg config.Postgres) (*pgxpool.Pool, func(), error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.PostgresURL())
	if err != nil {
		return nil, nil, fmt.Errorf("can't parse postgres config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("can't ping database: %w", err)
	}

	return pool, pool.Close, nil
}
