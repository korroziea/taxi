package bootsrap

import (
	"context"

	"github.com/korroziea/taxi/trip-service/internal/config"
	"github.com/korroziea/taxi/trip-service/internal/consumer"
	"go.uber.org/zap"
)

type App struct {
	l   *zap.Logger
	srv *consumer.AMQPServer
}

func New(l *zap.Logger, cfg config.Config) *App {
	srv := consumer.New(l, cfg.AMQP)

	app := &App{
		l:   l,
		srv: srv,
	}

	return app
}

func (a *App) Run(ctx context.Context) {
	go func() {
		a.srv.Consume()
	}()

	<-ctx.Done()
}
