//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
)

func TestGetProductByID(t *testing.T) {
	ctx := context.Background()

	// Подключение к БД напрямую (без сервиса)
	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Arrange: ручной INSERT валидного продукта
	const productID int64 = 42
	const insertSQL = `
		INSERT INTO product (id, name, description, price, imageurl, sku)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO NOTHING;
	`
	_, err = pool.Exec(ctx, insertSQL,
		productID, "Test Product", "Description", 100.0, "http://img.jpg", "SKU42")
	require.NoError(t, err)

	// Case 1: Продукт найден
	product, err := repo.GetProductByID(ctx, productID)
	assert.NoError(t, err)
	assert.Equal(t, productID, product.ID)
	assert.Equal(t, "Test Product", product.Name)

	// Case 2: Продукт не найден
	_, err = repo.GetProductByID(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
