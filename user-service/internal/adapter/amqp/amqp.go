package amqp

import (
	"fmt"
	"log"

	"github.com/korroziea/taxi/user-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(cfg config.AMQP) (*amqp.Connection, *amqp.Channel, func(), error) {
	conn, err := amqp.Dial(cfg.AMQPURL())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("amqp.Dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("conn.Channel: %w", err)
	}

	amqpDeferFunc := func() {
		if err := ch.Close(); err != nil {
			log.Fatal(err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}

	return conn, ch, amqpDeferFunc, nil
}
