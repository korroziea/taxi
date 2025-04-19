package config

import "fmt"

type Config struct {
	AMQP AMQP
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
