package bootstrap

import (
	"context"
	"fmt"
	"time"

	trippublisher "github.com/korroziea/taxi/driver-service/internal/adapter/amqp/publisher/trip"
	"github.com/korroziea/taxi/driver-service/internal/config"
	"github.com/korroziea/taxi/driver-service/internal/consumer"
	tripconsumer "github.com/korroziea/taxi/driver-service/internal/consumer/trip"
	"github.com/korroziea/taxi/driver-service/internal/handler"
	driverhndl "github.com/korroziea/taxi/driver-service/internal/handler/driver"
	"github.com/korroziea/taxi/driver-service/internal/handler/middleware"
	"github.com/korroziea/taxi/driver-service/internal/repository/psql"
	driverrepo "github.com/korroziea/taxi/driver-service/internal/repository/psql/driver"
	triprepo "github.com/korroziea/taxi/driver-service/internal/repository/psql/trip"
	"github.com/korroziea/taxi/driver-service/internal/repository/redis"
	drivercache "github.com/korroziea/taxi/driver-service/internal/repository/redis/driver"
	httpserver "github.com/korroziea/taxi/driver-service/internal/server/http"
	driversrv "github.com/korroziea/taxi/driver-service/internal/service/driver"
	tripsrv "github.com/korroziea/taxi/driver-service/internal/service/trip"
	"github.com/korroziea/taxi/driver-service/pkg/hashing"
	"go.uber.org/zap"
)

const shutdownTimeout = 3 * time.Second

type App struct {
	l            *zap.Logger
	srv          *httpserver.Server
	tripConsumer *tripconsumer.Consumer
}

func New(l *zap.Logger, cfg config.Config) (*App, error) {
	postgresDB, _, err := psql.Connect(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("psql.Connect: %w", err)
	}

	redisDB, err := redis.Connect(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis.Connect: %w", err)
	}

	amqpConn, _, err := consumer.Connect(cfg.AMQP)
	if err != nil {
		return nil, fmt.Errorf("amqp.Connect: %w", err)
	}

	driverRepo := driverrepo.New(postgresDB)
	tripRepo := triprepo.New(postgresDB)
	cache := drivercache.New(redisDB)

	tripPublisher := trippublisher.New(amqpConn)

	argon := hashing.New(cfg.Hashing)

	driverService := driversrv.New(argon, driverRepo)
	tripService := tripsrv.New(tripRepo, driverRepo, tripPublisher)

	tripConsumer := tripconsumer.New(l, amqpConn, tripService)

	authMiddleware := middleware.New(cfg, cache)

	driverHandler := driverhndl.New(l, cfg, authMiddleware, cache, driverService)

	handler := handler.New(driverHandler).InitRoutes()

	srv := httpserver.New(cfg.HTTPPort, handler)

	app := &App{
		l:            l,
		srv:          srv,
		tripConsumer: tripConsumer,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) {
	a.l.Info("Application is started")

	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			a.l.Error("srv.ListenAndServe failed", zap.Error(err))
		}
	}()

	go func() {
		a.tripConsumer.Consume(ctx)
	}()

	go func() {
		a.tripConsumer.ConsumeCancelTrip(ctx)
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.l.Error("srv.Shutdown failed", zap.Error(err))
	}
}
