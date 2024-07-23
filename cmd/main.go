package main

import (
	"ChaikaGoods/internal/config"
	"ChaikaGoods/internal/models"
	repo "ChaikaGoods/internal/repository/postgresql"
	"context"
	"database/sql"
	"flag"
	"fmt"
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
	// test repository
	product, err := repository.GetProductByID(ctx, 1)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to get product by ID", "err", err)
	}
	fmt.Printf("Product: %v\n", product)

	products, err := repository.GetAllProducts(ctx)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to get all products", "err", err)
	}
	fmt.Printf("Products: %v\n", products)
	// Clear all changes
	version, err := repository.GetCurrentDevVersion(ctx)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to get current dev version", "err", err)
	}
	changes, err := repository.GetAllChanges(ctx, version)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to get all changes", "err", err)
	}
	for _, change := range changes {
		err = repository.DeleteChange(ctx, change.ID)
		if err != nil {
			_ = level.Error(logger).Log("message", "Failed to clear change", "err", err)
		}
	}

	// Add new product, with sql.Null* fields
	product = &models.Product{
		Name:        sql.NullString{String: "Test Product", Valid: true},
		Description: sql.NullString{String: "Test Description", Valid: true},
		Price:       sql.NullFloat64{Float64: 100.0, Valid: true},
		ImageURL:    sql.NullString{String: "images/test.png", Valid: true},
		SKU:         sql.NullString{String: "TEST123", Valid: true},
	}
	err = repository.AddQueryToCreateProduct(ctx, product)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to add query to create product", "err", err)
	}
	// Get Current Dev Version
	version, err = repository.GetCurrentDevVersion(ctx)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to get current dev version", "err", err)
	}
	print("Current dev version: ", version.VersionID)
	// Apply changes
	err = repository.ApplyChanges(ctx, version)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to apply changes", "err", err)
	}
	// Get Current Dev Version
	version, err = repository.CreateNewDevVersion(ctx)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to create new dev version", "err", err)
	}
	version, err = repository.GetCurrentDevVersion(ctx)
	if err != nil {
		_ = level.Error(logger).Log("message", "Failed to get current dev version", "err", err)
	}
	print("Current dev version: ", version.VersionID)

}
