package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/korroziea/taxi/driver-service/internal/config"
	"github.com/redis/go-redis/v9"
)

const connectionTimeout = 3 * time.Second

func Connect(cfg config.Redis) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := db.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	return db, nil
}
