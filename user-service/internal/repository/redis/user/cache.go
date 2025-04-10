package user

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

func (c *Cache) SetToken(ctx context.Context, userID, token string) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	if err := c.db.Set(ctx, userID, token, ttl).Err(); err != nil {
		return fmt.Errorf("db.Set: %w", err)
	}

	return nil
}

func (c *Cache) GetToken(ctx context.Context, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	val, err := c.db.Get(ctx, userID).Result()
	if err != nil {
		return "", fmt.Errorf("db.Get: %w", err)
	}

	return val, nil
}
