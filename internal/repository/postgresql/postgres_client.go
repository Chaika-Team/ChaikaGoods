package postgresql

import (
	"ChaikaGoods/internal/config"
	"ChaikaGoods/internal/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

// Client is a subset of the pgx.Conn interface.
// It provides methods for executing SQL queries and transactions.
type Client interface {
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

	// Query sends a query to the database and returns the rows.
	// The args are for any placeholder parameters in the query.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow sends a query to the database and returns a single row.
	// The args are for any placeholder parameters in the query.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// NewClient creates a new Client from a pgx.Conn.
// It connects to the database using the provided connection details and returns a pool of connections.
// The function will attempt to connect to the database maxAttempts times before failing.
func NewClient(ctx context.Context, conn config.StorageConfig, maxAttempts int) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conn.Host, conn.User, conn.Password, conn.Database)
	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DSN: %v", err)
	}

	// Пытаемся подключиться к базе данных с заданным количеством попыток
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.NewWithConfig(ctx, connConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %v", err)
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatalf("Failed to connect to the database after %d attempts: %v", maxAttempts, err)
	}
	return pool, nil
}
