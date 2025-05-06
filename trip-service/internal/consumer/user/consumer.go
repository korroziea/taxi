package user

import (
	"context"
	"encoding/json"

	"github.com/korroziea/taxi/trip-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	startTripQueueName = "start-trip"
)

type TripService interface {
	StartTrip(ctx context.Context, trip domain.StartTrip) error
}

type Consumer struct {
	l           *zap.Logger
	ch          *amqp.Channel
	tripServive TripService
}

func New(
	l *zap.Logger,
	ch *amqp.Channel,
	tripServive TripService,
) *Consumer {
	consumer := &Consumer{
		l:           l,
		ch:          ch,
		tripServive: tripServive,
	}

	return consumer
}

func (c *Consumer) Consume(ctx context.Context) {
	q, err := c.ch.QueueDeclare(
		startTripQueueName,
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
			var trip startTrip
			if err := json.Unmarshal(m.Body, &trip); err != nil {
				c.l.Error("json.Unmarshal: %w", zap.Error(err))

				// todo: ack
			}

			if err := c.tripServive.StartTrip(context.Background(), trip.toDomain()); err != nil {
				c.l.Error("tripServive.StartTrip: %w", zap.Error(err))
			}
		}
	}()

	<-forever
}
