package postgresql

import (
	"ChaikaGoods/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	Begin(ctx context.Context) (pgx.Tx, error)
}

// NewClient создает новый пул соединений с базой данных.
// Подключение к базе данных осуществляется с использованием предоставленных параметров подключения.
// Функция будет пытаться подключиться к базе данных заданное количество попыток.
func NewClient(ctx context.Context, conn config.StorageConfig, maxAttempts int) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conn.Host, conn.User, conn.Password, conn.Database)
	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %v", err)
	}
	//
	// Настройка дополнительных параметров конфигурации, если необходимо
	// connConfig.MaxConns = 10                        // Максимальное количество соединений
	// connConfig.MinConns = 2                         // Минимальное количество соединений
	// connConfig.HealthCheckPeriod = 30 * time.Second // Период проверки состояния соединений

	for i := 1; i <= maxAttempts; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, connConfig)
		if err == nil {
			return pool, nil // Успешное подключение
		}

		if i < maxAttempts {
			log.Printf("Failed to connect to database, attempt %d/%d: %v", i, maxAttempts, err)
			time.Sleep(5 * time.Second) // Пауза перед следующей попыткой
		}
	}

	return nil, fmt.Errorf("failed to connect to the database after %d attempts: %v", maxAttempts, err)
}
