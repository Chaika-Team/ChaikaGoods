//go:build integration
// +build integration

package integration

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистка таблицы перед тестом
	_, err = pool.Exec(ctx, `DELETE FROM product`)
	require.NoError(t, err)

	product := &models.Product{
		Name:        "New Product",
		Description: "Description",
		Price:       99.99,
		ImageURL:    "http://example.com/image.jpg",
		SKU:         "SKU123",
	}

	t.Run("Успешное создание", func(t *testing.T) {
		id, err := repo.CreateProduct(ctx, product)
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})

	t.Run("Конфликт по SKU", func(t *testing.T) {
		// Создаём временное уникальное ограничение
		_, err := pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS test_unique_sku ON product(sku);`)
		require.NoError(t, err)

		// Повторная вставка с тем же SKU
		_, err = repo.CreateProduct(ctx, product)
		assert.Error(t, err, "ожидалась ошибка, но err == nil")

		if err != nil {
			var pgErr *pgconn.PgError
			ok := errors.As(err, &pgErr)
			assert.True(t, ok)
			assert.Equal(t, "23505", pgErr.Code) // Unique violation
		}

		// Удаляем индекс после теста
		_, _ = pool.Exec(ctx, `DROP INDEX IF EXISTS test_unique_sku;`)
	})

	t.Run("Ошибка соединения", func(t *testing.T) {
		pool.Close()
		newProduct := &models.Product{
			Name:        "P",
			Description: "D",
			Price:       1.0,
			ImageURL:    "url",
			SKU:         "SKU_OTHER",
		}
		_, err := repo.CreateProduct(ctx, newProduct)
		assert.Error(t, err)
	})
}
