package consumer

import (
	"fmt"
	"log"

	"github.com/korroziea/taxi/driver-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(cfg config.AMQP) (*amqp.Connection, func(), error) {
	conn, err := amqp.Dial(cfg.AMQPURL())
	if err != nil {
		return nil, nil, fmt.Errorf("amqp.Dial: %w", err)
	}

	amqpDeferFunc := func() {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}

	return conn, amqpDeferFunc, nil
}
