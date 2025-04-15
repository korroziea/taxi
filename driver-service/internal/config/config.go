package config

type Config struct {
	HTTPPort string `env:"HTTP_PORT, required"`
}
