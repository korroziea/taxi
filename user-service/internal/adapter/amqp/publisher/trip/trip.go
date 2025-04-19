package trip

import (
	"context"
	"fmt"
	"time"

	"github.com/korroziea/taxi/user-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

const publishTimeout = 3 * time.Second

type Adapter struct {
	conn *amqp.Connection
}

func New(conn *amqp.Connection) *Adapter {
	adapter := &Adapter{
		conn: conn,
	}

	return adapter
}

func (a *Adapter) StartTrip(ctx context.Context, trip domain.StartTrip) error {
	ch, err := a.conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"start-trip-req",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello"),
		},
	)
	if err != nil {
		return fmt.Errorf("ch.PublishWithContext: %w", err)
	}

	return nil
}

func (a *Adapter) CancelTrip(ctx context.Context) error {
	return nil
}
