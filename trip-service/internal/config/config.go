package config

import "fmt"

type Config struct {
	HTTPPort string `env:"HTTP_PORT, required"`
	
	Postgres Postgres

	AMQP AMQP
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

type AMQP struct {
	User     string `env:"AMQP_USER, required"`
	Password string `env:"AMQP_PASSWORD, required"`
	Host     string `env:"AMQP_HOST, required"`
	Port     string `env:"AMQP_PORT, required"`
}

func (a AMQP) AMQPURL() string {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		a.User,
		a.Password,
		a.Host,
		a.Port,
	)

	return url
}
