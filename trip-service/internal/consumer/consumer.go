package consumer

import (
	"github.com/korroziea/taxi/trip-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type AMQPServer struct {
	l   *zap.Logger
	cfg config.AMQP
}

func New(l *zap.Logger, cfg config.AMQP) *AMQPServer {
	server := &AMQPServer{
		l:   l,
		cfg: cfg,
	}

	return server
}

func (s *AMQPServer) Consume() {
	conn, err := amqp.Dial(s.cfg.AMQPURL())
	if err != nil {
		s.l.Panic("amqp.Dial", zap.Error(err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		s.l.Panic("conn.Channel", zap.Error(err))
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"start-trip-req",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.l.Panic("ch.QueueDeclare", zap.Error(err))
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
		s.l.Panic("ch.Consume", zap.Error(err))
	}

	var forever chan struct{}

	go func() {
		for m := range msgs {
			s.l.Info("message", zap.String(m.Timestamp.String(), string(m.Body)))
		}
	}()

	<-forever
}
