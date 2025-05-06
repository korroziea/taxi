package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/korroziea/taxi/driver-service/internal/domain"
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

func (a *Adapter) AcceptOrder(ctx context.Context, resp domain.AcceptOrderResp) error {
	ch, err := a.conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"accept-trip",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare: %w", err)
	}

	body, err := json.Marshal(toAcceptTripResp(resp))
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	fmt.Println("publisher - ", resp)

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
