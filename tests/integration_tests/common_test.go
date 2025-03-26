//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"

	"github.com/Chaika-Team/ChaikaGoods/internal/config"
	repo "github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"
)

func setupService(t *testing.T, keyspace string) service.Service {
	t.Helper()

	startupCfg := config.Config{
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

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "component", "test")

	ctx := context.Background()
	basePool, err := repo.NewClient(ctx, startupCfg.Storage, logger)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	_, err = basePool.Exec(ctx, `DROP SCHEMA IF EXISTS `+keyspace+` CASCADE`)
	if err != nil {
		t.Fatalf("failed to drop keyspace: %v", err)
	}

	_, err = basePool.Exec(ctx, `SELECT clone_schema('public', $1);`, keyspace)
	if err != nil {
		t.Fatalf("failed to clone DB schema: %v", err)
	}
	basePool.Close()

	cfg := startupCfg
	cfg.Storage.Schema = keyspace
	pool, err := repo.NewClient(ctx, cfg.Storage, logger)
	if err != nil {
		t.Fatalf("failed to connect to test keyspace: %v", err)
	}

	repository := repo.NewGoodsRepository(pool, logger)
	return service.NewService(repository, logger)
}
