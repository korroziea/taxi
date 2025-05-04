package bootsrap

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/config"
	"github.com/korroziea/taxi/trip-service/internal/consumer"
	userconsumer "github.com/korroziea/taxi/trip-service/internal/consumer/user"
	"github.com/korroziea/taxi/trip-service/internal/repository/psql"
	triprepo "github.com/korroziea/taxi/trip-service/internal/repository/psql/trip"
	tripservice "github.com/korroziea/taxi/trip-service/internal/service/user"
	"go.uber.org/zap"
)

type App struct {
	l            *zap.Logger
	userConsumer *userconsumer.Consumer
}

func New(l *zap.Logger, cfg config.Config) (*App, error) {
	postgresDB, _, err := psql.Connect(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("psql.Connect: %w", err)
	}

	amqpConn, _, err := consumer.Connect(cfg.AMQP)
	if err != nil {
		return nil, fmt.Errorf("consumer.Connect: %w", err)
	}

	tripRepo := triprepo.New(postgresDB)

	tripService := tripservice.New(tripRepo)

	userConsumer := userconsumer.New(l, amqpConn, tripService)

	app := &App{
		l:            l,
		userConsumer: userConsumer,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) {
	go func() {
		a.userConsumer.Consume(ctx)
	}()

	<-ctx.Done()
}
