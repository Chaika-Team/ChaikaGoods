package main

import (
	"ChaikaGoods/internal/config"
	repo "ChaikaGoods/internal/repository/postgresql"
	"context"
	"flag"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"os"
)

//	@title			ChaikaGoods API
//	@version		0.0.1
//	@description	This is a simple API to manage goods for the Chaika app.
//	@BasePath		/
//	@schemes		http
//	@produce		json
//	@consumes		json

//	@contact.name	Chaika
//	@contact.email	chaika.contact@yandex.ru
//	@license.name	MIT
//	@license.url	http://opensource.org/licenses/MIT

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger, "service", "ChaikaGoods", "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
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
	//Init database client
	pool, err := repo.NewClient(ctx, cfg.Storage, 5)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to connect to the database", "err", err)
		return
	}
	defer pool.Close()

	//Init repository
	repository := repo.NewGoodsRepository(pool, logger)
	if v, err := repository.GetCurrentDevVersion(ctx); err != nil {
		_ = level.Error(logger).Log("message", "Failed to get current dev version", "err", err)
	} else {
		_ = level.Info(logger).Log("message", "Current dev version", "version", v.VersionID)
	}

}
