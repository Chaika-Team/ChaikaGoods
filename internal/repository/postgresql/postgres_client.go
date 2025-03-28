package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Chaika-Team/ChaikaGoods/internal/config"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is an interface for interacting with PostgreSQL database.
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// PGClient wraps pgxpool.Pool to implement the Client interface.
type PGClient struct {
	pool   *pgxpool.Pool
	logger log.Logger
}

// NewClient creates a new PostgreSQL client with connection pooling.
// It attempts to connect to the database up to maxAttempts times with a delay between attempts.
func NewClient(ctx context.Context, cfg config.StorageConfig, logger log.Logger) (Client, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}

	// Настройка дополнительных параметров пула из конфигурации
	connConfig.MaxConns = cfg.MaxConns
	connConfig.MinConns = cfg.MinConns
	connConfig.MaxConnLifetime = cfg.MaxConnLifetime
	connConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	connConfig.HealthCheckPeriod = cfg.HealthCheckPeriod
	connConfig.MaxConnLifetimeJitter = cfg.MaxConnLifetimeJitter

	var pool *pgxpool.Pool
	for i := 1; i <= cfg.MaxAttempts; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, connConfig)
		if err == nil {
			// Тестирование подключения
			err = pool.Ping(ctx)
			if err == nil {
				_ = level.Info(logger).Log("msg", "Connected to PostgreSQL database", "name", cfg.Database)
				return &PGClient{pool: pool, logger: logger}, nil
			}
			pool.Close()
		}

		if i < cfg.MaxAttempts {
			_ = level.Warn(logger).Log("msg", "Failed to connect to database", "attempt", i, "max_attempts", cfg.MaxAttempts, "err", err)
			time.Sleep(5 * time.Second) // Пауза перед следующей попыткой
		}
	}

	return nil, fmt.Errorf("failed to connect to the database after %d attempts: %w", cfg.MaxAttempts, err)
}

// Exec executes a query without returning any rows.
func (c *PGClient) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, arguments...)
}

// Query sends a query to the database and returns the rows.
func (c *PGClient) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

// QueryRow sends a query to the database and returns a single row.
func (c *PGClient) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

// SendBatch sends a batch of queries to the database.
func (c *PGClient) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.pool.SendBatch(ctx, b)
}

// Begin starts a transaction.
func (c *PGClient) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

// Close closes the connection pool.
func (c *PGClient) Close() {
	c.pool.Close()
}
