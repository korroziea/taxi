package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/korroziea/taxi/trip-service/internal/bootsrap"
	"github.com/korroziea/taxi/trip-service/internal/config"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)

	var cfg config.Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal(err)
	}

	l, _ := zap.NewProduction()
	defer l.Sync()
	l.Info("Logger initialized")

	app, err := bootsrap.New(l, cfg)
	if err != nil {
		l.Fatal("Application initialization failed", zap.Error(err))

		return
	}
	l.Info("Application initialized")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		osCall := <-quitSignal
		l.Info(fmt.Sprintf("\nSystem Call: %+v", osCall))
		cancel()
	}()

	app.Run(ctx)
}
