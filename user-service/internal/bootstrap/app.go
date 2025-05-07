package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/korroziea/taxi/user-service/internal/adapter/amqp"
	trippublisher "github.com/korroziea/taxi/user-service/internal/adapter/amqp/publisher/trip"
	"github.com/korroziea/taxi/user-service/internal/config"
	tripconsumer "github.com/korroziea/taxi/user-service/internal/consumer/trip"
	"github.com/korroziea/taxi/user-service/internal/handler"
	"github.com/korroziea/taxi/user-service/internal/handler/middleware"
	triphndl "github.com/korroziea/taxi/user-service/internal/handler/trip"
	userhndl "github.com/korroziea/taxi/user-service/internal/handler/user"
	wallethndl "github.com/korroziea/taxi/user-service/internal/handler/wallet"
	"github.com/korroziea/taxi/user-service/internal/repository/psql"
	userrepo "github.com/korroziea/taxi/user-service/internal/repository/psql/user"
	walletrepo "github.com/korroziea/taxi/user-service/internal/repository/psql/wallet"
	"github.com/korroziea/taxi/user-service/internal/repository/redis"
	usercache "github.com/korroziea/taxi/user-service/internal/repository/redis/user"
	httpserver "github.com/korroziea/taxi/user-service/internal/server/http"
	tripsrv "github.com/korroziea/taxi/user-service/internal/service/trip"
	usersrv "github.com/korroziea/taxi/user-service/internal/service/user"
	walletsrv "github.com/korroziea/taxi/user-service/internal/service/wallet"
	"github.com/korroziea/taxi/user-service/pkg/hashing"
	"go.uber.org/zap"
)

const shutdownTimeout = 3 * time.Second

type App struct {
	l   *zap.Logger
	cfg config.Config // todo: need it?

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

	amqpConn, amqpCh, _, err := amqp.Connect(cfg.AMQP)
	if err != nil {
		return nil, fmt.Errorf("amqp.Connect: %w", err)
	}

	userRepo := userrepo.New(postgresDB)
	walletRepo := walletrepo.New(postgresDB)
	cache := usercache.New(redisDB)

	tripPublisher := trippublisher.New(amqpConn)

	argon := hashing.New(cfg.Hashing)

	userService := usersrv.New(argon, userRepo)
	walletService := walletsrv.New(walletRepo)
	tripService := tripsrv.New(tripPublisher)

	tripConsumer := tripconsumer.New(l, amqpCh)

	authMiddleware := middleware.New(cfg, cache)

	userHandler := userhndl.New(l, cfg, cache, userService)
	walletHandler := wallethndl.New(l, authMiddleware, walletService)
	tripHandler := triphndl.New(l, authMiddleware, tripService)

	handler := handler.New(
		userHandler,
		walletHandler,
		tripHandler,
	).InitRoutes()

	srv := httpserver.New(cfg.HTTPPort, handler)

	app := &App{
		l:            l,
		cfg:          cfg,
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

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.l.Error("srv.Shutdown failed", zap.Error(err))
	}
}
