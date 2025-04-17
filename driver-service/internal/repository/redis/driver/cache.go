package driver

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ttl          = 3600*time.Second + 1*time.Second
	queryTimeout = 3 * time.Second
)

type Cache struct {
	db *redis.Client
}

func New(db *redis.Client) *Cache {
	cache := &Cache{
		db: db,
	}

	return cache
}

func (c *Cache) SetToken(ctx context.Context, driverID, token string) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	if err := c.db.Set(ctx, driverID, token, ttl).Err(); err != nil {
		return fmt.Errorf("db.Set: %w", err)
	}

	return nil
}
