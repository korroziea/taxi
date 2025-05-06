package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/korroziea/taxi/trip-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	publishTimeout = 3 * time.Second

	findDriverQueueName = "find-driver"
)

type Adapter struct {
	conn *amqp.Connection
}

func New(conn *amqp.Connection) *Adapter {
	adapter := &Adapter{
		conn: conn,
	}

	return adapter
}

func (a *Adapter) FindDriver(ctx context.Context, req domain.FindDriverReq) error {
	ch, err := a.conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		findDriverQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare: %w", err)
	}

	body, err := json.Marshal(toFindDriverBody(req))
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	fmt.Println("publisher - ", req)

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
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("ch.PublishWithContext: %w", err)
	}

	return nil
}
