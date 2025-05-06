package driver

import (
	"context"
	"encoding/json"

	"github.com/korroziea/taxi/trip-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	acceptTripQueueName = "accept-trip"
)

type DriverService interface {
	AcceptOrder(ctx context.Context, req domain.AcceptOrderReq) error
}

type Consumer struct {
	l             *zap.Logger
	ch            *amqp.Channel
	driverService DriverService
}

func New(
	l *zap.Logger,
	ch *amqp.Channel,
	driverService DriverService,
) *Consumer {
	consumer := &Consumer{
		l:             l,
		ch:            ch,
		driverService: driverService,
	}

	return consumer
}

func (c *Consumer) Consume(ctx context.Context) {
	q, err := c.ch.QueueDeclare(
		acceptTripQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.l.Fatal("ch.QueueDeclare", zap.Error(err))
	}

	msgs, err := c.ch.Consume(
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
			var req acceptTripReq
			if err := json.Unmarshal(m.Body, &req); err != nil {
				c.l.Error("json.Unmarshal: %w", zap.Error(err))

				// todo: ack
			}

			if err := c.driverService.AcceptOrder(ctx, req.toDomain()); err != nil {
				c.l.Error("driverService.AcceptOrder: %w", zap.Error(err))
			}
		}
	}()

	<-forever
}
