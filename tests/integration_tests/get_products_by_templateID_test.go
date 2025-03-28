//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductsByTemplateID(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистка таблиц
	_, err = pool.Exec(ctx, `DELETE FROM packagecontent`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM package`)
	require.NoError(t, err)

	t.Run("Шаблон содержит продукты", func(t *testing.T) {
		templateID := int64(201)

		// Вставка продуктов
		_, err := pool.Exec(ctx, `
			INSERT INTO product (id, name, description, price, imageurl, sku)
			VALUES 
				(301, 'Product A', 'Desc A', 10.0, 'url', 'SKU301'),
				(302, 'Product B', 'Desc B', 20.0, 'url', 'SKU302')
			ON CONFLICT (id) DO NOTHING;
		`)
		require.NoError(t, err)

		// Вставка шаблона
		_, err = pool.Exec(ctx, `
			INSERT INTO package (packageid, packagename, description)
			VALUES ($1, $2, $3)
		`, templateID, "Filled Template", "Has products")
		require.NoError(t, err)

		// Вставка содержимого
		_, err = pool.Exec(ctx, `
			INSERT INTO packagecontent (packageid, productid, quantity)
			VALUES 
				($1, 301, 3),
				($1, 302, 7)
		`, templateID)
		require.NoError(t, err)

		// Тестируем
		contents, err := repo.GetProductsByTemplateID(ctx, templateID)
		assert.NoError(t, err)
		assert.Len(t, contents, 2)
		assert.Equal(t, int64(301), contents[0].ProductID)
		assert.Equal(t, 3, contents[0].Quantity)
		assert.Equal(t, int64(302), contents[1].ProductID)
		assert.Equal(t, 7, contents[1].Quantity)
	})

	t.Run("Шаблон пуст (без продуктов)", func(t *testing.T) {
		templateID := int64(202)

		_, err := pool.Exec(ctx, `
			INSERT INTO package (packageid, packagename, description)
			VALUES ($1, $2, $3)
		`, templateID, "Empty Template", "No products")
		require.NoError(t, err)

		contents, err := repo.GetProductsByTemplateID(ctx, templateID)
		assert.NoError(t, err)
		assert.Len(t, contents, 0)
	})

	t.Run("Ошибка соединения с БД", func(t *testing.T) {
		pool.Close()
		_, err := repo.GetProductsByTemplateID(ctx, 201)
		assert.Error(t, err)
	})
}
