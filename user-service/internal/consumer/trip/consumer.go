package trip

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/korroziea/taxi/user-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	startTripRespQueueName = "start-trip-resp"
)

type Handler interface {
	TripsTemp(trips []domain.Trip) gin.HandlerFunc
}

type Consumer struct {
	l  *zap.Logger
	ch *amqp.Channel
}

func New(
	l *zap.Logger,
	ch *amqp.Channel,
) *Consumer {
	consumer := &Consumer{
		l:  l,
		ch: ch,
	}

	return consumer
}

func (c *Consumer) ConsumeStartTrip(ctx context.Context) {
	q, err := c.ch.QueueDeclare(
		startTripRespQueueName,
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
			fmt.Println(string(m.Body))
		}
	}()

	<-forever
}
