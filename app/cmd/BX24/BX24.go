package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/bulgil/pravv-bx24/app/internal/bx24"
	"github.com/bulgil/pravv-bx24/app/internal/config"
	"github.com/bulgil/pravv-bx24/app/package/logger"
)

func main() {
	cfg := config.GetConfig()
	log := logger.New(cfg.Env, cfg.LogFolder)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bx24 := bx24.New(log, bx24.BX24Opts{
		Domain:                        cfg.BX24.Domain,
		ClientTimeout:                 cfg.BX24.Timeout,
		ClientRequestCounterMax:       cfg.BX24.RequestCounter.Max,
		ClientRequestCounterDecrement: cfg.BX24.RequestCounter.Decrement,
		ServerHost:                    cfg.BX24.ServerHost,
		ServerPort:                    cfg.BX24.ServerPort,
	})
	go bx24.Run(ctx)

	<-ctx.Done()
	time.Sleep(time.Second)
	log.Info("app closed")
}
