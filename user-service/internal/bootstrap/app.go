package bootstrap

import (
	"context"
	"time"

	"github.com/korroziea/taxi/user-service/internal/config"
	"github.com/korroziea/taxi/user-service/internal/handler"
	authhndl "github.com/korroziea/taxi/user-service/internal/handler/auth"
	httpserver "github.com/korroziea/taxi/user-service/internal/server/http"
	authsrv "github.com/korroziea/taxi/user-service/internal/service/auth"
	"go.uber.org/zap"
)

const shutdownTimeout = 3 * time.Second

type App struct {
	l   *zap.Logger
	cfg config.Config // todo: need it?
	srv *httpserver.Server
}

func New(l *zap.Logger, cfg config.Config) *App {
	authService := authsrv.New()

	authHandler := authhndl.New(authService)

	handler := handler.New(authHandler).InitRoutes()

	srv := httpserver.New(cfg.HTTPPort, handler)

	app := &App{
		l:   l,
		cfg: cfg,
		srv: srv,
	}

	return app
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
