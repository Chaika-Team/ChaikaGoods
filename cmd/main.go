package main

import (
	"ChaikaGoods/internal/config"
	"ChaikaGoods/internal/handler"
	repo "ChaikaGoods/internal/repository/postgresql"
	"ChaikaGoods/internal/service"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

//	@title			ChaikaGoods API
//	@version		0.0.2
//	@description	This is a simple API to manage goods for the Chaika app.
//	@host			127.0.0.1:8080
//	@BasePath		/
//	@schemes		http
//	@produce		json
//	@consumes		json

//	@contact.name	Chaika
//	@contact.email	chaika.contact@yandex.ru
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger, "service", "explorer", "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
	_ = level.Info(logger).Log("message", "Service started")
	defer func(info log.Logger, keyvals ...interface{}) {
		_ = info.Log(keyvals)
	}(level.Info(logger), "message", "Service ended")

	flag.Parse()
	cfg := config.GetConfig(logger)

	// Корневой контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Init database client
	pool, err := repo.NewClient(ctx, cfg.Storage, 5)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to connect to the database", "err", err)
		return
	} else {
		_ = level.Info(logger).Log("message", "Connection to the database is successful")
	}
	defer pool.Close()

	// Create new service
	var svc service.GoodsService
	{
		rep := repo.NewGoodsRepository(pool, logger)
		svc = service.NewService(rep, logger)
	}
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		_ = level.Info(logger).Log("message", "Received signal", "signal", sig)
		cancel() // Cancel context if needed
		errs <- fmt.Errorf("service stopped due to received signal: %s", sig)
	}()
	_ = level.Info(logger).Log("message", "Service is ready to listen and serve", "type", cfg.Listen.Type, "bind_ip", cfg.Listen.BindIP, "port", cfg.Listen.Port)

	endpoints := handler.MakeEndpoints(log.With(logger, "endpoint"), svc)

	go func() {
		address := cfg.Listen.BindIP + ":" + cfg.Listen.Port
		_ = level.Info(logger).Log("message", "HTTP server is starting", "address", address)
		httpHandler := handler.NewHTTPServer(logger, endpoints)
		serverErr := http.ListenAndServe(address, httpHandler)
		if serverErr != nil {
			errs <- serverErr
		}
	}()
	_ = level.Error(logger).Log("exit", <-errs)

}
