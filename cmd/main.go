package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chaika-Team/ChaikaGoods/internal/config"
	"github.com/Chaika-Team/ChaikaGoods/internal/handler"
	repo "github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"

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
	// Определение флага -config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yml", "Path to configuration file")
	flag.Parse()

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

	// Загрузка конфигурации
	cfg := config.GetConfigWithPath(logger, configPath)

	// Корневой контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация клиента базы данных
	pool, err := repo.NewClient(ctx, cfg.Storage, logger)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to connect to the database", "err", err)
		return
	}
	defer pool.Close()
	_ = level.Info(logger).Log("message", "Connection to the database is successful")

	// Создание нового сервиса
	var svc service.Service
	{
		rep := repo.NewGoodsRepository(pool, logger)
		svc = service.NewService(rep, logger)
	}

	// Канал для ошибок
	errs := make(chan error)

	// Горутин для обработки сигналов завершения
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		_ = level.Info(logger).Log("message", "Received signal", "signal", sig)
		cancel() // Отмена контекста при необходимости
		errs <- fmt.Errorf("service stopped due to received signal: %s", sig)
	}()

	_ = level.Info(logger).Log("message", "Service is ready to listen and serve", "type", cfg.Listen.Type, "bind_ip", cfg.Listen.BindIP, "port", cfg.Listen.Port)

	// Создание эндпоинтов
	endpoints := handler.MakeEndpoints(log.With(logger, "endpoint"), svc)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		address := cfg.Listen.BindIP + ":" + cfg.Listen.Port
		_ = level.Info(logger).Log("message", "HTTP server is starting", "address", address)
		httpHandler := handler.NewHTTPServer(logger, endpoints)
		serverErr := http.ListenAndServe(address, httpHandler)
		if serverErr != nil {
			errs <- serverErr
		}
	}()

	// Ожидание ошибок или сигналов завершения
	_ = level.Error(logger).Log("exit", <-errs)
}
