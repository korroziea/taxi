package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/korroziea/taxi/driver-service/internal/config"
	"github.com/korroziea/taxi/driver-service/internal/handler"
	driverhndl "github.com/korroziea/taxi/driver-service/internal/handler/driver"
	"github.com/korroziea/taxi/driver-service/internal/repository/psql"
	driverrepo "github.com/korroziea/taxi/driver-service/internal/repository/psql/driver"
	"github.com/korroziea/taxi/driver-service/internal/repository/redis"
	drivercache "github.com/korroziea/taxi/driver-service/internal/repository/redis/driver"
	httpserver "github.com/korroziea/taxi/driver-service/internal/server/http"
	driversrv "github.com/korroziea/taxi/driver-service/internal/service/driver"
	"github.com/korroziea/taxi/driver-service/pkg/hashing"
	"go.uber.org/zap"
)

const shutdownTimeout = 3 * time.Second

type App struct {
	l   *zap.Logger
	srv *httpserver.Server
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

	driverRepo := driverrepo.New(postgresDB)
	cache := drivercache.New(redisDB)

	argon := hashing.New(cfg.Hashing)

	driverService := driversrv.New(argon, driverRepo)

	driverHandler := driverhndl.New(l, cfg, cache, driverService)

	handler := handler.New(driverHandler).InitRoutes()

	srv := httpserver.New(cfg.HTTPPort, handler)

	app := &App{
		l:   l,
		srv: srv,
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

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.l.Error("srv.Shutdown failed", zap.Error(err))
	}
}
