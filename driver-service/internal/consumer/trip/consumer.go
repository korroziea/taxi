package trip

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	findDriverQueueName = "find-driver"
)

type TripService interface {
	AcceptOrder(ctx context.Context, tripID string) error
}

type Consumer struct {
	l           *zap.Logger
	conn        *amqp.Connection
	tripServive TripService
}

func New(
	l *zap.Logger,
	conn *amqp.Connection,
	tripServive TripService,
) *Consumer {
	consumer := &Consumer{
		l:           l,
		conn:        conn,
		tripServive: tripServive,
	}

	return consumer
}

func (c *Consumer) Consume(ctx context.Context) {
	ch, err := c.conn.Channel()
	if err != nil {
		c.l.Fatal("conn.Channel", zap.Error(err))
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		findDriverQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.l.Fatal("ch.QueueDeclare", zap.Error(err))
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.l.Fatal("ch.Consume", zap.Error(err))
	}

	var forever chan struct{}

	go func() {
		for m := range msgs {
			var req findDriverReq
			if err := json.Unmarshal(m.Body, &req); err != nil {
				c.l.Info(string(m.Body))
				fmt.Println("consumer - ", string(m.Body))
				c.l.Error("json.Unmarshal: %w", zap.Error(err))

				// todo: ack
			}

			if err := c.tripServive.AcceptOrder(context.Background(), req.UserID); err != nil {
				c.l.Error("tripServive.AcceptOrder: %w", zap.Error(err))
			}
		}
	}()

	<-forever
}
