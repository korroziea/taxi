package bootsrap

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/adapter/amqp"
	driverpublisher "github.com/korroziea/taxi/trip-service/internal/adapter/amqp/publisher/driver"
	userpublisher "github.com/korroziea/taxi/trip-service/internal/adapter/amqp/publisher/user"
	"github.com/korroziea/taxi/trip-service/internal/config"
	driverconsumer "github.com/korroziea/taxi/trip-service/internal/consumer/driver"
	userconsumer "github.com/korroziea/taxi/trip-service/internal/consumer/user"
	"github.com/korroziea/taxi/trip-service/internal/repository/psql"
	triprepo "github.com/korroziea/taxi/trip-service/internal/repository/psql/trip"
	driverservice "github.com/korroziea/taxi/trip-service/internal/service/driver"
	userservice "github.com/korroziea/taxi/trip-service/internal/service/user"
	"go.uber.org/zap"
)

type App struct {
	l              *zap.Logger
	userConsumer   *userconsumer.Consumer
	driverConsumer *driverconsumer.Consumer
}

func New(l *zap.Logger, cfg config.Config) (*App, error) {
	postgresDB, _, err := psql.Connect(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("psql.Connect: %w", err)
	}

	amqpConn, amqpCh, _, err := amqp.Connect(cfg.AMQP)
	if err != nil {
		return nil, fmt.Errorf("amqp.Connect: %w", err)
	}

	tripRepo := triprepo.New(postgresDB)

	driverPublisher := driverpublisher.New(amqpConn)
	userPublisher := userpublisher.New(amqpConn)

	userService := userservice.New(tripRepo, driverPublisher)
	driverService := driverservice.New(tripRepo, userPublisher)

	userConsumer := userconsumer.New(l, amqpCh, userService)
	driverConsumer := driverconsumer.New(l, amqpCh, driverService)

	app := &App{
		l:              l,
		userConsumer:   userConsumer,
		driverConsumer: driverConsumer,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) {
	go func() {
		a.userConsumer.Consume(ctx)
	}()

	go func() {
		a.driverConsumer.Consume(ctx)
	}()

	<-ctx.Done()
}
