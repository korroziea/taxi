package config

import "fmt"

type Config struct {
	HTTPPort string `env:"HTTP_PORT, default=:3000"`

	Postgres Postgres
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST, default=localhost"`
	Port     int    `env:"POSTGRES_PORT, default=5432"`
	Database string `env:"POSTGRES_DATABASE, default=taxi-user"`
	User     string `env:"POSTGRES_USER, default=postgres"`
	Password string `env:"POSTGRES_PASSWORD, default=secret"`
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
