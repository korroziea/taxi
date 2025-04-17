package config

import "fmt"

type Config struct {
	HTTPPort string `env:"HTTP_PORT, required"`

	SecretKey string `env:"SECRET_KEY, required"`

	Hashing Hashing

	Postgres Postgres
	Redis    Redis
}

type Hashing struct {
	Memory      uint32 `env:"HASHING_MEMORY, required"`
	Iterations  uint32 `env:"HASHING_ITERATIONS, required"`
	Parallelism uint8  `env:"HASHING_PARALLELISM, required"`
	SaltLength  uint32 `env:"HASHING_SALT_LEN, required"`
	KeyLength   uint32 `env:"HASHING_KEY_LEN, required"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST, required"`
	Port     int    `env:"POSTGRES_PORT, required"`
	Database string `env:"POSTGRES_DATABASE, required"`
	User     string `env:"POSTGRES_USER, required"`
	Password string `env:"POSTGRES_PASSWORD, required"`
	SSLMode  string `env:"POSTGRES_SSLMODE, default=disable"`
}

func (p Postgres) PostgresURL() string {
	url := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.Database,
		p.SSLMode,
	)

	return url
}

type Redis struct {
	User     string `env:"REDIS_USER, required"`
	Password string `env:"REDIS_PASSWORD, required"`
	Addr     string `env:"REDIS_ADDR, required"`
	DB       int    `env:"REDIS_DB, required"`
}

func (r Redis) RedisURL() string { // todo: use
	url := fmt.Sprintf(
		"redis://%s:%s@%s/%d",
		r.User,
		r.Password,
		r.Addr,
		r.DB,
	)

	return url
}
