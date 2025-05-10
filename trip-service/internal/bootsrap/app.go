package bootsrap

import (
	"context"
	"fmt"
	"time"

	"github.com/korroziea/taxi/trip-service/internal/adapter/amqp"
	driverpublisher "github.com/korroziea/taxi/trip-service/internal/adapter/amqp/publisher/driver"
	userpublisher "github.com/korroziea/taxi/trip-service/internal/adapter/amqp/publisher/user"
	"github.com/korroziea/taxi/trip-service/internal/config"
	driverconsumer "github.com/korroziea/taxi/trip-service/internal/consumer/driver"
	userconsumer "github.com/korroziea/taxi/trip-service/internal/consumer/user"
	"github.com/korroziea/taxi/trip-service/internal/handler"
	userhandler "github.com/korroziea/taxi/trip-service/internal/handler/user"
	"github.com/korroziea/taxi/trip-service/internal/repository/psql"
	triprepo "github.com/korroziea/taxi/trip-service/internal/repository/psql/trip"
	httpserver "github.com/korroziea/taxi/trip-service/internal/server/http"
	driverservice "github.com/korroziea/taxi/trip-service/internal/service/driver"
	userservice "github.com/korroziea/taxi/trip-service/internal/service/user"
	"go.uber.org/zap"
)

const shutdownTimeout = 3 * time.Second

type App struct {
	l              *zap.Logger
	srv            *httpserver.Server
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

	userService := userservice.New(tripRepo, driverPublisher, userPublisher)
	driverService := driverservice.New(tripRepo, userPublisher)

	userConsumer := userconsumer.New(l, amqpCh, userService)
	driverConsumer := driverconsumer.New(l, amqpCh, driverService)

	userHandler := userhandler.New(l, cfg, userService)

	handler := handler.New(userHandler).InitRoutes()

	srv := httpserver.New(cfg.HTTPPort, handler)

	app := &App{
		l:   l,
		srv: srv,
		userConsumer:   userConsumer,
		driverConsumer: driverConsumer,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) {
	go func() {
		a.userConsumer.ConsumeStartTrip(ctx)
	}()

	// go func() {
	// 	a.userConsumer.ConsumeTrips(ctx)
	// }()

	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			a.l.Error("srv.ListenAndServe failed", zap.Error(err))
		}
	}()

	go func() {
		a.driverConsumer.Consume(ctx)
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.l.Error("srv.Shutdown failed", zap.Error(err))
	}
}
