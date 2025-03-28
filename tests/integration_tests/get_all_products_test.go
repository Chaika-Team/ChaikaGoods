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

func TestGetAllProducts(t *testing.T) {
	ctx := context.Background()

	// Подключение к БД
	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистим таблицу перед тестом
	_, err = pool.Exec(ctx, `DELETE FROM product`)
	require.NoError(t, err)

	// Case 1: Пустая таблица
	products, err := repo.GetAllProducts(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 0)

	// Вставим 2 записи
	const insertSQL = `
		INSERT INTO product (id, name, description, price, imageurl, sku)
		VALUES 
			(1, 'Product 1', 'Desc 1', 10.0, 'http://img1.jpg', 'SKU1'),
			(2, 'Product 2', 'Desc 2', 20.0, 'http://img2.jpg', 'SKU2')
		ON CONFLICT (id) DO NOTHING;
	`
	_, err = pool.Exec(ctx, insertSQL)
	require.NoError(t, err)

	// Case 2: Две записи
	products, err = repo.GetAllProducts(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	assert.Equal(t, int64(1), products[0].ID)
	assert.Equal(t, "Product 1", products[0].Name)

	assert.Equal(t, int64(2), products[1].ID)
	assert.Equal(t, "Product 2", products[1].Name)
}
