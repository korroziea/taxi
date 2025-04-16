package driver

import "github.com/redis/go-redis/v9"

type Cache struct {
	db *redis.Client
}

func New(db *redis.Client) *Cache {
	cache := &Cache{
		db: db,
	}

	return cache
}
