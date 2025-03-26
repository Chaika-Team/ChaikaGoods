//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/log"

	"github.com/Chaika-Team/ChaikaGoods/internal/config"
	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	repo "github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"
)

func setupService(t *testing.T, keyspace string) service.Service {
	t.Helper()

	cfg := config.Config{
		Storage: config.StorageConfig{
			Host:                  os.Getenv("DB_HOST"),
			Port:                  os.Getenv("DB_PORT"),
			User:                  os.Getenv("DB_USER"),
			Password:              os.Getenv("DB_PASS"),
			Database:              os.Getenv("DB_NAME"),
			Schema:                "public",
			MaxAttempts:           5,
			MaxConns:              5,
			MinConns:              1,
			MaxConnLifetime:       3600,
			MaxConnIdleTime:       1800,
			HealthCheckPeriod:     60,
			MaxConnLifetimeJitter: 0,
		},
	}
	logger := log.With(log.NewLogfmtLogger(os.Stderr), "component", "test")

	ctx := context.Background()
	createDB(ctx, t, cfg, keyspace, logger)

	cfg.Storage.Database = keyspace
	pool := mustConnect(ctx, t, cfg, logger)

	t.Cleanup(func() {
		cleanupDB(ctx, t, cfg, keyspace, pool, logger)
	})

	return service.NewService(repo.NewGoodsRepository(pool, logger), logger)
}

func createDB(ctx context.Context, t *testing.T, cfg config.Config, keyspace string, logger log.Logger) {
	pool := mustConnect(ctx, t, cfg, logger)
	defer pool.Close()

	query := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s",
		keyspace, cfg.Storage.Database, cfg.Storage.User)
	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
}

func mustConnect(ctx context.Context, t *testing.T, cfg config.Config, logger log.Logger) postgresql.Client {
	pool, err := repo.NewClient(ctx, cfg.Storage, logger)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	return pool
}

func cleanupDB(ctx context.Context, t *testing.T, cfg config.Config, keyspace string, pool postgresql.Client, logger log.Logger) {
	pool.Close()

	cfg.Storage.Database = "postgres"
	cleanupPool := mustConnect(ctx, t, cfg, logger)
	defer cleanupPool.Close()

	if _, err := cleanupPool.Exec(ctx, "DROP DATABASE IF EXISTS "+keyspace); err != nil {
		t.Errorf("Failed to drop test DB: %v", err)
	}
}
