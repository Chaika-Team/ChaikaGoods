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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateProduct(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистка таблицы и подготовка уникального индекса
	_, err = pool.Exec(ctx, `DELETE FROM product`)
	require.NoError(t, err)
	_, _ = pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS test_unique_sku ON product(sku)`)

	// Вставка двух продуктов
	insertSQL := `
		INSERT INTO product (id, name, description, price, imageurl, sku)
		VALUES 
			(1, 'Product A', 'Desc A', 10.0, 'http://a.jpg', 'SKU_A'),
			(2, 'Product B', 'Desc B', 20.0, 'http://b.jpg', 'SKU_B')
		ON CONFLICT (id) DO NOTHING;
	`
	_, err = pool.Exec(ctx, insertSQL)
	require.NoError(t, err)

	t.Run("Успешное обновление", func(t *testing.T) {
		product := &models.Product{
			ID:          1,
			Name:        "Updated A",
			Description: "Updated Desc A",
			Price:       15.5,
			ImageURL:    "http://updated.jpg",
			SKU:         "SKU_A",
		}
		err := repo.UpdateProduct(ctx, product)
		assert.NoError(t, err)
	})

	t.Run("Конфликт по SKU", func(t *testing.T) {
		product := &models.Product{
			ID:          1,
			Name:        "Conflict Name",
			Description: "Conflict Desc",
			Price:       55.5,
			ImageURL:    "http://conflict.jpg",
			SKU:         "SKU_B", // уже у продукта с id=2
		}
		err := repo.UpdateProduct(ctx, product)
		var pgErr *pgconn.PgError
		assert.Error(t, err)
		assert.True(t, errors.As(err, &pgErr))
		assert.Equal(t, pgerrcode.UniqueViolation, pgErr.Code)
	})

	t.Run("Продукт не найден", func(t *testing.T) {
		product := &models.Product{
			ID:          999,
			Name:        "Ghost",
			Description: "No Desc",
			Price:       1.0,
			ImageURL:    "http://ghost.jpg",
			SKU:         "SKU_X",
		}
		err := repo.UpdateProduct(ctx, product)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	// Чистим индекс
	_, _ = pool.Exec(ctx, `DROP INDEX IF EXISTS test_unique_sku`)
}
