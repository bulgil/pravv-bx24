package bx24

import (
	"context"
	"time"

	"github.com/bulgil/pravv-bx24/app/internal/transport/client"
	"github.com/bulgil/pravv-bx24/app/internal/transport/server"
	"github.com/bulgil/pravv-bx24/app/package/logger"
)

type BX24 struct {
	client *client.Client
	server *server.Server
	logger logger.Logger

	domain           string
	mainHandlerRoute string
}

type BX24Opts struct {
	Domain           string
	MainHandlerRoute string

	ClientTimeout                 time.Duration
	ClientRequestCounterMax       int
	ClientRequestCounterDecrement int

	ServerHost string
	ServerPort string
}

func New(log logger.Logger, opts BX24Opts) *BX24 {
	client := client.New(log, client.ClientOpts{
		Timeout:          opts.ClientTimeout,
		CounterMax:       opts.ClientRequestCounterMax,
		CounterDecrement: opts.ClientRequestCounterDecrement,
	})

	server := server.New(log, server.ServerOption{
		Host: opts.ServerHost,
		Port: opts.ServerPort,
	})
	server.RegisterRoute("POST", opts.MainHandlerRoute)

	return &BX24{
		client:           client,
		server:           server,
		logger:           log,
		domain:           opts.Domain,
		mainHandlerRoute: opts.MainHandlerRoute,
	}
}

func (bx BX24) Run(ctx context.Context) error {
	go bx.client.Run(ctx)
	go bx.server.Run(ctx)

	bx.logger.Info("bx24 is running")

	<-ctx.Done()
	return nil
}
